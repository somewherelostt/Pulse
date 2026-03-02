package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MoodLogRow struct {
	ID       string
	Date     time.Time
	Score    int
	Energy   *int
	Anxiety  *int
	Note     *string
	Tags     []string
	LoggedAt time.Time
}

func CreateOrUpdateMood(ctx context.Context, pool *pgxpool.Pool, userID string, date time.Time, score int, energy, anxiety *int, note *string, tags []string) (*MoodLogRow, error) {
	_, err := pool.Exec(ctx, `
		insert into public.mood_logs (user_id, date, score, energy, anxiety, note, tags)
		values ($1::uuid, $2, $3, $4, $5, $6, $7)
		on conflict (user_id, date) do update set
			score = excluded.score, energy = excluded.energy, anxiety = excluded.anxiety,
			note = excluded.note, tags = excluded.tags, logged_at = now()
	`, userID, date, score, energy, anxiety, note, tags)
	if err != nil {
		return nil, err
	}
	return GetMoodByDate(ctx, pool, userID, date)
}

func GetMoodByDate(ctx context.Context, pool *pgxpool.Pool, userID string, date time.Time) (*MoodLogRow, error) {
	var m MoodLogRow
	err := pool.QueryRow(ctx, `
		select id::text, date, score, energy, anxiety, note, tags, logged_at
		from public.mood_logs where user_id = $1::uuid and date = $2
	`, userID, date).Scan(&m.ID, &m.Date, &m.Score, &m.Energy, &m.Anxiety, &m.Note, &m.Tags, &m.LoggedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func GetMoodRange(ctx context.Context, pool *pgxpool.Pool, userID string, from, to time.Time) ([]MoodLogRow, error) {
	rows, err := pool.Query(ctx, `
		select id::text, date, score, energy, anxiety, note, tags, logged_at
		from public.mood_logs where user_id = $1::uuid and date >= $2 and date <= $3 order by date
	`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []MoodLogRow
	for rows.Next() {
		var m MoodLogRow
		if err := rows.Scan(&m.ID, &m.Date, &m.Score, &m.Energy, &m.Anxiety, &m.Note, &m.Tags, &m.LoggedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
