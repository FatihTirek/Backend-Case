package service

import (
	"context"

	"github.com/FatihTirek/football-league/internal/domain"
)

// WeekResult groups the simulated matches for one week. Used as the response body
// for PlayNextWeek and each element of the PlayAll response.
type WeekResult struct {
	Week    int            `json:"week"`
	Matches []domain.Match `json:"matches"`
}

// LeagueService is the interface the HTTP handler depends on.
type LeagueService interface {
	GetTeams(ctx context.Context) ([]domain.Team, error)
	GetStandings(ctx context.Context) ([]domain.Standing, error)
	GetAllMatches(ctx context.Context) ([]domain.Match, error)
	GetWeekMatches(ctx context.Context, week int) ([]domain.Match, error)
	PlayNextWeek(ctx context.Context) (*WeekResult, error)
	PlayAll(ctx context.Context) ([]WeekResult, error)
	EditMatchResult(ctx context.Context, matchID, homeScore, awayScore int) ([]domain.Standing, error)
}
