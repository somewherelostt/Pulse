package circadian

import "math"

// SocialJetlag computes the absolute difference in hours between the
// mean mid-sleep on workdays vs weekends, returned in minutes.
// Returns 0 if either set is empty.
func SocialJetlag(days []SleepDay) float64 {
	var workMids, weekendMids []float64
	for _, d := range days {
		m := MidSleepHour(d)
		if m == nil {
			continue
		}
		if d.IsWeekend {
			weekendMids = append(weekendMids, *m)
		} else {
			workMids = append(workMids, *m)
		}
	}
	if len(workMids) == 0 || len(weekendMids) == 0 {
		return 0
	}
	workMean := mean(workMids)
	weekendMean := mean(weekendMids)
	lagHours := math.Abs(weekendMean - workMean)
	return lagHours * 60
}

func mean(v []float64) float64 {
	var s float64
	for _, x := range v {
		s += x
	}
	return s / float64(len(v))
}
