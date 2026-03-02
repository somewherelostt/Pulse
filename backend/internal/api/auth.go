package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthHandler struct {
	Pool        *pgxpool.Pool
	SupabaseURL string
}

func NewAuthHandler(pool *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{
		Pool:        pool,
		SupabaseURL: os.Getenv("SUPABASE_URL"),
	}
}

type SupabaseAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID string `json:"id"`
	} `json:"user"`
}

func (h *AuthHandler) Anonymous(w http.ResponseWriter, r *http.Request) {
	// Proxy anonymous sign-in to Supabase
	supabaseURL := os.Getenv("SUPABASE_URL")
	anonKey := os.Getenv("SUPABASE_ANON_KEY")

	if supabaseURL == "" || anonKey == "" {
		slog.Error("missing supabase config", "url", supabaseURL, "key_set", anonKey != "")
		writeErr(w, http.StatusInternalServerError, "supabase not configured", "CONFIG_ERROR")
		return
	}

	// Supabase anonymous sign-in via signInAnonymously endpoint
	authURL := supabaseURL + "/auth/v1/signup"
	
	// Empty JSON body for anonymous signup
	payload := []byte(`{"data":{}}`)
	req, err := http.NewRequestWithContext(r.Context(), "POST", authURL, 	bytes.NewBuffer(payload))
	if err != nil {
		slog.Error("failed to create auth request", "error", err)
		writeErr(w, http.StatusInternalServerError, "auth request failed", "INTERNAL_ERROR")
		return
	}

	req.Header.Set("apikey", anonKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+anonKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("failed to call supabase auth", "error", err, "url", authURL)
		writeErr(w, http.StatusInternalServerError, "auth call failed", "INTERNAL_ERROR")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		slog.Error("supabase auth failed", "status", resp.StatusCode)
		writeErr(w, http.StatusInternalServerError, "auth failed", "AUTH_ERROR")
		return
	}

	var authResp SupabaseAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		slog.Error("failed to decode auth response", "error", err)
		writeErr(w, http.StatusInternalServerError, "auth parse failed", "INTERNAL_ERROR")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  authResp.AccessToken,
		"refresh_token": authResp.RefreshToken,
		"expires_in":    authResp.ExpiresIn,
		"user_id":       authResp.User.ID,
	})
}
