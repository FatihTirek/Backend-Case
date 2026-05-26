package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/FatihTirek/football-league/internal/domain"
	"github.com/FatihTirek/football-league/internal/engine"
	"github.com/FatihTirek/football-league/internal/repository"
)

const totalWeeks = 6

var (
	ErrLeagueComplete = errors.New("league is complete: all 6 weeks have been played")
	ErrMatchNotFound  = errors.New("match not found")
)

type leagueService struct {
	teamRepo     repository.TeamRepository
	matchRepo    repository.MatchRepository
	standingRepo repository.StandingRepository
	engine       *engine.MatchEngine
}

// NewLeagueService wires dependencies. Accepts interfaces so any implementation
// (real DB, in-memory, mock) can be injected without touching this file.
func NewLeagueService(
	tr repository.TeamRepository,
	mr repository.MatchRepository,
	sr repository.StandingRepository,
	eng *engine.MatchEngine,
) *leagueService {
	return &leagueService{
		teamRepo:     tr,
		matchRepo:    mr,
		standingRepo: sr,
		engine:       eng,
	}
}

func (s *leagueService) GetTeams(ctx context.Context) ([]domain.Team, error) {
	return s.teamRepo.GetAll(ctx)
}

func (s *leagueService) GetStandings(ctx context.Context) ([]domain.Standing, error) {
	return s.standingRepo.GetAll(ctx)
}

func (s *leagueService) GetAllMatches(ctx context.Context) ([]domain.Match, error) {
	return s.matchRepo.GetAll(ctx)
}

func (s *leagueService) GetWeekMatches(ctx context.Context, week int) ([]domain.Match, error) {
	if week < 1 || week > totalWeeks {
		return nil, fmt.Errorf("week must be between 1 and %d", totalWeeks)
	}
	return s.matchRepo.GetByWeek(ctx, week)
}

// PlayNextWeek finds the next unplayed week and simulates it.
func (s *leagueService) PlayNextWeek(ctx context.Context) (*WeekResult, error) {
	current, err := s.matchRepo.GetCurrentWeek(ctx)
	if err != nil {
		return nil, err
	}
	if current >= totalWeeks {
		return nil, ErrLeagueComplete
	}
	return s.playWeek(ctx, current + 1)
}

// PlayAll simulates every remaining unplayed week sequentially.
func (s *leagueService) PlayAll(ctx context.Context) ([]WeekResult, error) {
	current, err := s.matchRepo.GetCurrentWeek(ctx)
	if err != nil {
		return nil, err
	}
	if current >= totalWeeks {
		return nil, ErrLeagueComplete
	}

	var results []WeekResult
	for week := current + 1; week <= totalWeeks; week++ {
		wr, err := s.playWeek(ctx, week)
		if err != nil {
			return results, fmt.Errorf("week %d: %w", week, err)
		}
		results = append(results, *wr)
	}
	return results, nil
}

func (s *leagueService) ResetLeague(ctx context.Context) error {
	if err := s.matchRepo.Reset(ctx); err != nil {
		return err
	}
	if err := s.standingRepo.Reset(ctx); err != nil {
		return err
	}
	return nil
}

// EditMatchResult updates one match score then recalculates all standings from scratch.
// We never surgically adjust standings — we always recompute from the matches table,
// which is the single source of truth. This avoids any stale-data bugs.
func (s *leagueService) EditMatchResult(ctx context.Context, matchID, homeScore, awayScore int) ([]domain.Standing, error) {
	if homeScore < 0 || awayScore < 0 {
		return nil, fmt.Errorf("scores cannot be negative")
	}
	if err := s.matchRepo.UpdateResult(ctx, matchID, homeScore, awayScore); err != nil {
		if errors.Is(err, repository.ErrMatchNotFound) {
			return nil, ErrMatchNotFound
		}
		return nil, err
	}
	if err := s.standingRepo.Recalculate(ctx); err != nil {
		return nil, err
	}
	return s.standingRepo.GetAll(ctx)
}

// playWeek is the core internal method: simulate → persist → recalculate.
func (s *leagueService) playWeek(ctx context.Context, week int) (*WeekResult, error) {
	fixtures, err := s.matchRepo.GetByWeek(ctx, week)
	if err != nil {
		return nil, err
	}

	teams, err := s.teamRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	teamMap := make(map[int]domain.Team, len(teams))
	for _, t := range teams {
		teamMap[t.ID] = t
	}

	// "for i := range" — modifies the actual slice element, not a copy.
	// "for _, f := range" would modify a throwaway copy and the slice would be unchanged.
	for i := range fixtures {
		if fixtures[i].Played {
			continue
		}
		hs, as := s.engine.SimulateMatch(teamMap[fixtures[i].HomeTeamID], teamMap[fixtures[i].AwayTeamID])
		fixtures[i].HomeScore = &hs
		fixtures[i].AwayScore = &as
		fixtures[i].Played = true
	}

	if err := s.matchRepo.SaveWeekResults(ctx, fixtures); err != nil {
		return nil, err
	}
	if err := s.standingRepo.Recalculate(ctx); err != nil {
		return nil, err
	}
	return &WeekResult{Week: week, Matches: fixtures}, nil
}