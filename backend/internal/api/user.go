package api

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"pulse-api/internal/db"
	"pulse-api/internal/middleware"
)

type UserHandler struct {
	Pool *pgxpool.Pool
}

func NewUserHandler(pool *pgxpool.Pool) *UserHandler {
	return &UserHandler{Pool: pool}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	supabaseUID := middleware.GetUserID(r.Context())
	if supabaseUID == "" {
		writeErr(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}
	u, err := db.UserBySupabaseUID(r.Context(), h.Pool, supabaseUID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "failed to get user", "INTERNAL_ERROR")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":               u.ID,
		"supabase_uid":     u.SupabaseUID,
		"timezone":         u.Timezone,
		"onboarding_done":  u.OnboardingDone,
		"consent_calendar": u.ConsentCalendar,
		"work_start_hour":  u.WorkStartHour,
		"work_end_hour":    u.WorkEndHour,
	})
}

func (h *UserHandler) UpsertMe(w http.ResponseWriter, r *http.Request) {
	supabaseUID := middleware.GetUserID(r.Context())
	if supabaseUID == "" {
		writeErr(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}
	var body struct {
		WorkStartHour *int    `json:"work_start_hour"`
		WorkEndHour   *int    `json:"work_end_hour"`
		Timezone      *string `json:"timezone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid body", "INVALID_REQUEST")
		return
	}
	tz := "UTC"
	if body.Timezone != nil {
		tz = *body.Timezone
	}
	workStart, workEnd := 9, 18
	if body.WorkStartHour != nil {
		workStart = *body.WorkStartHour
	}
	if body.WorkEndHour != nil {
		workEnd = *body.WorkEndHour
	}
	u, err := db.UpsertUser(r.Context(), h.Pool, supabaseUID, tz, workStart, workEnd)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "failed to upsert user", "INTERNAL_ERROR")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":               u.ID,
		"supabase_uid":     u.SupabaseUID,
		"timezone":         u.Timezone,
		"onboarding_done":  u.OnboardingDone,
		"consent_calendar": u.ConsentCalendar,
		"work_start_hour":  u.WorkStartHour,
		"work_end_hour":    u.WorkEndHour,
	})
}
