package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FatihTirek/football-league/internal/domain"
)

type standingRepository struct {
	db *sql.DB
}

// NewStandingRepository returns a StandingRepository backed by PostgreSQL.
func NewStandingRepository(db *sql.DB) *standingRepository {
	return &standingRepository{db: db}
}

// GetAll returns the league table sorted by Premier League tiebreaker rules.
func (r *standingRepository) GetAll(ctx context.Context) ([]domain.Standing, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.team_id, t.name, s.played, s.won, s.drawn, s.lost,
		       s.goals_for, s.goals_against, s.goal_diff, s.points
		FROM standings s
		JOIN teams t ON t.id = s.team_id
		ORDER BY s.points DESC, s.goal_diff DESC, s.goals_for DESC`)
	if err != nil {
		return nil, fmt.Errorf("standingRepo.GetAll: %w", err)
	}
	defer rows.Close()

	var standings []domain.Standing
	for rows.Next() {
		var s domain.Standing
		if err := rows.Scan(
			&s.TeamID, &s.TeamName, &s.Played, &s.Won, &s.Drawn, &s.Lost,
			&s.GoalsFor, &s.GoalsAgainst, &s.GoalDiff, &s.Points,
		); err != nil {
			return nil, fmt.Errorf("standingRepo.GetAll scan: %w", err)
		}
		standings = append(standings, s)
	}
	return standings, rows.Err()
}

// Recalculate recomputes all standings from scratch using only the matches table.
//
// A two-stage CTE handles this:
//   - match_stats: generates one row per team per played match (home + away perspective via UNION ALL)
//   - aggregated:  sums those per-match rows into one row per team
//
// The INSERT ... ON CONFLICT DO UPDATE makes this idempotent — safe to call after any
// match edit without risk of double-counting. The matches table is the single source of truth.
func (r *standingRepository) Recalculate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		WITH match_stats AS (
			SELECT
				home_team_id                                           AS team_id,
				1                                                      AS played,
				home_score                                             AS goals_for,
				away_score                                             AS goals_against,
				CASE WHEN home_score > away_score THEN 1 ELSE 0 END    AS won,
				CASE WHEN home_score = away_score THEN 1 ELSE 0 END    AS drawn,
				CASE WHEN home_score < away_score THEN 1 ELSE 0 END    AS lost,
				CASE
					WHEN home_score > away_score THEN 3
					WHEN home_score = away_score THEN 1
					ELSE 0
				END                                                    AS points
			FROM matches WHERE played = true
			UNION ALL
			SELECT
				away_team_id,
				1,
				away_score,
				home_score,
				CASE WHEN away_score > home_score THEN 1 ELSE 0 END,
				CASE WHEN away_score = home_score THEN 1 ELSE 0 END,
				CASE WHEN away_score < home_score THEN 1 ELSE 0 END,
				CASE
					WHEN away_score > home_score THEN 3
					WHEN away_score = home_score THEN 1
					ELSE 0
				END
			FROM matches WHERE played = true
		),
		aggregated AS (
			SELECT
				team_id,
				SUM(played)                         AS played,
				SUM(won)                            AS won,
				SUM(drawn)                          AS drawn,
				SUM(lost)                           AS lost,
				SUM(goals_for)                      AS goals_for,
				SUM(goals_against)                  AS goals_against,
				SUM(goals_for) - SUM(goals_against) AS goal_diff,
				SUM(points)                         AS points
			FROM match_stats
			GROUP BY team_id
		)
		INSERT INTO standings
			(team_id, played, won, drawn, lost, goals_for, goals_against, goal_diff, points)
		SELECT team_id, played, won, drawn, lost, goals_for, goals_against, goal_diff, points
		FROM aggregated
		ON CONFLICT (team_id) DO UPDATE SET
			played        = EXCLUDED.played,
			won           = EXCLUDED.won,
			drawn         = EXCLUDED.drawn,
			lost          = EXCLUDED.lost,
			goals_for     = EXCLUDED.goals_for,
			goals_against = EXCLUDED.goals_against,
			goal_diff     = EXCLUDED.goal_diff,
			points        = EXCLUDED.points,
			updated_at    = NOW()`)
	if err != nil {
		return fmt.Errorf("standingRepo.Recalculate: %w", err)
	}
	return nil
}