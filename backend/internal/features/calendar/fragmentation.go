package calendar

import (
	"sort"
	"time"
)

const FocusBlockMinMins = 15

// Focus blocks: gaps >= 15 min between events (within work day).
// FragmentationScore = back_to_back_count / meeting_count (clamped 0-1).
// BackToBack: consecutive pair with gap < 15 min.
func FragmentationAndFocus(events []EventSlot, workStart, workEnd time.Time) (
	backToBackCount int,
	avgFocusBlockMins, maxFocusBlockMins float64,
	focusBlocks []float64,
) {
	if len(events) == 0 {
		return 0, 0, 0, nil
	}
	sorted := make([]EventSlot, len(events))
	copy(sorted, events)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Start.Before(sorted[j].Start) })

	var gaps []float64
	for i := 0; i < len(sorted)-1; i++ {
		gap := sorted[i+1].Start.Sub(sorted[i].End).Minutes()
		if gap < FocusBlockMinMins {
			backToBackCount++
		} else {
			gaps = append(gaps, gap)
		}
	}
	// Gap from work start to first event
	if len(sorted) > 0 && sorted[0].Start.After(workStart) {
		gap := sorted[0].Start.Sub(workStart).Minutes()
		if gap >= FocusBlockMinMins {
			gaps = append(gaps, gap)
		}
	}
	// Gap from last event to work end
	if len(sorted) > 0 && sorted[len(sorted)-1].End.Before(workEnd) {
		gap := workEnd.Sub(sorted[len(sorted)-1].End).Minutes()
		if gap >= FocusBlockMinMins {
			gaps = append(gaps, gap)
		}
	}

	var sum float64
	var max float64
	for _, g := range gaps {
		sum += g
		if g > max {
			max = g
		}
	}
	n := len(gaps)
	if n == 0 {
		return backToBackCount, 0, 0, nil
	}
	avg := sum / float64(n)
	fragScore := float64(backToBackCount) / float64(len(sorted))
	if fragScore > 1 {
		fragScore = 1
	}
	_ = fragScore
	return backToBackCount, avg, max, gaps
}

func FragmentationScore(meetingCount, backToBackCount int) float64 {
	if meetingCount == 0 {
		return 0
	}
	v := float64(backToBackCount) / float64(meetingCount)
	if v > 1 {
		v = 1
	}
	return v
}
