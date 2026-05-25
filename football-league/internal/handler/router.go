package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/FatihTirek/football-league/internal/service"
)

// NewRouter builds and returns the fully configured HTTP handler.
// It is the only place in the project where routes are defined.
func NewRouter(ls service.LeagueService, ps service.PredictionService) http.Handler {
	r := chi.NewRouter()

	// Middleware runs on every request in the order it is added.
	r.Use(middleware.RequestID)  // injects X-Request-Id header for distributed tracing
	r.Use(middleware.Logger)     // logs method, path, status code, and latency
	r.Use(middleware.Recoverer)  // catches panics and returns 500 instead of crashing
	r.Use(middleware.Timeout(15 * time.Second)) // sets a hard timeout for all requests to prevent hanging

	lh := newLeagueHandler(ls)
	ph := newPredictionHandler(ps)

	r.Route("/api/v1", func(r chi.Router) {
		// Healthcheck — used by Docker and load balancers to verify the app is alive
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		})

		r.Get("/teams", lh.GetTeams)
		r.Get("/standings", lh.GetStandings)
		r.Get("/matches", lh.GetAllMatches)
		r.Get("/weeks/{week}/matches", lh.GetWeekMatches)
		r.Get("/predictions", ph.GetPredictions)

		r.Post("/weeks/next/play", lh.PlayNextWeek)
		r.Post("/play-all", lh.PlayAll)

		r.Put("/matches/{id}/result", lh.EditMatchResult)
	})

	return r
}