package handler

import (
	"errors"
	"net/http"

	"github.com/FatihTirek/football-league/internal/service"
)

type predictionHandler struct {
	svc service.PredictionService
}

func newPredictionHandler(svc service.PredictionService) *predictionHandler {
	return &predictionHandler{svc: svc}
}

// GetPredictions godoc
// GET /api/v1/predictions
// Runs 10,000 Monte Carlo simulations and returns each team's championship probability.
// Returns 400 Bad Request if fewer than 4 weeks have been played.
func (h *predictionHandler) GetPredictions(w http.ResponseWriter, r *http.Request) {
	predictions, err := h.svc.GetPredictions(r.Context())
	if errors.Is(err, service.ErrPredictionUnavailable) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, predictions)
}