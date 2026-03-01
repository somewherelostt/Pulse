package api

import (
	"encoding/json"
	"net/http"

	"pulse-api/internal/middleware"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	_ = middleware.GetUserID(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"id":                "",
		"supabase_uid":      "",
		"timezone":          "UTC",
		"onboarding_done":    false,
		"consent_calendar":  false,
		"work_start_hour":   9,
		"work_end_hour":     18,
	})
}

func (h *UserHandler) UpsertMe(w http.ResponseWriter, r *http.Request) {
	_ = middleware.GetUserID(r.Context())
	var body struct {
		WorkStartHour *int   `json:"work_start_hour"`
		WorkEndHour   *int   `json:"work_end_hour"`
		Timezone      *string `json:"timezone"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"id":                "",
		"supabase_uid":      "",
		"timezone":          "UTC",
		"onboarding_done":    false,
		"consent_calendar":  false,
		"work_start_hour":   9,
		"work_end_hour":     18,
	})
}
