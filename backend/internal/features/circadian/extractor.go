package circadian

import "time"

// ExtractAll computes CircadianFeatures for each day in the provided window.
// days must be ordered by Date ascending.
func ExtractAll(days []SleepDay) []CircadianFeatures {
	out := make([]CircadianFeatures, 0, len(days))
	for i, day := range days {
		f := CircadianFeatures{Date: day.Date}

		// Duration
		if day.TotalMins != nil {
			v := float64(*day.TotalMins)
			f.SleepDurationMins = &v
		}

		// Efficiency
		f.SleepEfficiencyPct = SleepEfficiency(day)

		// Mid-sleep hour
		f.MidSleepHour = MidSleepHour(day)

		// Sleep debt: 7-day rolling window
		window := days[max7Start(i):i+1]
		var durationWindow []float64
		for _, w := range window {
			if w.TotalMins != nil {
				durationWindow = append(durationWindow, float64(*w.TotalMins))
			}
		}
		if len(durationWindow) > 0 {
			debt := SleepDebt(durationWindow)
			f.SleepDebtMins = &debt
		}

		// Rhythm consistency: 7-day mid-sleep hours
		var midHours []float64
		for _, w := range window {
			m := MidSleepHour(w)
			if m != nil {
				midHours = append(midHours, *m)
			}
		}
		if len(midHours) >= 2 {
			rc := RhythmConsistency(midHours)
			f.RhythmConsistencyPct = &rc
		}

		// Social jetlag: requires weekday+weekend data in window
		sj := SocialJetlag(window)
		if sj > 0 {
			f.SocialJetlagMins = &sj
		}

		// Stage percentages
		if day.TotalMins != nil && *day.TotalMins > 0 {
			total := float64(*day.TotalMins)
			if day.REMMins != nil {
				v := float64(*day.REMMins) / total
				f.REMPct = &v
			}
			if day.DeepMins != nil {
				v := float64(*day.DeepMins) / total
				f.DeepPct = &v
			}
		}

		// Biometrics
		f.HRV = day.HRV
		f.RestingHR = day.RestingHR
		if day.SleepScore != nil {
			v := float64(*day.SleepScore)
			f.SleepScore = &v
		}

		// Light hygiene
		f.LightHygieneScore = LightHygiene(day)

		out = append(out, f)
	}
	return out
}

func max7Start(i int) int {
	if i < 7 {
		return 0
	}
	return i - 6
}

// BuildSleepDays converts raw session rows to SleepDay slice, one per date.
// sessions must be ordered by session_date ascending.
// IsWeekend is computed from the session date.
func BuildSleepDays(sessions []struct {
	Date      time.Time
	Bedtime   *time.Time
	WakeTime  *time.Time
	TotalMins *int
	REMMins   *int
	DeepMins  *int
	LightMins *int
	AwakeMins *int
	SleepScore *int
	HRV       *float64
	RestingHR *float64
}) []SleepDay {
	days := make([]SleepDay, 0, len(sessions))
	for _, s := range sessions {
		wd := s.Date.Weekday()
		days = append(days, SleepDay{
			Date:       s.Date,
			Bedtime:    s.Bedtime,
			WakeTime:   s.WakeTime,
			TotalMins:  s.TotalMins,
			REMMins:    s.REMMins,
			DeepMins:   s.DeepMins,
			LightMins:  s.LightMins,
			AwakeMins:  s.AwakeMins,
			SleepScore: s.SleepScore,
			HRV:        s.HRV,
			RestingHR:  s.RestingHR,
			IsWeekend:  wd == time.Saturday || wd == time.Sunday,
		})
	}
	return days
}
