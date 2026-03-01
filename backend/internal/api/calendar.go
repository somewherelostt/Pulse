package api

import (
	"encoding/json"
	"net/http"

	"pulse-api/internal/middleware"
)

type CalendarHandler struct {
	FrontendURL string
}

func NewCalendarHandler(frontendURL string) *CalendarHandler {
	return &CalendarHandler{FrontendURL: frontendURL}
}

func (h *CalendarHandler) Connect(w http.ResponseWriter, r *http.Request) {
	_ = middleware.GetUserID(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"url": ""})
}

func (h *CalendarHandler) Callback(w http.ResponseWriter, r *http.Request) {
	_ = r.URL.Query().Get("code")
	_ = r.URL.Query().Get("state")
	http.Redirect(w, r, h.FrontendURL+"/onboarding?step=2&syncing=true", http.StatusFound)
}

func (h *CalendarHandler) Status(w http.ResponseWriter, r *http.Request) {
	_ = middleware.GetUserID(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"connected":      false,
		"last_synced_at": nil,
		"status":         "never",
	})
}

func (h *CalendarHandler) Sync(w http.ResponseWriter, r *http.Request) {
	_ = middleware.GetUserID(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
