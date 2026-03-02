package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"pulse-api/internal/api"
	"pulse-api/internal/collectors/google"
	"pulse-api/internal/config"
	"pulse-api/internal/db"
	"pulse-api/internal/llm"
	"pulse-api/internal/middleware"
	"pulse-api/internal/pipeline"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	_ = os.Setenv("ENV", "development")
	if _, err := os.Stat(".env"); err == nil {
		if loadErr := config.LoadDotenv(); loadErr != nil {
			slog.Warn("could not load .env", "err", loadErr)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	// Build Google OAuth config
	oauthConfig := google.Config(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURI)

	// Build LLM client (Groq primary, Cerebras fallback)
	groqClient := llm.NewGroqClient(cfg.GroqAPIKey, cfg.GroqBaseURL, cfg.GroqModel)
	cerebrasClient := llm.NewCerebrasClient(cfg.CerebrasKey, cfg.CerebrasBaseURL, cfg.CerebrasModel)
	var llmClient llm.LLMClient = llm.NewFallbackClient(groqClient, cerebrasClient)

	// Build middleware
	authMw := middleware.Auth(cfg.SupabaseJWTSecret, cfg.SupabaseURL)
	corsMw := middleware.CORS(cfg.Frontend)
	logMw := middleware.Logger
	recMw := middleware.Recover

	// Build handlers
	userHandler := api.NewUserHandler(pool)
	calendarHandler := api.NewCalendarHandler(cfg.Frontend, pool, oauthConfig, cfg.CalendarLookbackDays, 9, 18)
	moodHandler := api.NewMoodHandler(pool)
	dashboardHandler := api.NewDashboardHandler(pool)
	insightsHandler := api.NewInsightsHandler(pool, llmClient)
	sleepHandler := api.NewSleepHandler(pool)
	circadianHandler := api.NewCircadianHandler(pool, llmClient)

	router := api.NewRouter(
		authMw, corsMw, logMw, recMw,
		userHandler, calendarHandler, moodHandler,
		dashboardHandler, insightsHandler,
		sleepHandler, circadianHandler,
	)

	// Start cron scheduler
	scheduler := pipeline.StartScheduler(ctx, pool, oauthConfig, cfg.CalendarLookbackDays)
	defer scheduler.Stop()

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: router,
	}
	go func() {
		slog.Info("server listening", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down")
	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutCtx)
}
