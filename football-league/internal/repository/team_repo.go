package repository

import (
	"context"

	"github.com/FatihTirek/football-league/internal/domain"
)

// TeamRepository provides read access to team data.
// Teams are seeded at startup and never mutated during a simulation.
type TeamRepository interface {
	GetAll(ctx context.Context) ([]domain.Team, error)
	GetByID(ctx context.Context, id int) (domain.Team, error)
}
