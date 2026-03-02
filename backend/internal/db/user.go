package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID               string
	SupabaseUID      string
	Timezone         string
	OnboardingDone   bool
	ConsentCalendar  bool
	WorkStartHour    int
	WorkEndHour      int
}

func UserBySupabaseUID(ctx context.Context, pool *pgxpool.Pool, supabaseUID string) (*User, error) {
	var u User
	err := pool.QueryRow(ctx, `
		select id::text, supabase_uid::text, timezone, onboarding_done, consent_calendar, work_start_hour, work_end_hour
		from public.users where supabase_uid = $1
	`, supabaseUID).Scan(&u.ID, &u.SupabaseUID, &u.Timezone, &u.OnboardingDone, &u.ConsentCalendar, &u.WorkStartHour, &u.WorkEndHour)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UpsertUser(ctx context.Context, pool *pgxpool.Pool, supabaseUID string, timezone string, workStart, workEnd int) (*User, error) {
	_, err := pool.Exec(ctx, `
		insert into public.users (supabase_uid, timezone, work_start_hour, work_end_hour, updated_at)
		values ($1::uuid, $2, $3, $4, now())
		on conflict (supabase_uid) do update set
			timezone = excluded.timezone,
			work_start_hour = excluded.work_start_hour,
			work_end_hour = excluded.work_end_hour,
			updated_at = now()
	`, supabaseUID, timezone, workStart, workEnd)
	if err != nil {
		return nil, fmt.Errorf("upsert user: %w", err)
	}
	return UserBySupabaseUID(ctx, pool, supabaseUID)
}

func UpdateUserOnboarding(ctx context.Context, pool *pgxpool.Pool, userID string, consentCalendar bool) error {
	_, err := pool.Exec(ctx, `
		update public.users set onboarding_done = true, consent_calendar = $2, updated_at = now() where id = $1::uuid
	`, userID, consentCalendar)
	return err
}
