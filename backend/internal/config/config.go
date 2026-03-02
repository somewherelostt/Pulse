package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     int
	Env      string
	Frontend string

	SupabaseURL        string
	SupabaseAnonKey    string
	SupabaseServiceKey string
	SupabaseJWTSecret  string
	DatabaseURL        string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string

	GroqAPIKey    string
	GroqModel     string
	GroqBaseURL   string
	CerebrasKey   string
	CerebrasModel string
	CerebrasBaseURL string
	LLMTimeoutSec int
	LLMMaxTokens  int

	SyncIntervalHours        int
	CalendarLookbackDays     int
	CorrelationMinDatapoints int

	// Sleep providers
	OuraToken          string
	FitbitClientID     string
	FitbitClientSecret string
	FitbitRedirectURI  string
}

func Load() (*Config, error) {
	c := &Config{
		Port:     getInt("PORT", 8080),
		Env:      getStr("ENV", "development"),
		Frontend: getStr("FRONTEND_URL", "http://localhost:3000"),

		SupabaseURL:        getStr("SUPABASE_URL", ""),
		SupabaseAnonKey:    getStr("SUPABASE_ANON_KEY", ""),
		SupabaseServiceKey: getStr("SUPABASE_SERVICE_KEY", ""),
		SupabaseJWTSecret:  getStr("SUPABASE_JWT_SECRET", ""),
		DatabaseURL:       getStr("DATABASE_URL", ""),

		GoogleClientID:     getStr("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getStr("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURI:  getStr("GOOGLE_REDIRECT_URI", "http://localhost:8080/api/v1/calendar/callback"),

		GroqAPIKey:      getStr("GROQ_API_KEY", ""),
		GroqModel:       getStr("GROQ_MODEL", "llama-3.3-70b-versatile"),
		GroqBaseURL:     getStr("GROQ_BASE_URL", "https://api.groq.com/openai/v1"),
		CerebrasKey:     getStr("CEREBRAS_API_KEY", ""),
		CerebrasModel:   getStr("CEREBRAS_MODEL", "llama3.1-70b"),
		CerebrasBaseURL: getStr("CEREBRAS_BASE_URL", "https://api.cerebras.ai/v1"),
		LLMTimeoutSec: getInt("LLM_TIMEOUT_SECONDS", 30),
		LLMMaxTokens:  getInt("LLM_MAX_TOKENS", 1024),

		SyncIntervalHours:        getInt("SYNC_INTERVAL_HOURS", 6),
		CalendarLookbackDays:      getInt("CALENDAR_LOOKBACK_DAYS", 30),
		CorrelationMinDatapoints:  getInt("CORRELATION_MIN_DATAPOINTS", 7),

		OuraToken:          getStr("OURA_PERSONAL_TOKEN", ""),
		FitbitClientID:     getStr("FITBIT_CLIENT_ID", ""),
		FitbitClientSecret: getStr("FITBIT_CLIENT_SECRET", ""),
		FitbitRedirectURI:  getStr("FITBIT_REDIRECT_URI", "http://localhost:8080/api/v1/sleep/fitbit/callback"),
	}
	return c, nil
}

func getStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return def
}

func LoadDotenv() error {
	return godotenv.Load()
}
