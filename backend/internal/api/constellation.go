package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"pulse-api/internal/constellation"
	"pulse-api/internal/db"
	"pulse-api/internal/llm"
	"pulse-api/internal/middleware"
)

// ConstellationHandler handles peer matching, safety checks, and WebRTC
// signaling for the Constellation peer support feature.
type ConstellationHandler struct {
	Pool      *pgxpool.Pool
	LLMClient llm.LLMClient
	Hub       *constellation.RoomHub
}

// NewConstellationHandler creates a fully initialised handler.
func NewConstellationHandler(pool *pgxpool.Pool, llmClient llm.LLMClient, hub *constellation.RoomHub) *ConstellationHandler {
	return &ConstellationHandler{Pool: pool, LLMClient: llmClient, Hub: hub}
}

// ---- POST /api/v1/constellation/join ----

// Join adds the calling user to the peer pool.
// Body: { "opt_in_confirmed": true }
func (h *ConstellationHandler) Join(w http.ResponseWriter, r *http.Request) {
	u := h.mustUser(w, r)
	if u == nil {
		return
	}

	var body struct {
		OptInConfirmed bool `json:"opt_in_confirmed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || !body.OptInConfirmed {
		writeErr(w, http.StatusBadRequest, "opt_in_confirmed must be true", "OPT_IN_REQUIRED")
		return
	}

	entry, err := constellation.JoinPool(r.Context(), u.ID, h.Pool)
	if err != nil {
		slog.Error("constellation: join pool", "err", err)
		writeErr(w, http.StatusInternalServerError, "failed to join pool", "INTERNAL_ERROR")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"pool_id": entry.ID,
		"status":  "searching",
	})
}

// ---- POST /api/v1/constellation/leave ----

// Leave marks the calling user as unavailable in the pool.
func (h *ConstellationHandler) Leave(w http.ResponseWriter, r *http.Request) {
	u := h.mustUser(w, r)
	if u == nil {
		return
	}

	if err := constellation.LeavePool(r.Context(), u.ID, h.Pool); err != nil {
		slog.Warn("constellation: leave pool", "err", err)
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "left"})
}

// ---- GET /api/v1/constellation/safety ----

// Safety runs the pre-session safety check and returns guidance on which
// support options to surface first.
func (h *ConstellationHandler) Safety(w http.ResponseWriter, r *http.Request) {
	u := h.mustUser(w, r)
	if u == nil {
		return
	}

	result, err := constellation.CheckSafety(r.Context(), u.ID, h.Pool)
	if err != nil {
		slog.Warn("constellation: safety check", "err", err)
	}
	writeJSON(w, http.StatusOK, result)
}

// ---- GET /api/v1/constellation/match ----

// Match finds the best peer for the requesting user.
// Returns a pending match (session ID used as match_id) or a retry hint.
func (h *ConstellationHandler) Match(w http.ResponseWriter, r *http.Request) {
	u := h.mustUser(w, r)
	if u == nil {
		return
	}

	// Rate-limit: 3 peer sessions per user per day.
	count, err := constellation.CheckRateLimit(r.Context(), h.Pool, u.ID)
	if err != nil {
		slog.Warn("constellation: rate limit check", "err", err)
	}
	if count >= constellation.DailySessionMax() {
		writeJSON(w, http.StatusTooManyRequests, map[string]interface{}{
			"match_found":  false,
			"retry_after":  86400,
			"reason":       "daily_limit_reached",
		})
		return
	}

	// Run the matching algorithm.
	seekerFP, candidate, err := constellation.FindMatch(r.Context(), u.ID, h.Pool)
	if err != nil {
		slog.Warn("constellation: find match", "err", err)
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"match_found": false,
			"retry_after": 60,
		})
		return
	}
	if candidate == nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"match_found": false,
			"retry_after": 60,
		})
		return
	}

	// Generate LLM match context (warm, anonymized description).
	contextHint := ""
	if h.LLMClient != nil {
		ctx, err := llm.GenerateMatchContext(
			r.Context(), h.LLMClient,
			fingerprintToLLM(seekerFP),
			fingerprintToLLM(candidate.Fingerprint),
			candidate.Similarity,
		)
		if err != nil {
			slog.Warn("constellation: match context LLM", "err", err)
		} else {
			contextHint = ctx
		}
	}

	// Create a pending session (no room yet, starts when client calls session/start).
	matchID, err := constellation.CreatePendingSession(
		r.Context(), h.Pool,
		u.ID, candidate.UserID,
		contextHint, candidate.Similarity,
	)
	if err != nil {
		slog.Error("constellation: create pending session", "err", err)
		writeErr(w, http.StatusInternalServerError, "failed to create session", "INTERNAL_ERROR")
		return
	}

	shared := constellation.SharedPatterns(seekerFP, candidate.Fingerprint)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"match_found":     true,
		"match_id":        matchID,
		"similarity":      candidate.Similarity,
		"shared_patterns": shared,
		"context_hint":    contextHint,
	})
}

// ---- POST /api/v1/constellation/session/start ----

// SessionStart activates a pending match, creates a signaling room, and
// returns the room_id so both peers can connect via WebSocket.
// Body: { "match_id": "uuid" }
func (h *ConstellationHandler) SessionStart(w http.ResponseWriter, r *http.Request) {
	u := h.mustUser(w, r)
	if u == nil {
		return
	}

	var body struct {
		MatchID string `json:"match_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.MatchID == "" {
		writeErr(w, http.StatusBadRequest, "match_id required", "INVALID_REQUEST")
		return
	}

	roomID := uuid.New().String()
	h.Hub.CreateRoom(roomID)

	contextHint, similarity, err := constellation.StartSession(
		r.Context(), h.Pool,
		body.MatchID, u.ID, roomID,
	)
	if err != nil {
		slog.Warn("constellation: start session", "err", err, "match_id", body.MatchID)
		writeErr(w, http.StatusBadRequest, "invalid or already started match", "INVALID_MATCH")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"room_id":    roomID,
		"context":    contextHint,
		"similarity": similarity,
	})
}

