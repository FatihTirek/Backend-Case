package service

import (
	"context"

	"github.com/FatihTirek/football-league/internal/domain"
)

// PredictionService is the interface the HTTP handler depends on.
type PredictionService interface {
	GetPredictions(ctx context.Context) ([]domain.PredictionResult, error)
}