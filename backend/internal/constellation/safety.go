package constellation

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SafetyResult guides which support options to surface to the user.
// Access is never blocked entirely; only the presentation order changes.
type SafetyResult struct {
	// ShowPeerOption indicates whether peer matching should be offered at all.
	ShowPeerOption bool `json:"show_peer_option"`
	// ShowProfessionalFirst indicates professional resources should appear
	// before the peer option (triggered by high need or drift signals).
	ShowProfessionalFirst bool `json:"show_professional_first"`
	// CrisisDetected indicates potential crisis language was found in recent
	// mood notes. When true, professional resources are surfaced prominently.
	// Peer support remains available but is secondary.
	CrisisDetected bool `json:"crisis_detected"`
}

// crisisKeywords are substrings matched case-insensitively against mood notes.
// These are sensitive signals — presence triggers ShowProfessionalFirst.
var crisisKeywords = []string{
	"suicid",
	"kill myself",
	"end my life",
	"don't want to live",
	"cant go on",
	"can't go on",
	"hopeless",
	"worthless",
	"self-harm",
	"hurt myself",
	"no reason to live",
}

// CheckSafety evaluates the last 7 days of mood data and behavioral features
// to determine appropriate support options for the user.
//
// Safety rules applied:
//   - drift_score (fragmentation) > 0.4 → recommend professional first
//   - mood_score < 3 for 3+ logged days → flag as high need
//   - average mood < 2.0 → high need
//   - crisis keywords in any recent mood note → crisis_detected
//
// This function always returns ShowPeerOption = true. Access is never blocked.
func CheckSafety(ctx context.Context, userID string, pool *pgxpool.Pool) (SafetyResult, error) {
	result := SafetyResult{ShowPeerOption: true}
	from := time.Now().AddDate(0, 0, -7)

	// Scan mood logs for scores and crisis keywords.
	moodRows, err := pool.Query(ctx, `
		SELECT score, note
		FROM public.mood_logs
		WHERE user_id = $1::uuid AND date >= $2
		ORDER BY date DESC
	`, userID, from)
	if err != nil {
		// Fail open — always allow access.
		return result, nil
	}
	defer moodRows.Close()

	var totalScore float64
	var moodCount, daysLow int
	for moodRows.Next() {
		var score int
		var note *string
		if err := moodRows.Scan(&score, &note); err != nil {
			continue
		}
		totalScore += float64(score)
		moodCount++
		if score < 3 {
			daysLow++
		}
		if note != nil {
			lower := strings.ToLower(*note)
			for _, kw := range crisisKeywords {
				if strings.Contains(lower, kw) {
					result.CrisisDetected = true
					result.ShowProfessionalFirst = true
					break
				}
			}
		}
	}
	_ = moodRows.Err()

	// 3+ days with low mood score → high need.
	if daysLow >= 3 {
		result.ShowProfessionalFirst = true
	}

	// Average mood below 2 → high need.
	if moodCount > 0 && totalScore/float64(moodCount) < 2.0 {
		result.ShowProfessionalFirst = true
	}

	// Check fragmentation drift from calendar features.
	var avgFragmentation *float64
	_ = pool.QueryRow(ctx, `
		SELECT AVG(fragmentation_score)
		FROM public.daily_features
		WHERE user_id = $1::uuid AND date >= $2
	`, userID, from).Scan(&avgFragmentation)
	if avgFragmentation != nil && *avgFragmentation > 0.4 {
		result.ShowProfessionalFirst = true
	}

	return result, nil
}
