package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/FatihTirek/football-league/internal/domain"
)

var ErrMatchNotFound = errors.New("match not found")

type matchRepository struct {
	db *sql.DB
}

// NewMatchRepository returns a MatchRepository backed by PostgreSQL.
func NewMatchRepository(db *sql.DB) *matchRepository {
	return &matchRepository{db: db}
}

// matchSelectQuery is the shared base query used by all match read methods.
// The double JOIN on teams (aliased ht/at) fetches both team names in one query.
const matchSelectQuery = `
	SELECT
		m.id, m.week, m.home_team_id, m.away_team_id,
		ht.name AS home_team,
		at.name AS away_team,
		m.home_score, m.away_score, m.played
	FROM matches m
	JOIN teams ht ON ht.id = m.home_team_id
	JOIN teams at ON at.id = m.away_team_id`

// scanMatch reads one match row into a domain.Match.
// home_score and away_score are NULL for unplayed matches, so we scan them
// into sql.NullInt64 first, then conditionally set the *int pointer fields.
func scanMatch(rows *sql.Rows) (domain.Match, error) {
	var m domain.Match
	var hs, as sql.NullInt64
	err := rows.Scan(
		&m.ID, &m.Week, &m.HomeTeamID, &m.AwayTeamID,
		&m.HomeTeam, &m.AwayTeam, &hs, &as, &m.Played,
	)
	if err != nil {
		return domain.Match{}, err
	}
	if hs.Valid {
		v := int(hs.Int64)
		m.HomeScore = &v
	}
	if as.Valid {
		v := int(as.Int64)
		m.AwayScore = &v
	}
	return m, nil
}

func (r *matchRepository) GetAll(ctx context.Context) ([]domain.Match, error) {
	rows, err := r.db.QueryContext(ctx, matchSelectQuery + ` ORDER BY m.week, m.id`)
	if err != nil {
		return nil, fmt.Errorf("matchRepo.GetAll: %w", err)
	}
	defer rows.Close()

	var matches []domain.Match
	for rows.Next() {
		m, err := scanMatch(rows)
		if err != nil {
			return nil, fmt.Errorf("matchRepo.GetAll scan: %w", err)
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

func (r *matchRepository) GetByWeek(ctx context.Context, week int) ([]domain.Match, error) {
	rows, err := r.db.QueryContext(ctx, matchSelectQuery + ` WHERE m.week = $1 ORDER BY m.id`, week)
	if err != nil {
		return nil, fmt.Errorf("matchRepo.GetByWeek: %w", err)
	}
	defer rows.Close()

	var matches []domain.Match
	for rows.Next() {
		m, err := scanMatch(rows)
		if err != nil {
			return nil, fmt.Errorf("matchRepo.GetByWeek scan: %w", err)
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

func (r *matchRepository) GetUnplayed(ctx context.Context) ([]domain.Match, error) {
	rows, err := r.db.QueryContext(ctx, matchSelectQuery + ` WHERE m.played = false ORDER BY m.week, m.id`)
	if err != nil {
		return nil, fmt.Errorf("matchRepo.GetUnplayed: %w", err)
	}
	defer rows.Close()

	var matches []domain.Match
	for rows.Next() {
		m, err := scanMatch(rows)
		if err != nil {
			return nil, fmt.Errorf("matchRepo.GetUnplayed scan: %w", err)
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

// GetCurrentWeek returns the MAX week number where played=true.
// COALESCE handles the case where no matches have been played yet (MAX returns NULL → 0).
func (r *matchRepository) GetCurrentWeek(ctx context.Context) (int, error) {
	var week int
	err := r.db.QueryRowContext(ctx,
		`SELECT COALESCE(MAX(week), 0) FROM matches WHERE played = true`,
	).Scan(&week)
	return week, err
}

// SaveWeekResults wraps all match score updates in a single transaction.
// If any update fails the entire batch is rolled back — partial saves are impossible.
func (r *matchRepository) SaveWeekResults(ctx context.Context, matches []domain.Match) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("matchRepo.SaveWeekResults begin: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck — Rollback on committed tx is a no-op

	stmt, err := tx.PrepareContext(ctx,
		`UPDATE matches SET home_score = $1, away_score = $2, played = true WHERE id = $3`)
	if err != nil {
		return fmt.Errorf("matchRepo.SaveWeekResults prepare: %w", err)
	}
	defer stmt.Close()

	for _, m := range matches {
		if m.HomeScore == nil || m.AwayScore == nil {
			continue
		}
		if _, err := stmt.ExecContext(ctx, *m.HomeScore, *m.AwayScore, m.ID); err != nil {
			return fmt.Errorf("matchRepo.SaveWeekResults exec match %d: %w", m.ID, err)
		}
	}
	return tx.Commit()
}

// UpdateResult edits a single match score and marks it as played.
// Returns repository.ErrMatchNotFound when the given ID doesn't exist in the database.
func (r *matchRepository) UpdateResult(ctx context.Context, matchID, homeScore, awayScore int) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE matches SET home_score = $1, away_score = $2, played = true WHERE id = $3`,
		homeScore, awayScore, matchID,
	)
	if err != nil {
		return fmt.Errorf("matchRepo.UpdateResult: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("matchRepo.UpdateResult rows affected: %w", err)
	}
	if n == 0 {
		return ErrMatchNotFound
	}
	return nil
}