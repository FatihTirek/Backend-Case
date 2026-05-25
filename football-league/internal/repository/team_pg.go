package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FatihTirek/football-league/internal/domain"
)

type teamRepository struct {
	db *sql.DB
}

// NewTeamRepository returns a TeamRepository backed by PostgreSQL.
func NewTeamRepository(db *sql.DB) *teamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) GetAll(ctx context.Context) ([]domain.Team, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, attack, defense FROM teams ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("teamRepo.GetAll: %w", err)
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var t domain.Team
		if err := rows.Scan(&t.ID, &t.Name, &t.Attack, &t.Defense); err != nil {
			return nil, fmt.Errorf("teamRepo.GetAll scan: %w", err)
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

func (r *teamRepository) GetByID(ctx context.Context, id int) (domain.Team, error) {
	var t domain.Team
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, attack, defense FROM teams WHERE id = $1`, id,
	).Scan(&t.ID, &t.Name, &t.Attack, &t.Defense)
	if err == sql.ErrNoRows {
		return domain.Team{}, fmt.Errorf("team %d not found", id)
	}
	if err != nil {
		return domain.Team{}, fmt.Errorf("teamRepo.GetByID: %w", err)
	}
	return t, nil
}