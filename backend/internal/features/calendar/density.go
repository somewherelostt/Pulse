package calendar

import (
	"time"
)

// MeetingDensityPct: work window minutes in meetings, 0-1.
// workStartHour, workEndHour are 0-23.
func MeetingDensityPct(events []EventSlot, date time.Time, workStartHour, workEndHour int) float64 {
	workMins := (workEndHour - workStartHour) * 60
	if workMins <= 0 {
		return 0
	}
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), workStartHour, 0, 0, 0, date.Location())
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), workEndHour, 0, 0, 0, date.Location())
	var meetingMins float64
	for _, e := range events {
		if e.IsAllDay {
			continue
		}
		s, en := e.Start, e.End
		if s.Before(startOfDay) {
			s = startOfDay
		}
		if en.After(endOfDay) {
			en = endOfDay
		}
		if !s.Before(en) {
			continue
		}
		meetingMins += en.Sub(s).Minutes()
	}
	v := meetingMins / float64(workMins)
	if v > 1 {
		v = 1
	}
	if v < 0 {
		v = 0
	}
	return v
}
