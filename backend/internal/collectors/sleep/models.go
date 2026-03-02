package sleep

import "time"

// SleepSession is the normalized sleep record returned by all collectors.
type SleepSession struct {
	Provider      string
	SessionDate   time.Time // date of sleep night start
	Bedtime       *time.Time
	WakeTime      *time.Time
	TotalMins     *int
	REMMins       *int
	DeepMins      *int
	LightMins     *int
	AwakeMins     *int
	SleepScore    *int
	HRV           *float64 // RMSSD in ms
	RestingHR     *float64 // bpm
	SourceID      string
}

// Collector fetches sleep sessions for a date range.
type Collector interface {
	FetchSessions(from, to time.Time) ([]SleepSession, error)
}
