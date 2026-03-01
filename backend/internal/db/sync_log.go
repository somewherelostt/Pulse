package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UpsertSyncLog(ctx context.Context, pool *pgxpool.Pool, userID, provider, status string, eventsFetched int, errMsg *string) error {
	_, err := pool.Exec(ctx, `
		insert into public.sync_log (user_id, provider, synced_at, events_fetched, status, error_message)
		values ($1::uuid, $2, now(), $3, $4, $5)
		on conflict (user_id, provider) do update set
			synced_at = now(), events_fetched = excluded.events_fetched, status = excluded.status, error_message = excluded.error_message
	`, userID, provider, eventsFetched, status, errMsg)
	return err
}

func GetSyncLog(ctx context.Context, pool *pgxpool.Pool, userID, provider string) (syncedAt *time.Time, eventsFetched int, status string, err error) {
	err = pool.QueryRow(ctx, `
		select synced_at, events_fetched, status from public.sync_log where user_id = $1::uuid and provider = $2
	`, userID, provider).Scan(&syncedAt, &eventsFetched, &status)
	return
}
