package constellation

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxRoomAge      = 45 * time.Minute
	roomCleanupTick = 5 * time.Minute
	maxPeersPerRoom = 2
	dailySessionMax = 3
)

// SignalMessage is a WebRTC signaling message relayed between peers.
// The server relays these verbatim and never inspects or logs the SDP/ICE content.
type SignalMessage struct {
	Type      string `json:"type"`                // "offer"|"answer"|"ice"|"end"|"heartbeat"
	SDP       string `json:"sdp,omitempty"`       // for offer/answer
	Candidate string `json:"candidate,omitempty"` // for ice
}

// room holds exactly 2 WebSocket connections for one signaling session.
type room struct {
	id      string
	peers   [maxPeersPerRoom]*websocket.Conn
	peerCnt int
	mu      sync.Mutex
	created time.Time
}

// RoomHub manages in-memory signaling rooms.
// Rooms are created on session start and automatically expire after 45 minutes.
type RoomHub struct {
	mu    sync.RWMutex
	rooms map[string]*room
	pool  *pgxpool.Pool
}

var wsUpgrader = websocket.Upgrader{
	// Origin checking is intentionally permissive here; JWT auth
	// is enforced before the WebSocket upgrade in the HTTP handler.
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// NewRoomHub creates a hub and starts the background cleanup goroutine.
func NewRoomHub(pool *pgxpool.Pool) *RoomHub {
	h := &RoomHub{
		rooms: make(map[string]*room),
		pool:  pool,
	}
	go h.cleanup()
	return h
}

// CreateRoom registers a new room with the given UUID.
// Called by the session/start handler before returning room_id to clients.
func (h *RoomHub) CreateRoom(roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.rooms[roomID] = &room{
		id:      roomID,
		created: time.Now(),
	}
}

// HandleSignal upgrades the HTTP connection to WebSocket, adds the peer to the
// room, and relays signaling messages between the two participants.
//
// The caller (HTTP handler) is responsible for:
//   - Verifying the Supabase JWT (via auth middleware)
//   - Confirming the user is assigned to the room (ValidateRoomAccess)
//   - Passing the resulting sessionID for audit logging
func (h *RoomHub) HandleSignal(w http.ResponseWriter, r *http.Request, roomID, sessionID string) {
	h.mu.RLock()
	rm, ok := h.rooms[roomID]
	h.mu.RUnlock()
	if !ok {
		http.Error(w, `{"error":"room not found","code":"ROOM_NOT_FOUND"}`, http.StatusNotFound)
		return
	}

	rm.mu.Lock()
	if rm.peerCnt >= maxPeersPerRoom {
		rm.mu.Unlock()
		http.Error(w, `{"error":"room full","code":"ROOM_FULL"}`, http.StatusForbidden)
		return
	}
	idx := rm.peerCnt
	rm.peerCnt++
	rm.mu.Unlock()

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Warn("constellation: ws upgrade failed", "room", roomID, "err", err)
		return
	}

	rm.mu.Lock()
	rm.peers[idx] = conn
	rm.mu.Unlock()

	logSessionEvent(r.Context(), h.pool, sessionID, "joined")

	defer func() {
		conn.Close()
		rm.mu.Lock()
		rm.peers[idx] = nil
		rm.mu.Unlock()
		logSessionEvent(context.Background(), h.pool, sessionID, "disconnected")
	}()

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg SignalMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "offer", "answer", "ice":
			logSessionEvent(r.Context(), h.pool, sessionID, msg.Type)
			relay(rm, idx, msg)

		case "heartbeat":
			// No relay needed; just reset the read deadline.

		case "end":
			logSessionEvent(r.Context(), h.pool, sessionID, "end")
			relay(rm, idx, msg)
			return

		default:
			// Silently discard unknown message types.
		}
	}
}

// relay sends msg to the other peer in the room.
func relay(rm *room, senderIdx int, msg SignalMessage) {
	rm.mu.Lock()
	other := rm.peers[1-senderIdx]
	rm.mu.Unlock()
	if other == nil {
		return
	}
	if err := other.WriteJSON(msg); err != nil {
		slog.Warn("constellation: relay failed", "err", err)
	}
}

// cleanup runs periodically and closes rooms that have exceeded maxRoomAge.
func (h *RoomHub) cleanup() {
	t := time.NewTicker(roomCleanupTick)
	defer t.Stop()
	for range t.C {
		cutoff := time.Now().Add(-maxRoomAge)
		h.mu.Lock()
		for id, rm := range h.rooms {
			if rm.created.Before(cutoff) {
				rm.mu.Lock()
				for _, conn := range rm.peers {
					if conn != nil {
						_ = conn.WriteJSON(SignalMessage{Type: "end"})
						conn.Close()
					}
				}
				rm.mu.Unlock()
				delete(h.rooms, id)
				slog.Info("constellation: room expired", "room_id", id)
			}
		}
		h.mu.Unlock()
	}
}

