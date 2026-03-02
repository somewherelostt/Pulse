package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	authMw func(http.Handler) http.Handler,
	corsMw func(http.Handler) http.Handler,
	logMw func(http.Handler) http.Handler,
	recMw func(http.Handler) http.Handler,
	authHandler *AuthHandler,
	userHandler *UserHandler,
	calendarHandler *CalendarHandler,
	moodHandler *MoodHandler,
	dashboardHandler *DashboardHandler,
	insightsHandler *InsightsHandler,
	sleepHandler *SleepHandler,
	circadianHandler *CircadianHandler,
	constellationHandler *ConstellationHandler,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chimw.RealIP)
	r.Use(chimw.RequestID)
	r.Use(recMw)
	r.Use(logMw)
	r.Use(corsMw)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Serve web pages (for standalone Go demo)
	r.Get("/", serveFile("web/dashboard.html"))
	r.Get("/connect", serveFile("web/connect.html"))
	r.Get("/log", serveFile("web/log.html"))
	r.Get("/circadian", serveFile("web/circadian.html"))
	r.Get("/constellation", serveFile("web/constellation.html"))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/anonymous", authHandler.Anonymous)
		r.Get("/calendar/callback", calendarHandler.Callback)

		r.Group(func(r chi.Router) {
			r.Use(authMw)
			r.Post("/users/me", userHandler.UpsertMe)
			r.Get("/users/me", userHandler.GetMe)
			r.Get("/calendar/connect", calendarHandler.Connect)
			r.Get("/calendar/status", calendarHandler.Status)
			r.Post("/calendar/sync", calendarHandler.Sync)
			r.Post("/mood", moodHandler.CreateOrUpdate)
			r.Get("/mood/today", moodHandler.GetToday)
			r.Get("/mood/range", moodHandler.GetRange)
			r.Get("/dashboard", dashboardHandler.Get)
			r.Get("/insights/latest", insightsHandler.Latest)
			r.Post("/insights/generate", insightsHandler.Generate)

			// Sleep endpoints
			r.Post("/sleep/manual", sleepHandler.LogManual)
			r.Get("/sleep/range", sleepHandler.GetRange)

			// Circadian endpoints
			r.Get("/circadian/dashboard", circadianHandler.Dashboard)
			r.Post("/circadian/extract", circadianHandler.ExtractFeatures)
			r.Post("/circadian/narrative", circadianHandler.GenerateNarrative)

			// Constellation peer matching + WebRTC signaling
			r.Post("/constellation/join", constellationHandler.Join)
			r.Post("/constellation/leave", constellationHandler.Leave)
			r.Get("/constellation/safety", constellationHandler.Safety)
			r.Get("/constellation/match", constellationHandler.Match)
			r.Post("/constellation/session/start", constellationHandler.SessionStart)
			r.Post("/constellation/session/{id}/end", constellationHandler.SessionEnd)
			r.Post("/constellation/session/{id}/rate", constellationHandler.SessionRate)
			r.Get("/constellation/signal/{room_id}", constellationHandler.Signal)
		})
	})
	return r
}

func serveFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}
