// Package constellation implements peer behavioral matching and WebRTC signaling
// for the Pulse peer support feature.
package constellation

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BehavioralFingerprint is a normalized 6-dimensional behavioral profile.
// All dimensions are in [0, 1].
type BehavioralFingerprint struct {
	// SleepQuality is derived from circadian_features sleep efficiency and score.
	SleepQuality float64
	// CalendarLoad is meeting density normalized over 14 days.
	CalendarLoad float64
	// DigitalEntropy captures fragmentation / compulsive scheduling ratio.
	DigitalEntropy float64
	// MoodTrend is the 7-day mood slope normalized to [0,1].
	// 0 = strongly declining, 0.5 = flat, 1 = strongly improving.
	MoodTrend float64
	// RhythmConsistency is circadian rhythm consistency from sleep data.
	RhythmConsistency float64
	// SocialSignals is communication / attendee volume normalized.
	SocialSignals float64
}

// BuildFingerprint pulls the last 14 days of behavioral data and computes
// a normalized 6-dimensional fingerprint for the given user.
func BuildFingerprint(ctx context.Context, userID string, pool *pgxpool.Pool) (BehavioralFingerprint, error) {
	from14 := time.Now().AddDate(0, 0, -14)
	from7 := time.Now().AddDate(0, 0, -7)

	var fp BehavioralFingerprint

	// --- Calendar / digital features ---
	calRows, err := pool.Query(ctx, `
		SELECT meeting_density_pct, fragmentation_score, attendee_avg
		FROM public.daily_features
		WHERE user_id = $1::uuid AND date >= $2
	`, userID, from14)
	if err != nil {
		return fp, fmt.Errorf("query daily_features: %w", err)
	}
	defer calRows.Close()

	var totalDensity, totalFrag, totalAttendee float64
	var calCount int
	for calRows.Next() {
		var density, frag, attendee *float64
		if err := calRows.Scan(&density, &frag, &attendee); err != nil {
			continue
		}
		calCount++
		if density != nil {
			totalDensity += *density
		}
		if frag != nil {
			totalFrag += *frag
		}
		if attendee != nil {
			totalAttendee += *attendee
		}
	}
	if err := calRows.Err(); err != nil {
		return fp, err
	}
	if calCount > 0 {
		// meeting_density_pct is 0-100, normalize to 0-1
		fp.CalendarLoad = clamp01((totalDensity / float64(calCount)) / 100.0)
		// fragmentation_score is already roughly 0-1
		fp.DigitalEntropy = clamp01(totalFrag / float64(calCount))
		// attendee_avg: typical max ~10 people, normalize
		fp.SocialSignals = clamp01((totalAttendee / float64(calCount)) / 10.0)
	}

	// --- Sleep / circadian features ---
	sleepRows, err := pool.Query(ctx, `
		SELECT sleep_efficiency_pct, rhythm_consistency_pct, sleep_score
		FROM public.circadian_features
		WHERE user_id = $1::uuid AND date >= $2
	`, userID, from14)
	if err != nil {
		return fp, fmt.Errorf("query circadian_features: %w", err)
	}
	defer sleepRows.Close()

	var totalEff, totalRhythm, totalScore float64
	var sleepCount int
	for sleepRows.Next() {
		var eff, rhythm, score *float64
		if err := sleepRows.Scan(&eff, &rhythm, &score); err != nil {
			continue
		}
		sleepCount++
		if eff != nil {
			totalEff += *eff
		}
		if rhythm != nil {
			totalRhythm += *rhythm
		}
		if score != nil {
			totalScore += *score
		}
	}
	if err := sleepRows.Err(); err != nil {
		return fp, err
	}
	if sleepCount > 0 {
		// Both pct fields are 0-100; sleep_score is 0-100
		avgEff := totalEff / float64(sleepCount)
		avgScore := totalScore / float64(sleepCount)
		fp.SleepQuality = clamp01(avgEff/100.0*0.5 + avgScore/100.0*0.5)
		fp.RhythmConsistency = clamp01((totalRhythm / float64(sleepCount)) / 100.0)
	}

	// --- Mood trend (7-day linear regression slope) ---
	moodRows, err := pool.Query(ctx, `
		SELECT score FROM public.mood_logs
		WHERE user_id = $1::uuid AND date >= $2
		ORDER BY date ASC
	`, userID, from7)
	if err != nil {
		return fp, fmt.Errorf("query mood_logs: %w", err)
	}
	defer moodRows.Close()

	var scores []float64
	for moodRows.Next() {
		var s int
		if err := moodRows.Scan(&s); err != nil {
			continue
		}
		scores = append(scores, float64(s))
	}
	if err := moodRows.Err(); err != nil {
		return fp, err
	}

	switch {
	case len(scores) >= 2:
		// slope is in units of mood-points/day; mood range 1-10 so slope ~[-5,+5]
		slope := linearSlope(scores)
		fp.MoodTrend = clamp01((slope + 5.0) / 10.0)
	case len(scores) == 1:
		fp.MoodTrend = clamp01(scores[0] / 10.0)
	default:
		fp.MoodTrend = 0.5 // neutral when no data
	}

	return fp, nil
}

// FingerprintToVector returns a 6-dimensional float32 slice suitable for
// pgvector storage and similarity search.
func FingerprintToVector(f BehavioralFingerprint) []float32 {
	return []float32{
		float32(f.SleepQuality),
		float32(f.CalendarLoad),
		float32(f.DigitalEntropy),
		float32(f.MoodTrend),
		float32(f.RhythmConsistency),
		float32(f.SocialSignals),
	}
}

// FormatVector formats a float32 slice as a pgvector literal, e.g. "[0.1,0.5,...]".
// The result can be passed to pgx as a string parameter with a ::vector cast.
func FormatVector(v []float32) string {
	if len(v) == 0 {
		return "[]"
	}
	var b strings.Builder
	b.WriteByte('[')
	for i, f := range v {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "%g", f)
	}
	b.WriteByte(']')
	return b.String()
}

// ParseVector parses a pgvector literal "[x,y,z,...]" into a float32 slice.
func ParseVector(s string) []float32 {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]float32, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		var f float32
		fmt.Sscanf(p, "%f", &f)
		out = append(out, f)
	}
	return out
}

// VectorToFingerprint reconstructs a BehavioralFingerprint from a 6-dim vector.
func VectorToFingerprint(v []float32) BehavioralFingerprint {
	get := func(i int) float64 {
		if i < len(v) {
			return float64(v[i])
		}
		return 0
	}
	return BehavioralFingerprint{
		SleepQuality:      get(0),
		CalendarLoad:      get(1),
		DigitalEntropy:    get(2),
		MoodTrend:         get(3),
		RhythmConsistency: get(4),
		SocialSignals:     get(5),
	}
}

// clamp01 constrains v to [0, 1], treating NaN/Inf as 0.
func clamp01(v float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// linearSlope computes the ordinary-least-squares slope of v vs. index.
func linearSlope(v []float64) float64 {
	n := float64(len(v))
	var sx, sy, sxy, sx2 float64
	for i, y := range v {
		x := float64(i)
		sx += x
		sy += y
		sxy += x * y
		sx2 += x * x
	}
	d := n*sx2 - sx*sx
	if d == 0 {
		return 0
	}
	return (n*sxy - sx*sy) / d
}
