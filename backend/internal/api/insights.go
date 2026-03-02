package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"pulse-api/internal/db"
	"pulse-api/internal/llm"
	"pulse-api/internal/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InsightsHandler struct {
	Pool      *pgxpool.Pool
	LLMClient llm.LLMClient
}

func NewInsightsHandler(pool *pgxpool.Pool, llmClient llm.LLMClient) *InsightsHandler {
	return &InsightsHandler{Pool: pool, LLMClient: llmClient}
}

func (h *InsightsHandler) Latest(w http.ResponseWriter, r *http.Request) {
	supabaseUID := middleware.GetUserID(r.Context())
	u, err := db.UserBySupabaseUID(r.Context(), h.Pool, supabaseUID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "user not found", "USER_NOT_FOUND")
		return
	}
	result := getLatestInsightJSON(r.Context(), h.Pool, u.ID)
	writeJSON(w, http.StatusOK, result)
}

func (h *InsightsHandler) Generate(w http.ResponseWriter, r *http.Request) {
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

	// Get mood logs
	moodLogs, err := db.GetMoodRange(ctx, h.Pool, u.ID, from30, now)
	if err != nil || len(moodLogs) < 3 {
		writeErr(w, http.StatusBadRequest, "not enough mood data — log at least 3 days first", "INSUFFICIENT_DATA")
		return
	}

	// Get daily features
	features, err := getDailyFeaturesRange(ctx, h.Pool, u.ID, from30, now)
	if err != nil {
		slog.Warn("get features failed", "err", err)
	}

	// Build index for fast lookup
	featByDate := map[string]DailyFeatureDB{}
	for _, f := range features {
		featByDate[f.Date.Format("2006-01-02")] = f
	}

	// Build DaySnapshot slice
	var days []llm.DaySnapshot
	for _, m := range moodLogs {
		key := m.Date.Format("2006-01-02")
		snap := llm.DaySnapshot{Date: key}
		score := float64(m.Score)
		snap.MoodScore = &score
		if m.Energy != nil {
			e := float64(*m.Energy)
			snap.EnergyScore = &e
		}
		if f, ok := featByDate[key]; ok {
			snap.MeetingDensityPct = f.MeetingDensityPct
			snap.FragmentationScore = f.FragmentationScore
			snap.AfterHoursMins = f.AfterHoursMins
			snap.BackToBackCount = f.BackToBackCount
			snap.AvgFocusBlockMins = f.AvgFocusBlockMins
			snap.AvgRecoveryMins = f.AvgRecoveryMins
			snap.AttendeeAvg = f.AttendeeAvg
		}
		days = append(days, snap)
	}

	input := llm.PatternAnalysisInput{
		UserID:   u.ID,
		Days:     days,
		Timezone: u.Timezone,
	}

	output, err := llm.AnalyzePatterns(ctx, h.LLMClient, input)
	if err != nil {
		slog.Error("LLM analysis failed", "err", err)
		writeErr(w, http.StatusServiceUnavailable, "analysis temporarily unavailable", "LLM_UNAVAILABLE")
		return
	}

	// Store in DB
	contentJSON, err := json.Marshal(output)
	if err == nil {
		weekStart := now.AddDate(0, 0, -int(now.Weekday()))
		_, _ = h.Pool.Exec(ctx, `
			insert into public.llm_insights (user_id, insight_type, week_start, content, model_used)
			values ($1::uuid, $2, $3, $4, $5)
			on conflict (user_id, insight_type, week_start) do update set content = excluded.content, model_used = excluded.model_used, created_at = now()
		`, u.ID, "pattern_analysis", weekStart, contentJSON, output.ModelUsed)
	}

	writeJSON(w, http.StatusOK, output)
}

// storeInsight helper (also used by scheduler)
func storeInsight(ctx context.Context, pool *pgxpool.Pool, userID string, output *llm.PatternAnalysisOutput) error {
	contentJSON, err := json.Marshal(output)
	if err != nil {
		return err
	}
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	_, err = pool.Exec(ctx, `
		insert into public.llm_insights (user_id, insight_type, week_start, content, model_used)
		values ($1::uuid, $2, $3, $4, $5)
		on conflict (user_id, insight_type, week_start) do update set content = excluded.content, model_used = excluded.model_used, created_at = now()
	`, userID, "pattern_analysis", weekStart, contentJSON, output.ModelUsed)
	return err
}
