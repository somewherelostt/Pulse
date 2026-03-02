package llm

import (
	"context"
	"fmt"
	"strings"
)

// MatchFingerprint mirrors constellation.BehavioralFingerprint to avoid
// a circular import between the llm and constellation packages.
type MatchFingerprint struct {
	SleepQuality      float64
	CalendarLoad      float64
	DigitalEntropy    float64
	MoodTrend         float64
	RhythmConsistency float64
	SocialSignals     float64
}

// GenerateMatchContext produces a warm, anonymized 2-3 sentence description
// of what the seeker and peer have in common, suitable for display in the
// Constellation UI before a support session begins.
//
// Privacy guarantees:
//   - Only pattern-level language is used (no specific values or numbers).
//   - No diagnosis or clinical framing.
//   - Content is derived purely from behavioral fingerprint dimensions.
func GenerateMatchContext(
	ctx context.Context,
	client LLMClient,
	seeker, peer MatchFingerprint,
	similarity float64,
) (string, error) {
	seekerDesc := describeFingerprint(seeker)
	peerDesc := describeFingerprint(peer)

	req := CompletionRequest{
		SystemPrompt: `You generate warm, brief match context for a peer support session.
Rules:
- Pattern-level language only. Never reference specific numbers or data values.
- Never diagnose or use clinical language.
- 2-3 sentences maximum.
- Warm and human, but not patronizing or condescending.
- Focus on shared behavioral patterns that create a natural connection.
- Write in second person ("You both...").`,
		UserPrompt: fmt.Sprintf(
			"Person seeking support — behavioral patterns: %s\n"+
				"Person offering support — behavioral patterns: %s\n"+
				"Pattern match strength: %.0f%%\n\n"+
				"Write the match context:",
			seekerDesc, peerDesc, similarity*100,
		),
		MaxTokens:   180,
		Temperature: 0.7,
	}

	resp, err := client.Complete(ctx, req)
	if err != nil {
		return "", fmt.Errorf("generate match context: %w", err)
	}
	return strings.TrimSpace(resp.Content), nil
}

// describeFingerprint converts a MatchFingerprint into a natural-language
// description suitable for the LLM prompt. Values are described qualitatively,
// never quantitatively.
func describeFingerprint(f MatchFingerprint) string {
	var traits []string

	// Sleep quality
	switch {
	case f.SleepQuality < 0.35:
		traits = append(traits, "disrupted sleep")
	case f.SleepQuality > 0.7:
		traits = append(traits, "good sleep quality")
	default:
		traits = append(traits, "variable sleep")
	}

	// Calendar load
	switch {
	case f.CalendarLoad > 0.7:
		traits = append(traits, "very high meeting load")
	case f.CalendarLoad > 0.45:
		traits = append(traits, "moderate calendar pressure")
	default:
		traits = append(traits, "lighter schedule")
	}

	// Digital entropy
	switch {
	case f.DigitalEntropy > 0.6:
		traits = append(traits, "fragmented focus time")
	case f.DigitalEntropy < 0.3:
		traits = append(traits, "good focus blocks")
	}

	// Mood trend
	switch {
	case f.MoodTrend < 0.4:
		traits = append(traits, "downward mood drift")
	case f.MoodTrend > 0.65:
		traits = append(traits, "improving mood")
	default:
		traits = append(traits, "steady mood")
	}

	// Circadian rhythm
	if f.RhythmConsistency < 0.45 {
		traits = append(traits, "irregular rhythm")
	} else {
		traits = append(traits, "consistent rhythm")
	}

	return strings.Join(traits, ", ")
}
