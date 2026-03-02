package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"pulse-api/internal/db"
	circadianfeat "pulse-api/internal/features/circadian"
	"pulse-api/internal/llm"
	"pulse-api/internal/middleware"
)

// CircadianHandler handles circadian feature and narrative endpoints.
type CircadianHandler struct {
	Pool      *pgxpool.Pool
	LLMClient llm.LLMClient
}

func NewCircadianHandler(pool *pgxpool.Pool, llmClient llm.LLMClient) *CircadianHandler {
	return &CircadianHandler{Pool: pool, LLMClient: llmClient}
}

// Dashboard returns 30 days of circadian features for the charts.
func (h *CircadianHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	supabaseUID := middleware.GetUserID(r.Context())
	if supabaseUID == "" {
		writeErr(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}
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
	from := now.AddDate(0, 0, -30)

	features, err := db.GetCircadianFeatures(ctx, h.Pool, u.ID, from, now)
	if err != nil {
		slog.Warn("get circadian features", "err", err)
	}

	sessions, err := db.GetSleepSessions(ctx, h.Pool, u.ID, from, now)
	if err != nil {
		slog.Warn("get sleep sessions", "err", err)
	}

	// Latest circadian insight
	narrative, interventionsJSON, modelUsed, weekStart, insightErr := db.GetLatestCircadianInsight(ctx, h.Pool, u.ID)
	var insightOut interface{}
	if insightErr == nil {
		var interventions interface{}
		_ = json.Unmarshal(interventionsJSON, &interventions)
		insightOut = map[string]interface{}{
			"narrative":     narrative,
			"interventions": interventions,
			"model_used":    modelUsed,
			"week_start":    weekStart.Format("2006-01-02"),
		}
	}

	// Build summary: averages over last 7 days
	summary := buildCircadianSummary(features)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"timeline":        buildCircadianTimeline(from, now, features, sessions),
		"summary":         summary,
		"latest_insight":  insightOut,
	})
}

