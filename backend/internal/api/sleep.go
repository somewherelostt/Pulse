package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"pulse-api/internal/collectors/sleep"
	"pulse-api/internal/db"
	"pulse-api/internal/middleware"
)

// SleepHandler handles sleep data endpoints.
type SleepHandler struct {
	Pool *pgxpool.Pool
}

func NewSleepHandler(pool *pgxpool.Pool) *SleepHandler {
	return &SleepHandler{Pool: pool}
}

// LogManual accepts a manual sleep entry from the user.
func (h *SleepHandler) LogManual(w http.ResponseWriter, r *http.Request) {
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
	var body struct {
		Date        string `json:"date"`         // "2024-01-15"
		BedtimeHour int    `json:"bedtime_hour"`  // 23
		BedtimeMin  int    `json:"bedtime_min"`   // 0
		WakeHour    int    `json:"wake_hour"`     // 7
		WakeMin     int    `json:"wake_min"`      // 30
		SleepScore  *int   `json:"sleep_score"`   // optional 0-100
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Date == "" {
		writeErr(w, http.StatusBadRequest, "date and bedtime/wake required", "INVALID_REQUEST")
		return
	}
	date, err := time.Parse("2006-01-02", body.Date)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid date format", "INVALID_DATE")
		return
	}
	entry := sleep.ManualEntry{
		Date:        date,
		BedtimeHour: body.BedtimeHour,
		BedtimeMin:  body.BedtimeMin,
		WakeHour:    body.WakeHour,
		WakeMin:     body.WakeMin,
		SleepScore:  body.SleepScore,
	}
	session := entry.ToSession()
	if err := db.UpsertSleepSessions(r.Context(), h.Pool, u.ID, []sleep.SleepSession{session}); err != nil {
		writeErr(w, http.StatusInternalServerError, "failed to save sleep session", "INTERNAL_ERROR")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"date":        session.SessionDate.Format("2006-01-02"),
		"total_mins":  session.TotalMins,
		"provider":    "manual",
	})
}

// GetRange returns sleep sessions for a date range.
func (h *SleepHandler) GetRange(w http.ResponseWriter, r *http.Request) {
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
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	if fromStr == "" {
		fromStr = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if toStr == "" {
		toStr = time.Now().Format("2006-01-02")
	}
	from, _ := time.Parse("2006-01-02", fromStr)
	to, _ := time.Parse("2006-01-02", toStr)
	sessions, err := db.GetSleepSessions(r.Context(), h.Pool, u.ID, from, to)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "failed to get sleep sessions", "INTERNAL_ERROR")
		return
	}
	out := make([]map[string]interface{}, len(sessions))
	for i, s := range sessions {
		out[i] = map[string]interface{}{
			"date":       s.SessionDate.Format("2006-01-02"),
			"provider":   s.Provider,
			"total_mins": s.TotalMins,
			"rem_mins":   s.REMMins,
			"deep_mins":  s.DeepMins,
			"bedtime":    fmtTimePtr(s.Bedtime),
			"wake_time":  fmtTimePtr(s.WakeTime),
			"sleep_score": s.SleepScore,
			"hrv":        s.HRV,
			"resting_hr": s.RestingHR,
		}
	}
	writeJSON(w, http.StatusOK, out)
}

func fmtTimePtr(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.Format(time.RFC3339)
}
