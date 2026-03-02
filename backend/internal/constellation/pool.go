package constellation

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PeerPoolEntry represents a user's record in the active peer pool.
type PeerPoolEntry struct {
	// ID is the pool record UUID.
	ID string
	// UserID is the internal user UUID.
	UserID string
	// IsAvailable indicates whether this peer is currently available for matching.
	IsAvailable bool
	// IsRecovering indicates the peer is in an active recovery phase.
	IsRecovering bool
	// MoodRecovered is true when the peer's mood trend improved after a similar drift.
	MoodRecovered bool
	// JoinedAt is when the user first joined the pool.
	JoinedAt time.Time
	// LastActive is when the user last signaled availability.
	LastActive time.Time
	// OptInAt is the timestamp of explicit consent.
	OptInAt time.Time
}

// JoinPool adds or refreshes the user in the peer pool.
// It computes a fresh behavioral fingerprint and upserts the pool record.
func JoinPool(ctx context.Context, userID string, pool *pgxpool.Pool) (*PeerPoolEntry, error) {
	fp, err := BuildFingerprint(ctx, userID, pool)
	if err != nil {
		// Non-fatal: allow join with neutral fingerprint so user isn't blocked.
		fp = BehavioralFingerprint{
			SleepQuality:      0.5,
			CalendarLoad:      0.5,
			DigitalEntropy:    0.5,
			MoodTrend:         0.5,
			RhythmConsistency: 0.5,
			SocialSignals:     0.5,
		}
	}

	vec := FingerprintToVector(fp)
	vecSQL := FormatVector(vec)

	// A peer is mood_recovered when their trajectory is upward and they have
	// reasonable sleep quality — indicating they're in a good position to support.
	moodRecovered := fp.MoodTrend > 0.55 && fp.SleepQuality > 0.3

	_, err = pool.Exec(ctx, `
		INSERT INTO public.peer_pool
			(user_id, fingerprint, is_available, mood_recovered, last_active, opt_in_at)
		VALUES ($1::uuid, $2::vector, true, $3, now(), now())
		ON CONFLICT (user_id) DO UPDATE SET
			fingerprint    = EXCLUDED.fingerprint,
			is_available   = true,
			mood_recovered = $3,
			last_active    = now()
	`, userID, vecSQL, moodRecovered)
	if err != nil {
		return nil, fmt.Errorf("join pool: %w", err)
	}

	return GetPoolEntry(ctx, userID, pool)
}

// LeavePool marks the user as unavailable, keeping historical data intact.
func LeavePool(ctx context.Context, userID string, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		UPDATE public.peer_pool SET is_available = false WHERE user_id = $1::uuid
	`, userID)
	return err
}

// HeartbeatPool refreshes last_active so the availability window resets.
// Should be called by the client on a regular interval while in the pool.
func HeartbeatPool(ctx context.Context, userID string, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		UPDATE public.peer_pool SET last_active = now()
		WHERE user_id = $1::uuid AND is_available = true
	`, userID)
	return err
}

// RefreshAvailability expires users who have not been active in the last 30 minutes.
// Intended to be called by the cron scheduler every 15 minutes.
func RefreshAvailability(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		UPDATE public.peer_pool
		SET is_available = false
		WHERE is_available = true
		  AND last_active < now() - interval '30 minutes'
	`)
	return err
}

// GetPoolEntry retrieves a user's current peer pool record.
// Returns pgx.ErrNoRows if the user has not joined the pool.
func GetPoolEntry(ctx context.Context, userID string, pool *pgxpool.Pool) (*PeerPoolEntry, error) {
	var e PeerPoolEntry
	err := pool.QueryRow(ctx, `
		SELECT id::text, user_id::text,
		       is_available, is_recovering, mood_recovered,
		       joined_at, last_active, opt_in_at
		FROM public.peer_pool
		WHERE user_id = $1::uuid
	`, userID).Scan(
		&e.ID, &e.UserID,
		&e.IsAvailable, &e.IsRecovering, &e.MoodRecovered,
		&e.JoinedAt, &e.LastActive, &e.OptInAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
