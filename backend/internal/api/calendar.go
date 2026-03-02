package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"pulse-api/internal/collectors/google"
	"pulse-api/internal/db"
	"pulse-api/internal/middleware"
	"pulse-api/internal/pipeline"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
)

type CalendarHandler struct {
	FrontendURL  string
	Pool         *pgxpool.Pool
	OAuthConfig  *oauth2.Config
	LookbackDays int
	WorkStartH   int
	WorkEndH     int
}

func NewCalendarHandler(frontendURL string, pool *pgxpool.Pool, oauthConfig *oauth2.Config, lookbackDays, workStartH, workEndH int) *CalendarHandler {
	return &CalendarHandler{
		FrontendURL:  frontendURL,
		Pool:         pool,
		OAuthConfig:  oauthConfig,
		LookbackDays: lookbackDays,
		WorkStartH:   workStartH,
		WorkEndH:     workEndH,
	}
}

func (h *CalendarHandler) Connect(w http.ResponseWriter, r *http.Request) {
	supabaseUID := middleware.GetUserID(r.Context())
	if supabaseUID == "" {
		writeErr(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}
	url := google.AuthCodeURL(h.OAuthConfig, supabaseUID)
	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *CalendarHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state") // state = supabaseUID
	if code == "" || state == "" {
		http.Redirect(w, r, h.FrontendURL+"/onboarding?error=oauth_failed", http.StatusFound)
		return
	}
	token, err := google.Exchange(r.Context(), h.OAuthConfig, code)
	if err != nil {
		slog.Error("google oauth exchange failed", "err", err)
		http.Redirect(w, r, h.FrontendURL+"/onboarding?error=oauth_failed", http.StatusFound)
		return
	}
	// Look up user by supabase UID (state)
	u, err := db.UserBySupabaseUID(r.Context(), h.Pool, state)
	if err != nil {
		// Create user first
		u, err = db.UpsertUser(r.Context(), h.Pool, state, "UTC", 9, 18)
		if err != nil {
			slog.Error("upsert user failed", "err", err)
			http.Redirect(w, r, h.FrontendURL+"/onboarding?error=user_failed", http.StatusFound)
			return
		}
	}
	// Store token
	var expiry *time.Time
	if !token.Expiry.IsZero() {
		t := token.Expiry
		expiry = &t
	}
	if err := db.SetOAuthToken(r.Context(), h.Pool, u.ID, "google", token.AccessToken, token.RefreshToken, expiry); err != nil {
		slog.Error("store oauth token failed", "err", err)
	}
	// Mark calendar consent
	_ = db.UpdateUserOnboarding(r.Context(), h.Pool, u.ID, true)
	
	// Trigger sync immediately in background
	go h.syncCalendar(u, token)
	
	http.Redirect(w, r, h.FrontendURL+"/onboarding?step=2&syncing=true", http.StatusFound)
}

func (h *CalendarHandler) Status(w http.ResponseWriter, r *http.Request) {
	supabaseUID := middleware.GetUserID(r.Context())
	u, err := db.UserBySupabaseUID(r.Context(), h.Pool, supabaseUID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "user not found", "USER_NOT_FOUND")
		return
	}
	syncedAt, eventsFetched, status, err := db.GetSyncLog(r.Context(), h.Pool, u.ID, "google")
	connected := err == nil && syncedAt != nil
	if err != nil {
		status = "never"
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"connected":      connected,
		"last_synced_at": syncedAt,
		"events_fetched": eventsFetched,
		"status":         status,
	})
}

