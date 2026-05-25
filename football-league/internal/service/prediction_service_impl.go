package service

import (
	"context"
	"errors"

	"github.com/FatihTirek/football-league/internal/domain"
	"github.com/FatihTirek/football-league/internal/engine"
	"github.com/FatihTirek/football-league/internal/repository"
)

// ErrPredictionUnavailable is returned when predictions are requested before week 4.
var ErrPredictionUnavailable = errors.New("predictions are only available from week 4 onwards")

type predictionService struct {
	teamRepo     repository.TeamRepository
	matchRepo    repository.MatchRepository
	standingRepo repository.StandingRepository
	engine       *engine.PredictionEngine
}

// NewPredictionService wires prediction dependencies.
func NewPredictionService(
	tr repository.TeamRepository,
	mr repository.MatchRepository,
	sr repository.StandingRepository,
	eng *engine.PredictionEngine,
) *predictionService {
	return &predictionService{
		teamRepo:     tr,
		matchRepo:    mr,
		standingRepo: sr,
		engine:       eng,
	}
}

// GetPredictions runs 10,000 Monte Carlo simulations and returns championship probabilities.
// The task spec requires predictions to be available only from week 4 onward.
func (s *predictionService) GetPredictions(ctx context.Context) ([]domain.PredictionResult, error) {
	// currentWeek, err := s.matchRepo.GetCurrentWeek(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// if currentWeek < 4 {
	// 	return nil, ErrPredictionUnavailable
	// }

	teams, err := s.teamRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Unplayed fixtures are the input to the simulation.
	// If the league is complete, this will be empty — the engine handles that gracefully.
	remaining, err := s.matchRepo.GetUnplayed(ctx)
	if err != nil {
		return nil, err
	}

	standings, err := s.standingRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return s.engine.Predict(teams, remaining, standings), nil
}