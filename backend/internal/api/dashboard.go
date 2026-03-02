package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"math"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"pulse-api/internal/correlation"
	"pulse-api/internal/db"
	"pulse-api/internal/middleware"
)

type DashboardHandler struct {
	Pool *pgxpool.Pool
}

func NewDashboardHandler(pool *pgxpool.Pool) *DashboardHandler {
	return &DashboardHandler{Pool: pool}
}

func (h *DashboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	supabaseUID := middleware.GetUserID(r.Context())
	u, err := db.UserBySupabaseUID(r.Context(), h.Pool, supabaseUID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "user not found", "USER_NOT_FOUND")
		return
	}

	ctx := r.Context()
	loc, _ := time.LoadLocation(u.Timezone)
	if loc == nil {
		loc = time.UTC
	}
	now := time.Now().In(loc)
	from30 := now.AddDate(0, 0, -30)
	from7 := now.AddDate(0, 0, -7)
	from14 := now.AddDate(0, 0, -14)

	// Get mood logs (30 days)
	moodLogs, err := db.GetMoodRange(ctx, h.Pool, u.ID, from30, now)
	if err != nil {
		slog.Warn("get mood range failed", "err", err)
		moodLogs = nil
	}

	// Get daily features (30 days)
	features30, err := getDailyFeaturesRange(ctx, h.Pool, u.ID, from30, now)
	if err != nil {
		slog.Warn("get daily features failed", "err", err)
		features30 = nil
	}

	// Build timeline (30 days, one entry per day with nulls for missing)
	timeline := buildTimeline(from30, now, moodLogs, features30)

	// Compute summary metrics
	summary := computeSummary(moodLogs, features30, from7, from14)

	// Compute correlations
	var topCorrelations []correlation.CorrelationResult
	if len(moodLogs) >= 7 && len(features30) >= 7 {
		featMap := buildFeatureMap(features30)
		moodMap := buildMoodMap(moodLogs)
		results, corrErr := correlation.BuildMatrix(featMap, moodMap)
		if corrErr == nil {
			// Return top 5 significant
			for _, cr := range results {
				if cr.Significant && len(topCorrelations) < 5 {
					topCorrelations = append(topCorrelations, cr)
				}
			}
		}
	}

	// Sync status
	syncedAt, eventsFetched, syncStatus, syncErr := db.GetSyncLog(ctx, h.Pool, u.ID, "google")
	if syncErr != nil {
		syncStatus = "never"
	}

	// Latest insight
	latestInsight := getLatestInsightJSON(ctx, h.Pool, u.ID)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"id":                 u.ID,
			"onboarding_done":    u.OnboardingDone,
			"calendar_connected": syncedAt != nil,
			"timezone":           u.Timezone,
		},
		"metrics":          summary,
		"timeline":         timeline,
		"top_correlations": topCorrelations,
		"latest_insight":   latestInsight,
		"sync_status": map[string]interface{}{
			"last_synced_at": syncedAt,
			"events_fetched": eventsFetched,
			"status":         syncStatus,
		},
	})
}

// DailyFeatureDB holds a daily feature row from the database.
type DailyFeatureDB struct {
	Date               time.Time
	MeetingDensityPct  *float64
	MeetingCount       *int
	AvgFocusBlockMins  *float64
	MaxFocusBlockMins  *float64
	FragmentationScore *float64
	AfterHoursMins     *float64
	BackToBackCount    *int
	AvgRecoveryMins    *float64
	AttendeeAvg        *float64
	SoloTimePct        *float64
}

