package circadian

// LightHygiene scores bedtime timing on a 0-100 scale.
// Optimal bedtime: 21:00–23:00 → 100. Each hour deviation reduces score by 15.
// Wake time also scored: 06:00–08:00 → bonus.
func LightHygiene(day SleepDay) *float64 {
	if day.Bedtime == nil {
		return nil
	}
	bedH := float64(day.Bedtime.Hour()) + float64(day.Bedtime.Minute())/60.0

	// Normalize bedtime hour: 22:00 is ideal (score 100)
	// Midnight (0) = same as 24 for distance calculation
	if bedH < 12 {
		bedH += 24 // treat early-morning hours as "next day"
	}
	idealBed := 22.0
	deviation := bedH - idealBed
	if deviation < 0 {
		deviation = -deviation
	}
	score := 100.0 - deviation*15.0
	if score < 0 {
		score = 0
	}

	// Wake time bonus: waking 06-08 adds up to 10 points
	if day.WakeTime != nil {
		wakeH := float64(day.WakeTime.Hour()) + float64(day.WakeTime.Minute())/60.0
		if wakeH >= 6 && wakeH <= 8 {
			score += 10.0 * (1.0 - (wakeH-7.0)/1.0*(wakeH-7.0)/1.0)
			if score > 100 {
				score = 100
			}
		}
	}
	return &score
}