// ---- Session DB helpers ----

// CreatePendingSession inserts a new session record without a room_id.
// The pending session is activated by the session/start endpoint.
func CreatePendingSession(ctx context.Context, pool *pgxpool.Pool, seekerID, supporterID, contextHint string, similarity float64) (string, error) {
	var id string
	err := pool.QueryRow(ctx, `
		INSERT INTO public.constellation_sessions
			(seeker_id, supporter_id, context_hint, similarity)
		VALUES ($1::uuid, $2::uuid, $3, $4)
		RETURNING id::text
	`, seekerID, supporterID, contextHint, similarity).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("create pending session: %w", err)
	}
	return id, nil
}

// StartSession sets the room_id and started_at on a pending session.
// Returns the context_hint and similarity for the API response.
func StartSession(ctx context.Context, pool *pgxpool.Pool, sessionID, seekerID, roomID string) (contextHint string, similarity float64, err error) {
	err = pool.QueryRow(ctx, `
		UPDATE public.constellation_sessions
		SET room_id    = $3::uuid,
		    started_at = now()
		WHERE id        = $1::uuid
		  AND seeker_id = $2::uuid
		  AND started_at IS NULL
		RETURNING COALESCE(context_hint, ''), COALESCE(similarity, 0)
	`, sessionID, seekerID, roomID).Scan(&contextHint, &similarity)
	if err != nil {
		return "", 0, fmt.Errorf("start session: %w", err)
	}
	return contextHint, similarity, nil
}

// EndSession marks the session as finished.
func EndSession(ctx context.Context, pool *pgxpool.Pool, sessionID, userID string) error {
	_, err := pool.Exec(ctx, `
		UPDATE public.constellation_sessions
		SET ended_at = now()
		WHERE id = $1::uuid
		  AND (seeker_id = $2::uuid OR supporter_id = $2::uuid)
		  AND ended_at IS NULL
	`, sessionID, userID)
	return err
}

// RateSession records the user's post-session rating.
func RateSession(ctx context.Context, pool *pgxpool.Pool, sessionID, userID string, rating int, wouldAgain bool) error {
	_, err := pool.Exec(ctx, `
		UPDATE public.constellation_sessions SET
			seeker_rating         = CASE WHEN seeker_id    = $2::uuid THEN $3 ELSE seeker_rating    END,
			supporter_rating      = CASE WHEN supporter_id = $2::uuid THEN $3 ELSE supporter_rating END,
			seeker_would_again    = CASE WHEN seeker_id    = $2::uuid THEN $4 ELSE seeker_would_again    END,
			supporter_would_again = CASE WHEN supporter_id = $2::uuid THEN $4 ELSE supporter_would_again END
		WHERE id = $1::uuid
	`, sessionID, userID, rating, wouldAgain)
	return err
}

// ValidateRoomAccess checks that userID is a participant in the session for roomID.
// Returns the session ID if access is granted, or an error otherwise.
func ValidateRoomAccess(ctx context.Context, pool *pgxpool.Pool, roomID, userID string) (string, error) {
	var sessionID string
	err := pool.QueryRow(ctx, `
		SELECT id::text
		FROM public.constellation_sessions
		WHERE room_id = $1::uuid
		  AND (seeker_id = $2::uuid OR supporter_id = $2::uuid)
	`, roomID, userID).Scan(&sessionID)
	if err != nil {
		return "", fmt.Errorf("room access denied: %w", err)
	}
	return sessionID, nil
}

// CheckRateLimit returns how many peer sessions the user has had in the last 24h.
// The daily limit is 3 (dailySessionMax).
func CheckRateLimit(ctx context.Context, pool *pgxpool.Pool, userID string) (int, error) {
	var count int
	err := pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM public.constellation_sessions
		WHERE (seeker_id = $1::uuid OR supporter_id = $1::uuid)
		  AND matched_at >= now() - interval '24 hours'
	`, userID).Scan(&count)
	return count, err
}

// DailySessionMax returns the maximum number of peer sessions allowed per user per day.
func DailySessionMax() int { return dailySessionMax }

// logSessionEvent records a content-free event for audit purposes.
// Only the event type and timestamp are stored — never message content.
func logSessionEvent(ctx context.Context, pool *pgxpool.Pool, sessionID, eventType string) {
	if pool == nil || sessionID == "" {
		return
	}
	_, _ = pool.Exec(ctx, `
		INSERT INTO public.constellation_session_log (session_id, event_type)
		VALUES ($1::uuid, $2)
	`, sessionID, eventType)
}
