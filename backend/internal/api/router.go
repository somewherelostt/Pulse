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
	userHandler *UserHandler,
	calendarHandler *CalendarHandler,
	moodHandler *MoodHandler,
	dashboardHandler *DashboardHandler,
	insightsHandler *InsightsHandler,
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

	r.Route("/api/v1", func(r chi.Router) {
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
		})
	})
	return r
}
