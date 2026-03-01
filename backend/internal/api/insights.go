package api

import (
	"encoding/json"
	"net/http"

	"pulse-api/internal/middleware"
)

type InsightsHandler struct{}

func NewInsightsHandler() *InsightsHandler {
	return &InsightsHandler{}
}

func (h *InsightsHandler) Latest(w http.ResponseWriter, r *http.Request) {
	_ = middleware.GetUserID(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(nil)
}

func (h *InsightsHandler) Generate(w http.ResponseWriter, r *http.Request) {
	_ = middleware.GetUserID(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"patterns":           []interface{}{},
		"summary":            "",
		"recommendation":     "",
		"data_quality_note": "",
		"disclaimer":         "",
		"generated_at":       "",
		"model_used":         "",
	})
}
