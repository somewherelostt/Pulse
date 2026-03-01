package calendar

// AvgRecoveryMins: (work_day_mins - meeting_mins) / meeting_count = avg "breathing room" per meeting.
func AvgRecoveryMins(workDayMins float64, meetingMins float64, meetingCount int) float64 {
	if meetingCount == 0 {
		return 0
	}
	recovery := workDayMins - meetingMins
	if recovery < 0 {
		recovery = 0
	}
	return recovery / float64(meetingCount)
}