func getDailyFeaturesRange(ctx context.Context, pool *pgxpool.Pool, userID string, from, to time.Time) ([]DailyFeatureDB, error) {
	rows, err := pool.Query(ctx, `
		select date, meeting_density_pct, meeting_count, avg_focus_block_mins, max_focus_block_mins,
		       fragmentation_score, after_hours_mins, back_to_back_count, avg_recovery_mins, attendee_avg, solo_time_pct
		from public.daily_features
		where user_id = $1::uuid and date >= $2 and date <= $3
		order by date
	`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DailyFeatureDB
	for rows.Next() {
		var f DailyFeatureDB
		if err := rows.Scan(&f.Date, &f.MeetingDensityPct, &f.MeetingCount, &f.AvgFocusBlockMins, &f.MaxFocusBlockMins,
			&f.FragmentationScore, &f.AfterHoursMins, &f.BackToBackCount, &f.AvgRecoveryMins, &f.AttendeeAvg, &f.SoloTimePct); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func buildTimeline(from, to time.Time, moods []db.MoodLogRow, features []DailyFeatureDB) []map[string]interface{} {
	// Index by date string
	moodByDate := map[string]db.MoodLogRow{}
	for _, m := range moods {
		moodByDate[m.Date.Format("2006-01-02")] = m
	}
	featByDate := map[string]DailyFeatureDB{}
	for _, f := range features {
		featByDate[f.Date.Format("2006-01-02")] = f
	}
	var timeline []map[string]interface{}
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		entry := map[string]interface{}{"date": key, "has_data": false}
		if m, ok := moodByDate[key]; ok {
			entry["mood_score"] = m.Score
			entry["energy_score"] = m.Energy
			entry["has_data"] = true
		}
		if f, ok := featByDate[key]; ok {
			entry["meeting_density_pct"] = f.MeetingDensityPct
			entry["fragmentation_score"] = f.FragmentationScore
			entry["avg_focus_block_mins"] = f.AvgFocusBlockMins
			entry["after_hours_mins"] = f.AfterHoursMins
			entry["back_to_back_count"] = f.BackToBackCount
			entry["has_data"] = true
		}
		timeline = append(timeline, entry)
	}
	return timeline
}

func computeSummary(moods []db.MoodLogRow, features []DailyFeatureDB, from7, from14 time.Time) map[string]interface{} {
	// Current week averages
	var moodScores7, moodScores714 []float64
	for _, m := range moods {
		if m.Date.After(from7) {
			moodScores7 = append(moodScores7, float64(m.Score))
		} else if m.Date.After(from14) {
			moodScores714 = append(moodScores714, float64(m.Score))
		}
	}
	var density7, density714 []float64
	for _, f := range features {
		if f.Date.After(from7) {
			if f.MeetingDensityPct != nil {
				density7 = append(density7, *f.MeetingDensityPct)
			}
		} else if f.Date.After(from14) {
			if f.MeetingDensityPct != nil {
				density714 = append(density714, *f.MeetingDensityPct)
			}
		}
	}
	avgMood7 := avgSlice(moodScores7)
	avgMoodPrev := avgSlice(moodScores714)
	avgDensity7 := avgSlice(density7)
	avgDensityPrev := avgSlice(density714)
	moodDelta := avgMood7 - avgMoodPrev
	densityDelta := avgDensity7 - avgDensityPrev

	moodTrend := "stable"
	if moodDelta > 0.5 {
		moodTrend = "up"
	} else if moodDelta < -0.5 {
		moodTrend = "down"
	}
	densityTrend := "stable"
	if densityDelta > 0.05 {
		densityTrend = "up"
	} else if densityDelta < -0.05 {
		densityTrend = "down"
	}

	return map[string]interface{}{
		"avg_meeting_density_pct": avgDensity7,
		"meeting_density_trend":   densityTrend,
		"meeting_density_delta":   densityDelta,
		"avg_mood_score":          avgMood7,
		"mood_trend":              moodTrend,
		"mood_delta":              moodDelta,
		"data_days":               len(moods),
	}
}

func avgSlice(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	var sum float64
	for _, x := range v {
		sum += x
	}
	return sum / float64(len(v))
}

func buildFeatureMap(features []DailyFeatureDB) map[string][]float64 {
	m := map[string][]float64{
		"meeting_density_pct":  {},
		"fragmentation_score":  {},
		"avg_focus_block_mins": {},
		"after_hours_mins":     {},
		"back_to_back_count":   {},
		"avg_recovery_mins":    {},
		"attendee_avg":         {},
		"solo_time_pct":        {},
	}
	for _, f := range features {
		appendIfNotNil := func(key string, v *float64) {
			if v != nil {
				m[key] = append(m[key], *v)
			} else {
				m[key] = append(m[key], math.NaN())
			}
		}
		appendIfNotNilInt := func(key string, v *int) {
			if v != nil {
				m[key] = append(m[key], float64(*v))
			} else {
				m[key] = append(m[key], math.NaN())
			}
		}
		appendIfNotNil("meeting_density_pct", f.MeetingDensityPct)
		appendIfNotNil("fragmentation_score", f.FragmentationScore)
		appendIfNotNil("avg_focus_block_mins", f.AvgFocusBlockMins)
		appendIfNotNil("after_hours_mins", f.AfterHoursMins)
		appendIfNotNilInt("back_to_back_count", f.BackToBackCount)
		appendIfNotNil("avg_recovery_mins", f.AvgRecoveryMins)
		appendIfNotNil("attendee_avg", f.AttendeeAvg)
		appendIfNotNil("solo_time_pct", f.SoloTimePct)
	}
	// Filter out series that are all NaN
	for k, v := range m {
		hasVal := false
		for _, x := range v {
			if !math.IsNaN(x) {
				hasVal = true
				break
			}
		}
		if !hasVal {
			delete(m, k)
		}
	}
	return m
}

func buildMoodMap(moods []db.MoodLogRow) map[string][]float64 {
	m := map[string][]float64{"mood_score": {}, "energy_score": {}}
	for _, ml := range moods {
		m["mood_score"] = append(m["mood_score"], float64(ml.Score))
		if ml.Energy != nil {
			m["energy_score"] = append(m["energy_score"], float64(*ml.Energy))
		} else {
			m["energy_score"] = append(m["energy_score"], math.NaN())
		}
	}
	return m
}

func getLatestInsightJSON(ctx context.Context, pool *pgxpool.Pool, userID string) interface{} {
	var content []byte
	err := pool.QueryRow(ctx, `
		select content from public.llm_insights where user_id = $1::uuid order by created_at desc limit 1
	`, userID).Scan(&content)
	if err != nil || content == nil {
		return nil
	}
	var out interface{}
	if err := json.Unmarshal(content, &out); err != nil {
		return nil
	}
	return out
}
