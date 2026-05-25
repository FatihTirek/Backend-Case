package repository

import (
	"context"

	"github.com/FatihTirek/football-league/internal/domain"
)

// StandingRepository manages the persistent league table.
type StandingRepository interface {
	// GetAll returns standings sorted by Points DESC, GoalDiff DESC, GoalsFor DESC.
	GetAll(ctx context.Context) ([]domain.Standing, error)

	// Recalculate wipes and recomputes all standings from the matches table.
	// This operation is IDEMPOTENT: calling it multiple times with the same match
	// data always produces the same result. This property makes it safe to call
	// after any match edit without worrying about double-counting or stale data.
	Recalculate(ctx context.Context) error
}
