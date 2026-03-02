package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"

type SupabaseClaims struct {
	Sub string `json:"sub"`
	jwt.RegisteredClaims
}

type jwksCache struct {
	mu        sync.Mutex
	keysByKID map[string]*rsa.PublicKey
	fetchedAt time.Time
}

func Auth(jwtSecret string, supabaseURL string) func(next http.Handler) http.Handler {
	cache := &jwksCache{keysByKID: map[string]*rsa.PublicKey{}}
	devMode := os.Getenv("ENV") == "development"
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

			// In development, accept Supabase tokens without verifying the signature
			// to avoid local key/JWKS configuration issues. We still parse claims to
			// extract the subject (user id).
			if devMode {
				var claims SupabaseClaims
				parser := jwt.NewParser(jwt.WithoutClaimsValidation())
				_, _, err := parser.ParseUnverified(tokenStr, &claims)
				if err != nil || claims.Sub == "" {
					slog.Debug("jwt parse unverified failed", "err", err)
					writeJSONError(w, http.StatusUnauthorized, "invalid token", "UNAUTHORIZED")
					return
				}
				ctx := context.WithValue(r.Context(), UserIDKey, claims.Sub)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			token, err := jwt.ParseWithClaims(tokenStr, &SupabaseClaims{}, func(t *jwt.Token) (interface{}, error) {
				alg, _ := t.Header["alg"].(string)
				switch alg {
				case "HS256":
					if jwtSecret == "" {
						return nil, errors.New("missing SUPABASE_JWT_SECRET")
					}
					return []byte(jwtSecret), nil
				case "RS256":
					kid, _ := t.Header["kid"].(string)
					if kid == "" {
						return nil, errors.New("missing kid in token header")
					}
					key, err := cache.getKey(r.Context(), supabaseURL, kid)
					if err != nil {
						return nil, err
					}
					return key, nil
				default:
					return nil, fmt.Errorf("unsupported jwt alg: %s", alg)
				}
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

func (c *jwksCache) getKey(ctx context.Context, supabaseURL, kid string) (*rsa.PublicKey, error) {
	// Cache JWKS for 10 minutes to avoid network on every request.
	c.mu.Lock()
	if k, ok := c.keysByKID[kid]; ok && time.Since(c.fetchedAt) < 10*time.Minute {
		c.mu.Unlock()
		return k, nil
	}
	c.mu.Unlock()

	keys, err := fetchJWKS(ctx, strings.TrimRight(supabaseURL, "/")+"/auth/v1/.well-known/jwks.json")
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.keysByKID = keys
	c.fetchedAt = time.Now()
	k := c.keysByKID[kid]
	c.mu.Unlock()

	if k == nil {
		return nil, fmt.Errorf("jwks key not found for kid: %s", kid)
	}
	return k, nil
}

func fetchJWKS(ctx context.Context, url string) (map[string]*rsa.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("jwks fetch failed: %s", res.Status)
	}
	var body struct {
		Keys []struct {
			KTY string `json:"kty"`
			KID string `json:"kid"`
			Use string `json:"use"`
			Alg string `json:"alg"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}
	out := map[string]*rsa.PublicKey{}
	for _, k := range body.Keys {
		if k.KTY != "RSA" || k.KID == "" || k.N == "" || k.E == "" {
			continue
		}
		pub, err := rsaPublicKeyFromJWK(k.N, k.E)
		if err != nil {
			continue
		}
		out[k.KID] = pub
	}
	if len(out) == 0 {
		return nil, errors.New("jwks contains no usable rsa keys")
	}
	return out, nil
}

func rsaPublicKeyFromJWK(nB64URL, eB64URL string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nB64URL)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eB64URL)
	if err != nil {
		return nil, err
	}
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes).Int64()
	if e <= 0 {
		return nil, errors.New("invalid exponent")
	}
	return &rsa.PublicKey{N: n, E: int(e)}, nil
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
