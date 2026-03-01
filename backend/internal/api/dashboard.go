package api

import (
	"encoding/json"
	"net/http"

	"pulse-api/internal/middleware"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	_ = middleware.GetUserID(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	payload := map[string]interface{}{
		"user": map[string]interface{}{
			"onboarding_done":    false,
			"calendar_connected": false,
			"timezone":           "UTC",
		},
		"metrics": map[string]interface{}{
			"avg_meeting_density_pct":  0,
			"meeting_density_trend":    "stable",
			"avg_fragmentation_score":  0,
			"avg_mood_score":           0,
			"mood_trend":                "stable",
			"data_days":                 0,
		},
		"timeline":              []interface{}{},
		"top_correlations":       []interface{}{},
		"latest_insight":         nil,
		"sync_status": map[string]interface{}{
			"last_synced_at":   nil,
			"events_fetched":   0,
			"status":           "never",
		},
	}
	_ = json.NewEncoder(w).Encode(payload)
}
