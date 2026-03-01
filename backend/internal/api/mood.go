package api

import (
	"encoding/json"
	"net/http"

	"pulse-api/internal/middleware"
)

type MoodHandler struct{}

func NewMoodHandler() *MoodHandler {
	return &MoodHandler{}
}

func (h *MoodHandler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	_ = middleware.GetUserID(r.Context())
	var body struct {
		Score  *int    `json:"score"`
		Energy *int    `json:"energy"`
		Anxiety *int   `json:"anxiety"`
		Note   *string `json:"note"`
		Tags   []string `json:"tags"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       "",
		"date":     "",
		"score":    0,
		"energy":   nil,
		"anxiety":  nil,
		"note":     nil,
		"tags":     []string{},
		"logged_at": "",
	})
}

func (h *MoodHandler) GetToday(w http.ResponseWriter, r *http.Request) {
	_ = middleware.GetUserID(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(nil)
}

func (h *MoodHandler) GetRange(w http.ResponseWriter, r *http.Request) {
	_ = middleware.GetUserID(r.Context())
	_ = r.URL.Query().Get("from")
	_ = r.URL.Query().Get("to")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode([]interface{}{})
}
