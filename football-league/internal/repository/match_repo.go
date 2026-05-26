package repository

import (
	"context"

	"github.com/FatihTirek/football-league/internal/domain"
)

// MatchRepository manages match fixtures and their results.
type MatchRepository interface {
	// GetAll returns every fixture in the league, ordered by week then match id.
	GetAll(ctx context.Context) ([]domain.Match, error)

	// GetByWeek returns the two fixtures scheduled for a specific week.
	GetByWeek(ctx context.Context, week int) ([]domain.Match, error)

	// GetUnplayed returns all fixtures that have not yet been simulated.
	// Used by the prediction engine to know which matches remain.
	GetUnplayed(ctx context.Context) ([]domain.Match, error)

	// GetCurrentWeek returns the highest week number that has been fully played.
	// Returns 0 if no matches have been played yet (start of season).
	GetCurrentWeek(ctx context.Context) (int, error)

	// SaveWeekResults persists simulated scores for all matches in a given week
	// inside a single database transaction. If any update fails, all are rolled back.
	SaveWeekResults(ctx context.Context, matches []domain.Match) error

	// UpdateResult edits the score of a specific, already-played match.
	// Implements the bonus "edit results" feature. Returns ErrMatchNotFound
	// if the given matchID does not exist in the database.
	UpdateResult(ctx context.Context, matchID, homeScore, awayScore int) error

	// Reset clears all match results and resets the league to its initial state.
	Reset(ctx context.Context) error
}
