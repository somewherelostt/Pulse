package pipeline

import (
	"context"
	"log/slog"
	"time"

	"pulse-api/internal/collectors/google"
	"pulse-api/internal/constellation"
	"pulse-api/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
	"golang.org/x/oauth2"
)

// StartScheduler starts the background cron jobs for calendar sync, feature extraction,
// and sleep sync. Returns the cron instance so the caller can stop it on shutdown.
func StartScheduler(ctx context.Context, pool *pgxpool.Pool, oauthConfig *oauth2.Config, lookbackDays int) *cron.Cron {
	c := cron.New()

	// Every 6 hours: sync all calendars
	_, _ = c.AddFunc("0 */6 * * *", func() {
		slog.Info("cron: starting calendar sync for all users")
		start := time.Now()
		syncAllCalendars(ctx, pool, oauthConfig, lookbackDays)
		slog.Info("cron: calendar sync complete", "duration_ms", time.Since(start).Milliseconds())
	})

	// Every midnight: extract calendar features for all users
	_, _ = c.AddFunc("0 0 * * *", func() {
		slog.Info("cron: starting feature extraction for all users")
		start := time.Now()
		extractFeaturesAllUsers(ctx, pool)
		slog.Info("cron: feature extraction complete", "duration_ms", time.Since(start).Milliseconds())
	})

	// 8am daily: extract circadian features for users who have sleep data
	_, _ = c.AddFunc("0 8 * * *", func() {
		slog.Info("cron: starting circadian feature extraction")
		start := time.Now()
		extractCircadianFeaturesAllUsers(ctx, pool)
		slog.Info("cron: circadian extraction complete", "duration_ms", time.Since(start).Milliseconds())
	})

	// Every 15 min: expire inactive peer pool entries (> 30 min without heartbeat)
	_, _ = c.AddFunc("*/15 * * * *", func() {
		if err := constellation.RefreshAvailability(ctx, pool); err != nil {
			slog.Warn("cron: peer pool refresh failed", "err", err)
		}
	})

	c.Start()
	return c
}

func syncAllCalendars(ctx context.Context, pool *pgxpool.Pool, oauthConfig *oauth2.Config, lookbackDays int) {
	// Get all users with google oauth tokens
	rows, err := pool.Query(ctx, `
		select u.id::text, u.timezone, u.work_start_hour, u.work_end_hour,
		       ot.access_token, ot.refresh_token, ot.token_expiry
		from public.users u
		join public.oauth_tokens ot on ot.user_id = u.id and ot.provider = 'google'
	`)
	if err != nil {
		slog.Error("cron: query users failed", "err", err)
		return
	}
	defer rows.Close()

	type userToken struct {
		UserID    string
		Timezone  string
		WorkStart int
		WorkEnd   int
		Access    string
		Refresh   string
		Expiry    *time.Time
	}

	var users []userToken
	for rows.Next() {
		var ut userToken
		if err := rows.Scan(&ut.UserID, &ut.Timezone, &ut.WorkStart, &ut.WorkEnd, &ut.Access, &ut.Refresh, &ut.Expiry); err != nil {
			continue
		}
		users = append(users, ut)
	}
	_ = rows.Err()

	sem := make(chan struct{}, 10)
	for _, ut := range users {
		ut := ut
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			syncOneUser(ctx, pool, oauthConfig, ut.UserID, ut.Timezone, ut.Access, ut.Refresh, ut.Expiry, lookbackDays, ut.WorkStart, ut.WorkEnd)
		}()
	}
	// Drain semaphore to wait for all goroutines
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
}

func syncOneUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	oauthConfig *oauth2.Config,
	userID, timezone, access, refresh string,
	expiry *time.Time,
	lookbackDays, workStart, workEnd int,
) {
	token := &oauth2.Token{AccessToken: access, RefreshToken: refresh}
	if expiry != nil {
		token.Expiry = *expiry
	}

	newTok, err := google.RefreshToken(ctx, oauthConfig, token)
	if err != nil {
		slog.Warn("cron: token refresh failed", "user_id", userID[:8], "err", err)
		_ = db.UpsertSyncLog(ctx, pool, userID, "google", "error", 0, ptrStr(err.Error()))
		return
	}
	if newTok.AccessToken != access {
		var newExp *time.Time
		if !newTok.Expiry.IsZero() {
			t := newTok.Expiry
			newExp = &t
		}
		_ = db.SetOAuthToken(ctx, pool, userID, "google", newTok.AccessToken, newTok.RefreshToken, newExp)
	}

	svc, err := google.NewCalendarService(ctx, newTok, oauthConfig)
	if err != nil {
		slog.Error("cron: calendar service failed", "user_id", userID[:8], "err", err)
		return
	}
	loc, _ := time.LoadLocation(timezone)
	if loc == nil {
		loc = time.UTC
	}
	events, err := google.FetchEvents(ctx, svc, lookbackDays, loc)
	if err != nil {
		_ = db.UpsertSyncLog(ctx, pool, userID, "google", "error", 0, ptrStr(err.Error()))
		slog.Error("cron: fetch events failed", "user_id", userID[:8], "err", err)
		return
	}

	var rawEvents []db.RawCalendarEvent
	for _, e := range events {
		titleHash := HashTitle(e.Title)
		durationMins := e.End.Sub(e.Start).Minutes()
		afterHours := e.Start.Hour() < 8 || e.Start.Hour() >= 18
		isWeekend := e.Start.Weekday() == 0 || e.Start.Weekday() == 6
		rawEvents = append(rawEvents, db.RawCalendarEvent{
			GoogleEventID: e.ID,
			TitleHash:     &titleHash,
			StartTime:     e.Start,
			EndTime:       e.End,
			AttendeeCount: e.Attendees,
			IsAllDay:      e.IsAllDay,
			IsRecurring:   e.IsRecurring,
			IsAfterHours:  afterHours,
			IsWeekend:     isWeekend,
			DurationMins:  durationMins,
		})
	}
	if err := db.UpsertRawEvents(ctx, pool, userID, rawEvents); err != nil {
		slog.Error("cron: upsert events failed", "user_id", userID[:8], "err", err)
	}
	_ = ExtractAndStoreFeatures(ctx, pool, userID, workStart, workEnd)
	_ = db.UpsertSyncLog(ctx, pool, userID, "google", "success", len(events), nil)
	slog.Info("cron: user sync done", "user_id", userID[:8], "events", len(events))
}

func extractFeaturesAllUsers(ctx context.Context, pool *pgxpool.Pool) {
	rows, err := pool.Query(ctx, `select id::text, work_start_hour, work_end_hour from public.users`)
	if err != nil {
		slog.Error("cron: query users failed", "err", err)
		return
	}
	defer rows.Close()
	type u struct {
		ID    string
		Start int
		End   int
	}
	var users []u
	for rows.Next() {
		var x u
		if err := rows.Scan(&x.ID, &x.Start, &x.End); err != nil {
			continue
		}
		users = append(users, x)
	}
	_ = rows.Err()
	for _, usr := range users {
		if err := ExtractAndStoreFeatures(ctx, pool, usr.ID, usr.Start, usr.End); err != nil {
			slog.Error("cron: extract features failed", "user_id", usr.ID[:8], "err", err)
		}
	}
}

func ptrStr(s string) *string { return &s }

func extractCircadianFeaturesAllUsers(ctx context.Context, pool *pgxpool.Pool) {
	rows, err := pool.Query(ctx, `select id::text from public.users`)
	if err != nil {
		slog.Error("cron: query users for circadian", "err", err)
		return
	}
	defer rows.Close()
	var userIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			userIDs = append(userIDs, id)
		}
	}
	_ = rows.Err()

	for _, uid := range userIDs {
		if err := ExtractAndStoreCircadianFeatures(ctx, pool, uid); err != nil {
			slog.Warn("cron: circadian extraction failed", "user_id", uid[:8], "err", err)
		}
	}
}
