package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UpsertRawEvents(ctx context.Context, pool *pgxpool.Pool, userID string, events []RawCalendarEvent) error {
	for _, e := range events {
		_, err := pool.Exec(ctx, `
			insert into public.raw_calendar_events (
				user_id, google_event_id, title_hash, start_time, end_time,
				attendee_count, is_all_day, is_recurring, is_after_hours, is_weekend, duration_mins
			) values ($1::uuid, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			on conflict (user_id, google_event_id) do update set
				start_time = excluded.start_time, end_time = excluded.end_time,
				attendee_count = excluded.attendee_count, is_all_day = excluded.is_all_day,
				is_recurring = excluded.is_recurring, is_after_hours = excluded.is_after_hours,
				is_weekend = excluded.is_weekend, duration_mins = excluded.duration_mins, fetched_at = now()
		`, userID, e.GoogleEventID, e.TitleHash, e.StartTime, e.EndTime,
			e.AttendeeCount, e.IsAllDay, e.IsRecurring, e.IsAfterHours, e.IsWeekend, e.DurationMins)
		if err != nil {
			return err
		}
	}
	return nil
}

type RawCalendarEvent struct {
	GoogleEventID string
	TitleHash     *string
	StartTime     time.Time
	EndTime       time.Time
	AttendeeCount int
	IsAllDay      bool
	IsRecurring   bool
	IsAfterHours  bool
	IsWeekend     bool
	DurationMins  float64
}

func RawEventsByUserAndDateRange(ctx context.Context, pool *pgxpool.Pool, userID string, from, to time.Time) ([]RawEventByDay, error) {
	rows, err := pool.Query(ctx, `
		select start_time, end_time, attendee_count, is_all_day
		from public.raw_calendar_events
		where user_id = $1::uuid and start_time < $3 and end_time > $2
		order by start_time
	`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []RawEventByDay
	for rows.Next() {
		var start, end time.Time
		var att int
		var allDay bool
		if err := rows.Scan(&start, &end, &att, &allDay); err != nil {
			return nil, err
		}
		day := start.Truncate(24 * time.Hour)
		found := false
		for i := range out {
			if out[i].Date.Equal(day) {
				out[i].Events = append(out[i].Events, RawEventRow{Start: start, End: end, AttendeeCount: att, IsAllDay: allDay})
				found = true
				break
			}
		}
		if !found {
			out = append(out, RawEventByDay{
				Date:   day,
				Events: []RawEventRow{{Start: start, End: end, AttendeeCount: att, IsAllDay: allDay}},
			})
		}
	}
	return out, rows.Err()
}

type RawEventByDay struct {
	Date   time.Time
	Events []RawEventRow
}

type RawEventRow struct {
	Start        time.Time
	End          time.Time
	AttendeeCount int
	IsAllDay     bool
}
