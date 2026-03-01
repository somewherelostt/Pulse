package calendar

import "time"

type EventRow struct {
	Start         time.Time
	End           time.Time
	AttendeeCount int
	IsAllDay      bool
}

type DailyFeatures struct {
	Date                  time.Time
	MeetingDensityPct     *float64
	MeetingCount         *int
	AvgFocusBlockMins    *float64
	MaxFocusBlockMins    *float64
	FragmentationScore   *float64
	AfterHoursMins       *float64
	WeekendMeetingMins   *float64
	BackToBackCount      *int
	AvgRecoveryMins      *float64
	AttendeeAvg          *float64
	SoloTimePct          *float64
	SourceEventCount     int
}

type EventSlot struct {
	Start time.Time
	End   time.Time
	// clipped to work window for density
	Attendees int
	IsAllDay  bool
}
