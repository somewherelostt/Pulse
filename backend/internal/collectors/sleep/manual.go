package sleep

import "time"

// ManualEntry accepts user-submitted sleep data and normalizes it.
type ManualEntry struct {
	Date          time.Time
	BedtimeHour   int // e.g. 23 for 11pm
	BedtimeMin    int
	WakeHour      int // e.g. 7 for 7am
	WakeMin       int
	SleepScore    *int     // optional self-reported quality 0-100
	Note          *string
}

// ToSession converts a manual entry to a SleepSession.
func (m ManualEntry) ToSession() SleepSession {
	loc := m.Date.Location()
	year, month, day := m.Date.Date()

	// Bedtime is on the night of Date; if hour >= 12 it's that night, else next day (crossed midnight)
	bedtimeDay := day
	if m.BedtimeHour < 12 {
		// Went to bed after midnight — bedtime is Date itself (or next day)
		bedtimeDay = day
	}
	bedtime := time.Date(year, month, bedtimeDay, m.BedtimeHour, m.BedtimeMin, 0, 0, loc)

	// Wake time: if wake hour is in the morning (< bedtime hour modulo) it is next day
	wakeDate := m.Date.AddDate(0, 0, 1)
	yr2, mo2, dy2 := wakeDate.Date()
	wakeTime := time.Date(yr2, mo2, dy2, m.WakeHour, m.WakeMin, 0, 0, loc)
	if wakeTime.Before(bedtime) {
		wakeTime = wakeTime.AddDate(0, 0, 1)
	}

	totalMins := int(wakeTime.Sub(bedtime).Minutes())
	if totalMins < 0 {
		totalMins = 0
	}

	s := SleepSession{
		Provider:    "manual",
		SessionDate: m.Date,
		Bedtime:     &bedtime,
		WakeTime:    &wakeTime,
		TotalMins:   &totalMins,
		SleepScore:  m.SleepScore,
		SourceID:    m.Date.Format("2006-01-02"),
	}
	return s
}
