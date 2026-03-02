package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"pulse-api/internal/collectors/sleep"
)

// UpsertSleepSessions inserts or updates sleep sessions for a user.
func UpsertSleepSessions(ctx context.Context, pool *pgxpool.Pool, userID string, sessions []sleep.SleepSession) error {
	for _, s := range sessions {
		_, err := pool.Exec(ctx, `
			insert into public.sleep_sessions
				(user_id, provider, session_date, bedtime, wake_time,
				 total_sleep_mins, rem_mins, deep_mins, light_mins, awake_mins,
				 sleep_score, hrv_rmssd, resting_hr, source_id)
			values ($1::uuid, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			on conflict (user_id, provider, session_date)
			do update set
				bedtime = excluded.bedtime,
				wake_time = excluded.wake_time,
				total_sleep_mins = excluded.total_sleep_mins,
				rem_mins = excluded.rem_mins,
				deep_mins = excluded.deep_mins,
				light_mins = excluded.light_mins,
				awake_mins = excluded.awake_mins,
				sleep_score = excluded.sleep_score,
				hrv_rmssd = excluded.hrv_rmssd,
				resting_hr = excluded.resting_hr,
				source_id = excluded.source_id,
				fetched_at = now()
		`,
			userID, s.Provider, s.SessionDate, s.Bedtime, s.WakeTime,
			s.TotalMins, s.REMMins, s.DeepMins, s.LightMins, s.AwakeMins,
			s.SleepScore, s.HRV, s.RestingHR, s.SourceID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// SleepSessionRow is a raw sleep session from the database.
type SleepSessionRow struct {
	SessionDate   time.Time
	Provider      string
	Bedtime       *time.Time
	WakeTime      *time.Time
	TotalMins     *int
	REMMins       *int
	DeepMins      *int
	LightMins     *int
	AwakeMins     *int
	SleepScore    *int
	HRV           *float64
	RestingHR     *float64
}

// GetSleepSessions returns sleep sessions for a user in a date range.
func GetSleepSessions(ctx context.Context, pool *pgxpool.Pool, userID string, from, to time.Time) ([]SleepSessionRow, error) {
	rows, err := pool.Query(ctx, `
		select session_date, provider, bedtime, wake_time,
		       total_sleep_mins, rem_mins, deep_mins, light_mins, awake_mins,
		       sleep_score, hrv_rmssd, resting_hr
		from public.sleep_sessions
		where user_id = $1::uuid and session_date >= $2 and session_date <= $3
		order by session_date
	`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SleepSessionRow
	for rows.Next() {
		var r SleepSessionRow
		if err := rows.Scan(
			&r.SessionDate, &r.Provider, &r.Bedtime, &r.WakeTime,
			&r.TotalMins, &r.REMMins, &r.DeepMins, &r.LightMins, &r.AwakeMins,
			&r.SleepScore, &r.HRV, &r.RestingHR,
		); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// CircadianFeatureRow is one row from circadian_features.
type CircadianFeatureRow struct {
	Date                 time.Time
	SleepDurationMins    *float64
	SleepEfficiencyPct   *float64
	SleepDebtMins        *float64
	MidSleepHour         *float64
	RhythmConsistencyPct *float64
	SocialJetlagMins     *float64
	REMPct               *float64
	DeepPct              *float64
	HRV                  *float64
	RestingHR            *float64
	SleepScore           *float64
	LightHygieneScore    *float64
}

// UpsertCircadianFeatures stores extracted circadian features.
func UpsertCircadianFeatures(ctx context.Context, pool *pgxpool.Pool, userID string, rows []CircadianFeatureRow) error {
	for _, r := range rows {
		_, err := pool.Exec(ctx, `
			insert into public.circadian_features
				(user_id, date, sleep_duration_mins, sleep_efficiency_pct, sleep_debt_mins,
				 mid_sleep_hour, rhythm_consistency_pct, social_jetlag_mins,
				 rem_pct, deep_pct, hrv_rmssd, resting_hr, sleep_score, light_hygiene_score)
			values ($1::uuid,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
			on conflict (user_id, date) do update set
				sleep_duration_mins    = excluded.sleep_duration_mins,
				sleep_efficiency_pct   = excluded.sleep_efficiency_pct,
				sleep_debt_mins        = excluded.sleep_debt_mins,
				mid_sleep_hour         = excluded.mid_sleep_hour,
				rhythm_consistency_pct = excluded.rhythm_consistency_pct,
				social_jetlag_mins     = excluded.social_jetlag_mins,
				rem_pct                = excluded.rem_pct,
				deep_pct               = excluded.deep_pct,
				hrv_rmssd              = excluded.hrv_rmssd,
				resting_hr             = excluded.resting_hr,
				sleep_score            = excluded.sleep_score,
				light_hygiene_score    = excluded.light_hygiene_score,
				computed_at            = now()
		`,
			userID, r.Date, r.SleepDurationMins, r.SleepEfficiencyPct, r.SleepDebtMins,
			r.MidSleepHour, r.RhythmConsistencyPct, r.SocialJetlagMins,
			r.REMPct, r.DeepPct, r.HRV, r.RestingHR, r.SleepScore, r.LightHygieneScore,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetCircadianFeatures returns circadian feature rows for a user in a date range.
func GetCircadianFeatures(ctx context.Context, pool *pgxpool.Pool, userID string, from, to time.Time) ([]CircadianFeatureRow, error) {
	rows, err := pool.Query(ctx, `
		select date, sleep_duration_mins, sleep_efficiency_pct, sleep_debt_mins,
		       mid_sleep_hour, rhythm_consistency_pct, social_jetlag_mins,
		       rem_pct, deep_pct, hrv_rmssd, resting_hr, sleep_score, light_hygiene_score
		from public.circadian_features
		where user_id = $1::uuid and date >= $2 and date <= $3
		order by date
	`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CircadianFeatureRow
	for rows.Next() {
		var r CircadianFeatureRow
		if err := rows.Scan(
			&r.Date, &r.SleepDurationMins, &r.SleepEfficiencyPct, &r.SleepDebtMins,
			&r.MidSleepHour, &r.RhythmConsistencyPct, &r.SocialJetlagMins,
			&r.REMPct, &r.DeepPct, &r.HRV, &r.RestingHR, &r.SleepScore, &r.LightHygieneScore,
		); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// UpsertCircadianInsight saves or replaces a weekly circadian narrative.
func UpsertCircadianInsight(ctx context.Context, pool *pgxpool.Pool, userID string, weekStart time.Time, narrative string, interventions []byte, modelUsed string) error {
	_, err := pool.Exec(ctx, `
		insert into public.circadian_insights (user_id, week_start, narrative, interventions, model_used)
		values ($1::uuid, $2, $3, $4, $5)
		on conflict (user_id, week_start) do update set
			narrative     = excluded.narrative,
			interventions = excluded.interventions,
			model_used    = excluded.model_used,
			created_at    = now()
	`, userID, weekStart, narrative, interventions, modelUsed)
	return err
}

// GetLatestCircadianInsight retrieves the most recent circadian insight for a user.
func GetLatestCircadianInsight(ctx context.Context, pool *pgxpool.Pool, userID string) (narrative string, interventions []byte, modelUsed string, weekStart time.Time, err error) {
	err = pool.QueryRow(ctx, `
		select narrative, interventions, model_used, week_start
		from public.circadian_insights
		where user_id = $1::uuid
		order by week_start desc limit 1
	`, userID).Scan(&narrative, &interventions, &modelUsed, &weekStart)
	return
}
