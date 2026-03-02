package circadian

import "math"

const targetSleepMins = 480.0 // 8 hours

// MidSleepHour returns the hour of day (fractional) at the midpoint of sleep.
// Returns nil if bedtime or wake time is missing.
func MidSleepHour(day SleepDay) *float64 {
	if day.Bedtime == nil || day.WakeTime == nil {
		return nil
	}
	mid := day.Bedtime.Add(day.WakeTime.Sub(*day.Bedtime) / 2)
	h := float64(mid.Hour()) + float64(mid.Minute())/60.0
	return &h
}

// RhythmConsistency computes a 0-1 score based on the standard deviation of
// mid-sleep hours in the window. Lower stdev → higher consistency.
// stdev 0 → 1.0, stdev ≥ 2h → 0.0
func RhythmConsistency(midSleepHours []float64) float64 {
	if len(midSleepHours) < 2 {
		return 1.0
	}
	var sum float64
	for _, h := range midSleepHours {
		sum += h
	}
	mean := sum / float64(len(midSleepHours))
	var sq float64
	for _, h := range midSleepHours {
		d := h - mean
		sq += d * d
	}
	stdev := math.Sqrt(sq / float64(len(midSleepHours)))
	// Normalize: 0h stdev → 1.0, 2h stdev → 0.0
	score := 1.0 - (stdev / 2.0)
	if score < 0 {
		score = 0
	}
	return score
}

// SleepDebt returns the rolling average sleep debt in minutes relative to targetSleepMins.
// Positive = slept less than target on average.
func SleepDebt(recentTotalMins []float64) float64 {
	if len(recentTotalMins) == 0 {
		return 0
	}
	var sum float64
	for _, m := range recentTotalMins {
		sum += m
	}
	avg := sum / float64(len(recentTotalMins))
	debt := targetSleepMins - avg
	if debt < 0 {
		return 0
	}
	return debt
}

// SleepEfficiency returns the ratio of total_sleep_mins to time_in_bed_mins.
func SleepEfficiency(day SleepDay) *float64 {
	if day.TotalMins == nil || day.Bedtime == nil || day.WakeTime == nil {
		return nil
	}
	timeInBed := day.WakeTime.Sub(*day.Bedtime).Minutes()
	if timeInBed <= 0 {
		return nil
	}
	eff := float64(*day.TotalMins) / timeInBed
	if eff > 1.0 {
		eff = 1.0
	}
	return &eff
}
