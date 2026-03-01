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
	"pulse-api/internal/config"
	"pulse-api/internal/middleware"
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

	authMw := middleware.Auth(cfg.SupabaseJWTSecret)
	corsMw := middleware.CORS(cfg.Frontend)
	logMw := middleware.Logger
	recMw := middleware.Recover

	userHandler := api.NewUserHandler()
	calendarHandler := api.NewCalendarHandler(cfg.Frontend)
	moodHandler := api.NewMoodHandler()
	dashboardHandler := api.NewDashboardHandler()
	insightsHandler := api.NewInsightsHandler()

	router := api.NewRouter(
		authMw, corsMw, logMw, recMw,
		userHandler, calendarHandler, moodHandler,
		dashboardHandler, insightsHandler,
	)

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}
