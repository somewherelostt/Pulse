package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UpsertDailyFeatures(ctx context.Context, pool *pgxpool.Pool, userID string, date time.Time, f DailyFeatureRow) error {
	_, err := pool.Exec(ctx, `
		insert into public.daily_features (
			user_id, date, meeting_density_pct, meeting_count, avg_focus_block_mins, max_focus_block_mins,
			fragmentation_score, after_hours_mins, weekend_meeting_mins, back_to_back_count,
			avg_recovery_mins, attendee_avg, solo_time_pct, source_event_count
		) values ($1::uuid, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		on conflict (user_id, date) do update set
			meeting_density_pct = excluded.meeting_density_pct, meeting_count = excluded.meeting_count,
			avg_focus_block_mins = excluded.avg_focus_block_mins, max_focus_block_mins = excluded.max_focus_block_mins,
			fragmentation_score = excluded.fragmentation_score, after_hours_mins = excluded.after_hours_mins,
			weekend_meeting_mins = excluded.weekend_meeting_mins, back_to_back_count = excluded.back_to_back_count,
			avg_recovery_mins = excluded.avg_recovery_mins, attendee_avg = excluded.attendee_avg,
			solo_time_pct = excluded.solo_time_pct, source_event_count = excluded.source_event_count, computed_at = now()
	`, userID, date, f.MeetingDensityPct, f.MeetingCount, f.AvgFocusBlockMins, f.MaxFocusBlockMins,
		f.FragmentationScore, f.AfterHoursMins, f.WeekendMeetingMins, f.BackToBackCount,
		f.AvgRecoveryMins, f.AttendeeAvg, f.SoloTimePct, f.SourceEventCount)
	return err
}

type DailyFeatureRow struct {
	MeetingDensityPct   *float64
	MeetingCount        *int
	AvgFocusBlockMins   *float64
	MaxFocusBlockMins   *float64
	FragmentationScore  *float64
	AfterHoursMins      *float64
	WeekendMeetingMins  *float64
	BackToBackCount     *int
	AvgRecoveryMins     *float64
	AttendeeAvg         *float64
	SoloTimePct         *float64
	SourceEventCount    int
}