// ---- POST /api/v1/constellation/session/{id}/end ----

// SessionEnd marks the session as completed.
func (h *ConstellationHandler) SessionEnd(w http.ResponseWriter, r *http.Request) {
	u := h.mustUser(w, r)
	if u == nil {
		return
	}
	sessionID := chi.URLParam(r, "id")
	if err := constellation.EndSession(r.Context(), h.Pool, sessionID, u.ID); err != nil {
		slog.Warn("constellation: end session", "err", err)
		writeErr(w, http.StatusBadRequest, "session not found or already ended", "NOT_FOUND")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ended"})
}

// ---- POST /api/v1/constellation/session/{id}/rate ----

// SessionRate records the user's rating for a completed session.
// Body: { "rating": 1-5, "would_talk_again": bool }
func (h *ConstellationHandler) SessionRate(w http.ResponseWriter, r *http.Request) {
	u := h.mustUser(w, r)
	if u == nil {
		return
	}
	sessionID := chi.URLParam(r, "id")

	var body struct {
		Rating        int  `json:"rating"`
		WouldTalkAgain bool `json:"would_talk_again"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid request body", "INVALID_REQUEST")
		return
	}
	if body.Rating < 1 || body.Rating > 5 {
		writeErr(w, http.StatusBadRequest, "rating must be 1-5", "INVALID_RATING")
		return
	}

	if err := constellation.RateSession(r.Context(), h.Pool, sessionID, u.ID, body.Rating, body.WouldTalkAgain); err != nil {
		slog.Warn("constellation: rate session", "err", err)
		writeErr(w, http.StatusBadRequest, "session not found", "NOT_FOUND")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "rated"})
}

// ---- GET /api/v1/constellation/signal/{room_id} (WebSocket) ----

// Signal upgrades the connection to WebSocket for WebRTC signaling.
// Authentication is enforced via the auth middleware before this handler runs.
// The user must be assigned to the room as seeker or supporter.
func (h *ConstellationHandler) Signal(w http.ResponseWriter, r *http.Request) {
	u := h.mustUser(w, r)
	if u == nil {
		return
	}
	roomID := chi.URLParam(r, "room_id")

	// Confirm the user was assigned to this room.
	sessionID, err := constellation.ValidateRoomAccess(r.Context(), h.Pool, roomID, u.ID)
	if err != nil {
		http.Error(w, `{"error":"access denied","code":"FORBIDDEN"}`, http.StatusForbidden)
		return
	}

	h.Hub.HandleSignal(w, r, roomID, sessionID)
}

// ---- helpers ----

// mustUser resolves the calling user from the JWT context.
// Returns nil and writes an error response if resolution fails.
func (h *ConstellationHandler) mustUser(w http.ResponseWriter, r *http.Request) *db.User {
	supabaseUID := middleware.GetUserID(r.Context())
	if supabaseUID == "" {
		writeErr(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return nil
	}
	u, err := db.UserBySupabaseUID(r.Context(), h.Pool, supabaseUID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "user not found", "USER_NOT_FOUND")
		return nil
	}
	return u
}

// fingerprintToLLM converts a constellation fingerprint to the llm package type.
func fingerprintToLLM(f constellation.BehavioralFingerprint) llm.MatchFingerprint {
	return llm.MatchFingerprint{
		SleepQuality:      f.SleepQuality,
		CalendarLoad:      f.CalendarLoad,
		DigitalEntropy:    f.DigitalEntropy,
		MoodTrend:         f.MoodTrend,
		RhythmConsistency: f.RhythmConsistency,
		SocialSignals:     f.SocialSignals,
	}
}