// GenerateNarrative runs LLM analysis on the last 14 days of sleep data.
func (h *CircadianHandler) GenerateNarrative(w http.ResponseWriter, r *http.Request) {
	supabaseUID := middleware.GetUserID(r.Context())
	if supabaseUID == "" {
		writeErr(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}
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
	from := now.AddDate(0, 0, -14)

	features, err := db.GetCircadianFeatures(ctx, h.Pool, u.ID, from, now)
	if err != nil || len(features) < 3 {
		writeErr(w, http.StatusBadRequest, "need at least 3 days of sleep data", "INSUFFICIENT_DATA")
		return
	}

	// Build LLM input
	days := make([]llm.CircadianInput, 0, len(features))
	for _, f := range features {
		days = append(days, llm.CircadianInput{
			Date:                 f.Date.Format("2006-01-02"),
			SleepDurationMins:    f.SleepDurationMins,
			SleepEfficiencyPct:   f.SleepEfficiencyPct,
			SleepDebtMins:        f.SleepDebtMins,
			MidSleepHour:         f.MidSleepHour,
			RhythmConsistencyPct: f.RhythmConsistencyPct,
			SocialJetlagMins:     f.SocialJetlagMins,
			REMPct:               f.REMPct,
			DeepPct:              f.DeepPct,
			HRV:                  f.HRV,
			RestingHR:            f.RestingHR,
			SleepScore:           f.SleepScore,
			LightHygieneScore:    f.LightHygieneScore,
		})
	}

	out, err := llm.GenerateCircadianNarrative(ctx, h.LLMClient, days)
	if err != nil {
		slog.Error("circadian LLM failed", "err", err)
		writeErr(w, http.StatusInternalServerError, "LLM analysis failed", "LLM_ERROR")
		return
	}

	// Persist
	interventionsJSON, _ := json.Marshal(out.Interventions)
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	if err := db.UpsertCircadianInsight(ctx, h.Pool, u.ID, weekStart, out.Narrative, interventionsJSON, out.ModelUsed); err != nil {
		slog.Warn("upsert circadian insight failed", "err", err)
	}

	writeJSON(w, http.StatusOK, out)
}

// ExtractFeatures triggers on-demand circadian feature extraction for the user.
func (h *CircadianHandler) ExtractFeatures(w http.ResponseWriter, r *http.Request) {
	supabaseUID := middleware.GetUserID(r.Context())
	if supabaseUID == "" {
		writeErr(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}
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
	from := now.AddDate(0, 0, -30)

	sessions, err := db.GetSleepSessions(ctx, h.Pool, u.ID, from, now)
	if err != nil || len(sessions) == 0 {
		writeErr(w, http.StatusBadRequest, "no sleep data found", "NO_DATA")
		return
	}

	// Convert db rows to circadian input struct
	rawSessions := make([]struct {
		Date       time.Time
		Bedtime    *time.Time
		WakeTime   *time.Time
		TotalMins  *int
		REMMins    *int
		DeepMins   *int
		LightMins  *int
		AwakeMins  *int
		SleepScore *int
		HRV        *float64
		RestingHR  *float64
	}, len(sessions))
	for i, s := range sessions {
		rawSessions[i].Date = s.SessionDate
		rawSessions[i].Bedtime = s.Bedtime
		rawSessions[i].WakeTime = s.WakeTime
		rawSessions[i].TotalMins = s.TotalMins
		rawSessions[i].REMMins = s.REMMins
		rawSessions[i].DeepMins = s.DeepMins
		rawSessions[i].LightMins = s.LightMins
		rawSessions[i].AwakeMins = s.AwakeMins
		rawSessions[i].SleepScore = s.SleepScore
		rawSessions[i].HRV = s.HRV
		rawSessions[i].RestingHR = s.RestingHR
	}
	sleepDays := circadianfeat.BuildSleepDays(rawSessions)
	extracted := circadianfeat.ExtractAll(sleepDays)

	// Convert to db rows
	dbRows := make([]db.CircadianFeatureRow, len(extracted))
	for i, e := range extracted {
		dbRows[i] = db.CircadianFeatureRow{
			Date:                 e.Date,
			SleepDurationMins:    e.SleepDurationMins,
			SleepEfficiencyPct:   e.SleepEfficiencyPct,
			SleepDebtMins:        e.SleepDebtMins,
			MidSleepHour:         e.MidSleepHour,
			RhythmConsistencyPct: e.RhythmConsistencyPct,
			SocialJetlagMins:     e.SocialJetlagMins,
			REMPct:               e.REMPct,
			DeepPct:              e.DeepPct,
			HRV:                  e.HRV,
			RestingHR:            e.RestingHR,
			SleepScore:           e.SleepScore,
			LightHygieneScore:    e.LightHygieneScore,
		}
	}
	if err := db.UpsertCircadianFeatures(ctx, h.Pool, u.ID, dbRows); err != nil {
		writeErr(w, http.StatusInternalServerError, "failed to store features", "INTERNAL_ERROR")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"extracted": len(extracted)})
}

// --- helpers ---

func buildCircadianTimeline(from, to time.Time, features []db.CircadianFeatureRow, sessions []db.SleepSessionRow) []map[string]interface{} {
	featByDate := map[string]db.CircadianFeatureRow{}
	for _, f := range features {
		featByDate[f.Date.Format("2006-01-02")] = f
	}
	sessByDate := map[string]db.SleepSessionRow{}
	for _, s := range sessions {
		sessByDate[s.SessionDate.Format("2006-01-02")] = s
	}
	var timeline []map[string]interface{}
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		entry := map[string]interface{}{"date": key}
		if f, ok := featByDate[key]; ok {
			entry["sleep_duration_mins"] = f.SleepDurationMins
			entry["sleep_efficiency_pct"] = f.SleepEfficiencyPct
			entry["sleep_debt_mins"] = f.SleepDebtMins
			entry["mid_sleep_hour"] = f.MidSleepHour
			entry["rhythm_consistency_pct"] = f.RhythmConsistencyPct
			entry["hrv"] = f.HRV
			entry["sleep_score"] = f.SleepScore
		}
		timeline = append(timeline, entry)
	}
	return timeline
}

func buildCircadianSummary(features []db.CircadianFeatureRow) map[string]interface{} {
	if len(features) == 0 {
		return nil
	}
	// Last 7 days
	cutoff := time.Now().AddDate(0, 0, -7)
	var durations, scores, hrv []float64
	var debt, consistency float64
	var n int
	for _, f := range features {
		if f.Date.Before(cutoff) {
			continue
		}
		n++
		if f.SleepDurationMins != nil {
			durations = append(durations, *f.SleepDurationMins)
		}
		if f.SleepScore != nil {
			scores = append(scores, *f.SleepScore)
		}
		if f.HRV != nil {
			hrv = append(hrv, *f.HRV)
		}
		if f.SleepDebtMins != nil {
			debt += *f.SleepDebtMins
		}
		if f.RhythmConsistencyPct != nil {
			consistency += *f.RhythmConsistencyPct
		}
	}
	avg := func(v []float64) *float64 {
		if len(v) == 0 {
			return nil
		}
		var s float64
		for _, x := range v {
			s += x
		}
		r := s / float64(len(v))
		return &r
	}
	var avgDebt, avgConsistency *float64
	if n > 0 {
		v := debt / float64(n)
		avgDebt = &v
		c := consistency / float64(n)
		avgConsistency = &c
	}
	return map[string]interface{}{
		"avg_sleep_duration_mins":    avg(durations),
		"avg_sleep_score":            avg(scores),
		"avg_hrv":                    avg(hrv),
		"avg_sleep_debt_mins":        avgDebt,
		"avg_rhythm_consistency_pct": avgConsistency,
		"days_with_data":             n,
	}
}
