package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"pulse-api/internal/db"
	"pulse-api/internal/middleware"
)

type MoodHandler struct {
	Pool *pgxpool.Pool
}

func NewMoodHandler(pool *pgxpool.Pool) *MoodHandler {
	return &MoodHandler{Pool: pool}
}

func (h *MoodHandler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
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
	var body struct {
		Score  *int     `json:"score"`
		Energy *int     `json:"energy"`
		Anxiety *int    `json:"anxiety"`
		Note   *string  `json:"note"`
		Tags   []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Score == nil || *body.Score < 1 || *body.Score > 10 {
		writeErr(w, http.StatusBadRequest, "score required (1-10)", "INVALID_REQUEST")
		return
	}
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	m, err := db.CreateOrUpdateMood(r.Context(), h.Pool, u.ID, date, *body.Score, body.Energy, body.Anxiety, body.Note, body.Tags)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "failed to save mood", "INTERNAL_ERROR")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":       m.ID,
		"date":     m.Date.Format("2006-01-02"),
		"score":    m.Score,
		"energy":   m.Energy,
		"anxiety":  m.Anxiety,
		"note":     m.Note,
		"tags":     m.Tags,
		"logged_at": m.LoggedAt.Format(time.RFC3339),
	})
}

func (h *MoodHandler) GetToday(w http.ResponseWriter, r *http.Request) {
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
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	m, err := db.GetMoodByDate(r.Context(), h.Pool, u.ID, date)
	if err != nil {
		writeJSON(w, http.StatusOK, nil)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":       m.ID,
		"date":     m.Date.Format("2006-01-02"),
		"score":    m.Score,
		"energy":   m.Energy,
		"anxiety":  m.Anxiety,
		"note":     m.Note,
		"tags":     m.Tags,
		"logged_at": m.LoggedAt.Format(time.RFC3339),
	})
}

func (h *MoodHandler) GetRange(w http.ResponseWriter, r *http.Request) {
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
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	if fromStr == "" || toStr == "" {
		writeErr(w, http.StatusBadRequest, "from and to required", "INVALID_REQUEST")
		return
	}
	from, _ := time.Parse("2006-01-02", fromStr)
	to, _ := time.Parse("2006-01-02", toStr)
	logs, err := db.GetMoodRange(r.Context(), h.Pool, u.ID, from, to)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "failed to get mood range", "INTERNAL_ERROR")
		return
	}
	out := make([]map[string]interface{}, len(logs))
	for i, m := range logs {
		out[i] = map[string]interface{}{
			"id":       m.ID,
			"date":     m.Date.Format("2006-01-02"),
			"score":    m.Score,
			"energy":   m.Energy,
			"anxiety":  m.Anxiety,
			"note":     m.Note,
			"tags":     m.Tags,
			"logged_at": m.LoggedAt.Format(time.RFC3339),
		}
	}
	writeJSON(w, http.StatusOK, out)
}
