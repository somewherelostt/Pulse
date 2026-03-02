package pipeline

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"pulse-api/internal/db"
	"pulse-api/internal/features/calendar"
)

func SyncUserCalendar(ctx context.Context, pool *pgxpool.Pool, userID, supabaseUID string, lookbackDays, workStartHour, workEndHour int) (eventsFetched int, err error) {
	// Get OAuth token
	access, refresh, expiry, err := db.GetOAuthToken(ctx, pool, userID, "google")
	if err != nil || access == "" {
		db.UpsertSyncLog(ctx, pool, userID, "google", "failed", 0, ptr("calendar not connected"))
		return 0, err
	}
	// Build token and get calendar service (we need config for oauth2.Config - pass from caller)
	// For now this is a placeholder; the actual sync is triggered from the handler with full deps.
	_ = refresh
	_ = expiry
	_ = supabaseUID
	_ = lookbackDays
	_ = workStartHour
	_ = workEndHour
	slog.Info("sync placeholder", "user_id", userID)
	return 0, nil
}

func HashTitle(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func ExtractAndStoreFeatures(ctx context.Context, pool *pgxpool.Pool, userID string, workStartHour, workEndHour int) error {
	to := time.Now()
	from := to.AddDate(0, 0, -31)
	days, err := db.RawEventsByUserAndDateRange(ctx, pool, userID, from, to)
	if err != nil {
		return err
	}
	for _, d := range days {
		eventRows := make([]calendar.EventRow, len(d.Events))
		for i, e := range d.Events {
			eventRows[i] = calendar.EventRow{
				Start:         e.Start,
				End:           e.End,
				AttendeeCount: e.AttendeeCount,
				IsAllDay:      e.IsAllDay,
			}
		}
		df := calendar.ExtractForDay(eventRows, d.Date, workStartHour, workEndHour)
		err := db.UpsertDailyFeatures(ctx, pool, userID, df.Date, db.DailyFeatureRow{
			MeetingDensityPct:  df.MeetingDensityPct,
			MeetingCount:       df.MeetingCount,
			AvgFocusBlockMins:  df.AvgFocusBlockMins,
			MaxFocusBlockMins:  df.MaxFocusBlockMins,
			FragmentationScore: df.FragmentationScore,
			AfterHoursMins:     df.AfterHoursMins,
			WeekendMeetingMins: df.WeekendMeetingMins,
			BackToBackCount:    df.BackToBackCount,
			AvgRecoveryMins:    df.AvgRecoveryMins,
			AttendeeAvg:        df.AttendeeAvg,
			SoloTimePct:        df.SoloTimePct,
			SourceEventCount:   df.SourceEventCount,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func ptr(s string) *string {
	return &s
}
