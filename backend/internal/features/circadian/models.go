package circadian

import "time"

// SleepDay is the normalized per-day input for circadian feature extraction.
type SleepDay struct {
	Date        time.Time
	Bedtime     *time.Time
	WakeTime    *time.Time
	TotalMins   *int
	REMMins     *int
	DeepMins    *int
	LightMins   *int
	AwakeMins   *int
	SleepScore  *int
	HRV         *float64
	RestingHR   *float64
	IsWeekend   bool
}

// CircadianFeatures holds extracted circadian features for a single day.
type CircadianFeatures struct {
	Date                 time.Time
	SleepDurationMins    *float64
	SleepEfficiencyPct   *float64
	SleepDebtMins        *float64
	MidSleepHour         *float64
	RhythmConsistencyPct *float64
	SocialJetlagMins     *float64
	REMPct               *float64
	DeepPct              *float64
	HRV                  *float64
	RestingHR            *float64
	SleepScore           *float64
	LightHygieneScore    *float64
}

func ptrF(v float64) *float64 { return &v }
