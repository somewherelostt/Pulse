package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"

type SupabaseClaims struct {
	Sub string `json:"sub"`
	jwt.RegisteredClaims
}

func Auth(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				writeJSONError(w, http.StatusUnauthorized, "missing authorization header", "UNAUTHORIZED")
				return
			}
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				writeJSONError(w, http.StatusUnauthorized, "invalid authorization header", "UNAUTHORIZED")
				return
			}
			tokenStr := parts[1]
			token, err := jwt.ParseWithClaims(tokenStr, &SupabaseClaims{}, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil {
				slog.Debug("jwt parse failed", "err", err)
				writeJSONError(w, http.StatusUnauthorized, "invalid token", "UNAUTHORIZED")
				return
			}
			claims, ok := token.Claims.(*SupabaseClaims)
			if !ok || !token.Valid {
				writeJSONError(w, http.StatusUnauthorized, "invalid token", "UNAUTHORIZED")
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, claims.Sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) string {
	v, _ := ctx.Value(UserIDKey).(string)
	return v
}

func writeJSONError(w http.ResponseWriter, status int, msg, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":"` + escapeJSON(msg) + `","code":"` + code + `"}`))
}

func escapeJSON(s string) string {
	var b []byte
	for _, r := range s {
		switch r {
		case '"', '\\':
			b = append(b, '\\', byte(r))
		default:
			b = append(b, byte(r))
		}
	}
	return string(b)
}
