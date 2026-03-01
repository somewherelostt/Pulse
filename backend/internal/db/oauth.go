package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetOAuthToken(ctx context.Context, pool *pgxpool.Pool, userID, provider string) (accessToken, refreshToken string, expiry *time.Time, err error) {
	err = pool.QueryRow(ctx, `
		select access_token, refresh_token, token_expiry from public.oauth_tokens
		where user_id = $1::uuid and provider = $2
	`, userID, provider).Scan(&accessToken, &refreshToken, &expiry)
	return
}

func SetOAuthToken(ctx context.Context, pool *pgxpool.Pool, userID, provider, accessToken, refreshToken string, expiry *time.Time) error {
	_, err := pool.Exec(ctx, `
		insert into public.oauth_tokens (user_id, provider, access_token, refresh_token, token_expiry, updated_at)
		values ($1::uuid, $2, $3, $4, $5, now())
		on conflict (user_id, provider) do update set
			access_token = excluded.access_token,
			refresh_token = coalesce(excluded.refresh_token, oauth_tokens.refresh_token),
			token_expiry = excluded.token_expiry,
			updated_at = now()
	`, userID, provider, accessToken, refreshToken, expiry)
	return err
}