func (h *CalendarHandler) Sync(w http.ResponseWriter, r *http.Request) {
	supabaseUID := middleware.GetUserID(r.Context())
	u, err := db.UserBySupabaseUID(r.Context(), h.Pool, supabaseUID)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "user not found — complete onboarding first", "USER_NOT_FOUND")
		return
	}
	// Get token
	access, refresh, expiry, err := db.GetOAuthToken(r.Context(), h.Pool, u.ID, "google")
	if err != nil || access == "" {
		writeErr(w, http.StatusBadRequest, "calendar not connected", "NOT_CONNECTED")
		return
	}
	token := &oauth2.Token{AccessToken: access, RefreshToken: refresh}
	if expiry != nil {
		token.Expiry = *expiry
	}
	// Refresh if needed
	newTok, err := google.RefreshToken(r.Context(), h.OAuthConfig, token)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "token refresh failed", "TOKEN_REFRESH_FAILED")
		return
	}
	// Update stored token if refreshed
	if newTok.AccessToken != access {
		var newExp *time.Time
		if !newTok.Expiry.IsZero() {
			t := newTok.Expiry
			newExp = &t
		}
		_ = db.SetOAuthToken(r.Context(), h.Pool, u.ID, "google", newTok.AccessToken, newTok.RefreshToken, newExp)
	}
	// Fetch calendar events
	svc, err := google.NewCalendarService(r.Context(), newTok, h.OAuthConfig)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "calendar service failed", "CALENDAR_ERROR")
		return
	}
	loc, _ := time.LoadLocation(u.Timezone)
	if loc == nil {
		loc = time.UTC
	}
	events, err := google.FetchEvents(r.Context(), svc, h.LookbackDays, loc)
	if err != nil {
		_ = db.UpsertSyncLog(r.Context(), h.Pool, u.ID, "google", "error", 0, ptr(err.Error()))
		writeErr(w, http.StatusInternalServerError, "calendar fetch failed", "FETCH_FAILED")
		return
	}
	// Store raw events
	var rawEvents []db.RawCalendarEvent
	for _, e := range events {
		titleHash := pipeline.HashTitle(e.Title)
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
	if err := db.UpsertRawEvents(r.Context(), h.Pool, u.ID, rawEvents); err != nil {
		slog.Error("upsert raw events failed", "err", err)
	}
	// Extract features
	_ = pipeline.ExtractAndStoreFeatures(r.Context(), h.Pool, u.ID, u.WorkStartHour, u.WorkEndHour)
	count := len(events)
	_ = db.UpsertSyncLog(r.Context(), h.Pool, u.ID, "google", "success", count, nil)
	writeJSON(w, http.StatusOK, map[string]interface{}{"status": "success", "events_fetched": count})
}

func ptr(s string) *string { return &s }

func (h *CalendarHandler) syncCalendar(u *db.User, token *oauth2.Token) {
	ctx := context.Background()
	
	// Refresh if needed
	newTok, err := google.RefreshToken(ctx, h.OAuthConfig, token)
	if err != nil {
		slog.Error("background sync: token refresh failed", "user_id", u.ID, "err", err)
		_ = db.UpsertSyncLog(ctx, h.Pool, u.ID, "google", "error", 0, ptr("token refresh failed"))
		return
	}
	
	// Update stored token if refreshed
	if newTok.AccessToken != token.AccessToken {
		var newExp *time.Time
		if !newTok.Expiry.IsZero() {
			t := newTok.Expiry
			newExp = &t
		}
		_ = db.SetOAuthToken(ctx, h.Pool, u.ID, "google", newTok.AccessToken, newTok.RefreshToken, newExp)
	}
	
	// Fetch calendar events
	svc, err := google.NewCalendarService(ctx, newTok, h.OAuthConfig)
	if err != nil {
		slog.Error("background sync: calendar service failed", "user_id", u.ID, "err", err)
		_ = db.UpsertSyncLog(ctx, h.Pool, u.ID, "google", "error", 0, ptr("calendar service failed"))
		return
	}
	
	loc, _ := time.LoadLocation(u.Timezone)
	if loc == nil {
		loc = time.UTC
	}
	events, err := google.FetchEvents(ctx, svc, h.LookbackDays, loc)
	if err != nil {
		slog.Error("background sync: fetch events failed", "user_id", u.ID, "err", err)
		_ = db.UpsertSyncLog(ctx, h.Pool, u.ID, "google", "error", 0, ptr(err.Error()))
		return
	}
	
	// Store raw events
	var rawEvents []db.RawCalendarEvent
	for _, e := range events {
		titleHash := pipeline.HashTitle(e.Title)
		durationMins := e.End.Sub(e.Start).Minutes()
		afterHours := e.Start.Hour() < u.WorkStartHour || e.Start.Hour() >= u.WorkEndHour
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
	
	if err := db.UpsertRawEvents(ctx, h.Pool, u.ID, rawEvents); err != nil {
		slog.Error("background sync: upsert raw events failed", "user_id", u.ID, "err", err)
	}
	
	// Extract features
	_ = pipeline.ExtractAndStoreFeatures(ctx, h.Pool, u.ID, u.WorkStartHour, u.WorkEndHour)
	
	count := len(events)
	_ = db.UpsertSyncLog(ctx, h.Pool, u.ID, "google", "success", count, nil)
	slog.Info("background sync completed", "user_id", u.ID, "events_fetched", count)
}
