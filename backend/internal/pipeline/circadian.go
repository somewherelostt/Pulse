package pipeline

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"pulse-api/internal/db"
	circadianfeat "pulse-api/internal/features/circadian"
)

// ExtractAndStoreCircadianFeatures fetches sleep sessions for a user,
// computes circadian features, and persists them.
func ExtractAndStoreCircadianFeatures(ctx context.Context, pool *pgxpool.Pool, userID string) error {
	now := time.Now()
	from := now.AddDate(0, 0, -37) // 30 days + 7 for rolling window

	sessions, err := db.GetSleepSessions(ctx, pool, userID, from, now)
	if err != nil || len(sessions) == 0 {
		return err
	}

	// Build generic structs required by circadian.BuildSleepDays
	rawSessions := make([]struct {
		Date       time.Time
		Bedtime    *time.Time
		WakeTime   *time.Time
		TotalMins  *int
		REMMins    *int
		DeepMins   *int
		LightMins  *int
		AwakeMins  *int
		SleepScore *int
		HRV        *float64
		RestingHR  *float64
	}, len(sessions))
	for i, s := range sessions {
		rawSessions[i].Date = s.SessionDate
		rawSessions[i].Bedtime = s.Bedtime
		rawSessions[i].WakeTime = s.WakeTime
		rawSessions[i].TotalMins = s.TotalMins
		rawSessions[i].REMMins = s.REMMins
		rawSessions[i].DeepMins = s.DeepMins
		rawSessions[i].LightMins = s.LightMins
		rawSessions[i].AwakeMins = s.AwakeMins
		rawSessions[i].SleepScore = s.SleepScore
		rawSessions[i].HRV = s.HRV
		rawSessions[i].RestingHR = s.RestingHR
	}

	sleepDays := circadianfeat.BuildSleepDays(rawSessions)
	extracted := circadianfeat.ExtractAll(sleepDays)

	dbRows := make([]db.CircadianFeatureRow, len(extracted))
	for i, e := range extracted {
		dbRows[i] = db.CircadianFeatureRow{
			Date:                 e.Date,
			SleepDurationMins:    e.SleepDurationMins,
			SleepEfficiencyPct:   e.SleepEfficiencyPct,
			SleepDebtMins:        e.SleepDebtMins,
			MidSleepHour:         e.MidSleepHour,
			RhythmConsistencyPct: e.RhythmConsistencyPct,
			SocialJetlagMins:     e.SocialJetlagMins,
			REMPct:               e.REMPct,
			DeepPct:              e.DeepPct,
			HRV:                  e.HRV,
			RestingHR:            e.RestingHR,
			SleepScore:           e.SleepScore,
			LightHygieneScore:    e.LightHygieneScore,
		}
	}
	return db.UpsertCircadianFeatures(ctx, pool, userID, dbRows)
}
