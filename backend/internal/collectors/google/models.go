package google

import "time"

type Event struct {
	ID           string
	Title        string
	Start        time.Time
	End          time.Time
	Attendees    int
	IsAllDay     bool
	IsRecurring  bool
}

type RawEventRow struct {
	GoogleEventID string
	TitleHash     string
	StartTime     time.Time
	EndTime       time.Time
	AttendeeCount int
	IsAllDay      bool
	IsRecurring   bool
	IsAfterHours  bool
	IsWeekend     bool
	DurationMins  float64
}
