package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/FatihTirek/football-league/internal/service"
)

type leagueHandler struct {
	svc service.LeagueService
}

func newLeagueHandler(svc service.LeagueService) *leagueHandler {
	return &leagueHandler{svc: svc}
}

// GetTeams godoc
// GET /api/v1/teams
// Returns all four teams with their strength metrics.
func (h *leagueHandler) GetTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.svc.GetTeams(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, teams)
}

// GetStandings godoc
// GET /api/v1/standings
// Returns the current league table sorted by Points, Goal Difference, Goals For.
func (h *leagueHandler) GetStandings(w http.ResponseWriter, r *http.Request) {
	standings, err := h.svc.GetStandings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, standings)
}

// GetAllMatches godoc
// GET /api/v1/matches
// Returns every fixture in the league, ordered by week.
func (h *leagueHandler) GetAllMatches(w http.ResponseWriter, r *http.Request) {
	matches, err := h.svc.GetAllMatches(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, matches)
}

// GetWeekMatches godoc
// GET /api/v1/weeks/{week}/matches
// Returns the two fixtures for the given week number (1–6).
func (h *leagueHandler) GetWeekMatches(w http.ResponseWriter, r *http.Request) {
	week, err := strconv.Atoi(chi.URLParam(r, "week"))
	if err != nil || week < 1 || week > 6 {
		writeError(w, http.StatusBadRequest, "week must be a number between 1 and 6")
		return
	}
	matches, err := h.svc.GetWeekMatches(r.Context(), week)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, matches)
}

// PlayNextWeek godoc
// POST /api/v1/weeks/next/play
// Simulates the next unplayed week and returns its match results.
// Returns 409 Conflict when all 6 weeks have already been played.
func (h *leagueHandler) PlayNextWeek(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.PlayNextWeek(r.Context())
	if errors.Is(err, service.ErrLeagueComplete) {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// PlayAll godoc
// POST /api/v1/play-all
// Simulates all remaining weeks and returns an array of week results.
// Returns 409 Conflict when the league was already complete before the call.
func (h *leagueHandler) PlayAll(w http.ResponseWriter, r *http.Request) {
	results, err := h.svc.PlayAll(r.Context())
	if errors.Is(err, service.ErrLeagueComplete) {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, results)
}

// editResultRequest is the expected JSON body for EditMatchResult.
type editResultRequest struct {
	HomeScore int `json:"home_score"`
	AwayScore int `json:"away_score"`
}

// EditMatchResult godoc
// PUT /api/v1/matches/{id}/result
// Body: { "home_score": 2, "away_score": 1 }
// Overwrites a match score and returns the fully recalculated standings table.
// Returns 404 when the match ID doesn't exist.
func (h *leagueHandler) EditMatchResult(w http.ResponseWriter, r *http.Request) {
	matchID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || matchID < 1 {
		writeError(w, http.StatusBadRequest, "invalid match id")
		return
	}

	var req editResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}
	if req.HomeScore < 0 || req.AwayScore < 0 {
		writeError(w, http.StatusBadRequest, "scores cannot be negative")
		return
	}

	standings, err := h.svc.EditMatchResult(r.Context(), matchID, req.HomeScore, req.AwayScore)
	if errors.Is(err, service.ErrMatchNotFound) {
		writeError(w, http.StatusNotFound, "match not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, standings)
}