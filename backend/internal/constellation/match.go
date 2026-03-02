package constellation

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// MatchCandidate is an internal representation of a potential peer match.
type MatchCandidate struct {
	// PoolID is the peer_pool record UUID.
	PoolID string
	// UserID is the peer's internal user UUID (never surfaced in API responses).
	UserID string
	// Similarity is the cosine similarity [0,1] between seeker and peer fingerprints.
	Similarity float64
	// MoodRecovered indicates the peer has an improving mood trajectory.
	MoodRecovered bool
	// Fingerprint is the peer's decoded behavioral fingerprint.
	Fingerprint BehavioralFingerprint
}

// FindMatch locates the best available, recovered peer for seekerID.
//
// It returns:
//   - seekerFP: the seeker's behavioral fingerprint (for context generation)
//   - candidate: the best peer, or nil if no suitable match exists
//   - err: database or parsing errors
//
// A match is considered valid only when cosine similarity >= 0.6.
// Peers matched with the seeker in the last 7 days are excluded.
func FindMatch(ctx context.Context, seekerID string, pool *pgxpool.Pool) (BehavioralFingerprint, *MatchCandidate, error) {
	// Retrieve seeker's stored fingerprint string from the pool table.
	var seekerFPStr string
	err := pool.QueryRow(ctx, `
		SELECT fingerprint::text
		FROM public.peer_pool
		WHERE user_id = $1::uuid AND is_available = true
	`, seekerID).Scan(&seekerFPStr)
	if err != nil {
		return BehavioralFingerprint{}, nil, fmt.Errorf("seeker not in pool or unavailable: %w", err)
	}

	// Query up to 5 candidates ordered by cosine distance (ascending = most similar first).
	// Only recovered peers are eligible as supporters.
	rows, err := pool.Query(ctx, `
		SELECT pp.id::text,
		       pp.user_id::text,
		       1 - (pp.fingerprint <=> $2::vector) AS cosine_sim,
		       pp.mood_recovered,
		       pp.fingerprint::text
		FROM public.peer_pool pp
		WHERE pp.is_available    = true
		  AND pp.user_id        != $1::uuid
		  AND pp.mood_recovered  = true
		  AND pp.fingerprint    IS NOT NULL
		  AND pp.user_id NOT IN (
		      SELECT supporter_id
		        FROM public.constellation_sessions
		       WHERE seeker_id = $1::uuid
		         AND matched_at >= now() - interval '7 days'
		      UNION
		      SELECT seeker_id
		        FROM public.constellation_sessions
		       WHERE supporter_id = $1::uuid
		         AND matched_at >= now() - interval '7 days'
		  )
		ORDER BY pp.fingerprint <=> $2::vector
		LIMIT 5
	`, seekerID, seekerFPStr)
	if err != nil {
		return BehavioralFingerprint{}, nil, fmt.Errorf("query candidates: %w", err)
	}
	defer rows.Close()

	var candidates []MatchCandidate
	for rows.Next() {
		var c MatchCandidate
		var fpStr string
		if err := rows.Scan(&c.PoolID, &c.UserID, &c.Similarity, &c.MoodRecovered, &fpStr); err != nil {
			continue
		}
		c.Fingerprint = VectorToFingerprint(ParseVector(fpStr))
		candidates = append(candidates, c)
	}
	if err := rows.Err(); err != nil {
		return BehavioralFingerprint{}, nil, err
	}

	if len(candidates) == 0 {
		return VectorToFingerprint(ParseVector(seekerFPStr)), nil, nil
	}

	// Score: 70% similarity + 30% recovery quality.
	best := candidates[0]
	bestScore := compositScore(best)
	for _, c := range candidates[1:] {
		if s := compositScore(c); s > bestScore {
			bestScore = s
			best = c
		}
	}

	seekerFP := VectorToFingerprint(ParseVector(seekerFPStr))

	// Reject if below similarity threshold.
	if best.Similarity < 0.6 {
		return seekerFP, nil, nil
	}

	return seekerFP, &best, nil
}

// SharedPatterns returns the behavioral dimension labels where seeker and peer
// fingerprints are within 15% of each other (normalized distance).
func SharedPatterns(seeker, peer BehavioralFingerprint) []string {
	type dim struct {
		label string
		s, p  float64
	}
	dims := []dim{
		{"sleep patterns", seeker.SleepQuality, peer.SleepQuality},
		{"calendar load", seeker.CalendarLoad, peer.CalendarLoad},
		{"digital habits", seeker.DigitalEntropy, peer.DigitalEntropy},
		{"mood trajectory", seeker.MoodTrend, peer.MoodTrend},
		{"circadian rhythm", seeker.RhythmConsistency, peer.RhythmConsistency},
		{"social signals", seeker.SocialSignals, peer.SocialSignals},
	}
	var out []string
	for _, d := range dims {
		diff := d.s - d.p
		if diff < 0 {
			diff = -diff
		}
		if diff < 0.15 {
			out = append(out, d.label)
		}
	}
	return out
}

func compositScore(c MatchCandidate) float64 {
	recovery := 0.0
	if c.MoodRecovered {
		recovery = 1.0
	}
	return c.Similarity*0.7 + recovery*0.3
}
