package calendar

import (
	"time"
)

// AfterHoursMins: event overlap with after work_end (same day) or before work_start (same day).
// Weekend: event on Saturday or Sunday.
func AfterHoursAndWeekend(events []EventSlot, workStartHour, workEndHour int, loc *time.Location) (afterHoursMins, weekendMins float64) {
	for _, e := range events {
		if e.IsAllDay {
			continue
		}
		s, en := e.Start.In(loc), e.End.In(loc)
		day := time.Date(s.Year(), s.Month(), s.Day(), 0, 0, 0, 0, loc)
		workEnd := time.Date(day.Year(), day.Month(), day.Day(), workEndHour, 0, 0, 0, loc)
		workStart := time.Date(day.Year(), day.Month(), day.Day(), workStartHour, 0, 0, 0, loc)
		// After hours: part of event after workEnd or before workStart
		if en.After(workEnd) {
			afterHoursMins += en.Sub(workEnd).Minutes()
		}
		if s.Before(workStart) {
			afterHoursMins += workStart.Sub(s).Minutes()
		}
		// Weekend
		wd := s.Weekday()
		if wd == time.Saturday || wd == time.Sunday {
			dur := en.Sub(s).Minutes()
			weekendMins += dur
		}
	}
	return afterHoursMins, weekendMins
}
