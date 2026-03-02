package calendar

import (
	"time"
)

// ExtractForDay computes daily features from events for one day.
func ExtractForDay(events []EventRow, date time.Time, workStartHour, workEndHour int) DailyFeatures {
	loc := date.Location()
	workStart := time.Date(date.Year(), date.Month(), date.Day(), workStartHour, 0, 0, 0, loc)
	workEnd := time.Date(date.Year(), date.Month(), date.Day(), workEndHour, 0, 0, 0, loc)
	workMins := (workEndHour - workStartHour) * 60
	if workMins <= 0 {
		workMins = 1
	}

	slots := make([]EventSlot, len(events))
	var meetingMins float64
	for i := range events {
		slots[i] = EventSlot{
			Start:     events[i].Start,
			End:       events[i].End,
			Attendees: events[i].AttendeeCount,
			IsAllDay:  events[i].IsAllDay,
		}
		if !events[i].IsAllDay {
			s, e := events[i].Start, events[i].End
			if s.Before(workStart) {
				s = workStart
			}
			if e.After(workEnd) {
				e = workEnd
			}
			if s.Before(e) {
				meetingMins += e.Sub(s).Minutes()
			}
		}
	}

	density := MeetingDensityPct(slots, date, workStartHour, workEndHour)
	backToBack, avgFocus, maxFocus, _ := FragmentationAndFocus(slots, workStart, workEnd)
	meetingCount := len(events)
	fragScore := FragmentationScore(meetingCount, backToBack)
	afterHours, weekendMins := AfterHoursAndWeekend(slots, workStartHour, workEndHour, loc)
	avgRecovery := AvgRecoveryMins(float64(workMins), meetingMins, meetingCount)

	var attendeeAvg float64
	if meetingCount > 0 {
		var sum int
		for _, s := range slots {
			sum += s.Attendees
		}
		attendeeAvg = float64(sum) / float64(meetingCount)
	}
	soloPct := 1 - density
	if soloPct < 0 {
		soloPct = 0
	}

	return DailyFeatures{
		Date:                date,
		MeetingDensityPct:   ptr(density),
		MeetingCount:        ptr(meetingCount),
		AvgFocusBlockMins:   ptr(avgFocus),
		MaxFocusBlockMins:   ptr(maxFocus),
		FragmentationScore:  ptr(fragScore),
		AfterHoursMins:      ptr(afterHours),
		WeekendMeetingMins:  ptr(weekendMins),
		BackToBackCount:     ptr(backToBack),
		AvgRecoveryMins:     ptr(avgRecovery),
		AttendeeAvg:         ptr(attendeeAvg),
		SoloTimePct:         ptr(soloPct),
		SourceEventCount:    meetingCount,
	}
}

func ptr[T any](v T) *T {
	return &v
}
