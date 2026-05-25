package engine

import (
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/FatihTirek/football-league/internal/domain"
)

const defaultSimCount = 10_000

// PredictionEngine runs Monte Carlo simulations to compute championship probabilities.
type PredictionEngine struct {
	simCount int
}

// NewPredictionEngine creates an engine that runs 10,000 simulations by default.
func NewPredictionEngine() *PredictionEngine {
	return &PredictionEngine{simCount: defaultSimCount}
}

// Predict distributes simCount simulations across all CPU cores and returns
// each team's championship win percentage.
//
// Concurrency model:
//  1. A buffered channel is pre-filled with simCount job tokens, then closed.
//  2. workerCount goroutines (one per CPU core) drain the channel.
//  3. Each goroutine owns its own MatchEngine (private rand.Rand) — no mutex needed.
//  4. Winners are tallied via atomic.AddInt64 — no mutex needed.
//  5. wg.Wait() blocks until every goroutine finishes.
func (p *PredictionEngine) Predict(
	teams []domain.Team,
	remainingFixtures []domain.Match,
	currentStandings []domain.Standing,
) []domain.PredictionResult {

	winCounts := make([]int64, len(teams))

	// teamIndex: TeamID → position in the teams/winCounts slice (O(1) lookup)
	teamIndex := make(map[int]int, len(teams))
	for i, t := range teams {
		teamIndex[t.ID] = i
	}

	// teamMap: TeamID → Team struct for engine lookups inside each simulation
	teamMap := make(map[int]domain.Team, len(teams))
	for _, t := range teams {
		teamMap[t.ID] = t
	}

	jobs := make(chan struct{}, p.simCount)
	for i := 0; i < p.simCount; i++ {
		jobs <- struct{}{}
	}
	close(jobs) // closing signals workers to stop once the channel is drained

	var wg sync.WaitGroup
	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Private rand source per goroutine: rand.Rand is NOT goroutine-safe.
			localEngine := &MatchEngine{rng: rand.New(rand.NewSource(time.Now().UnixNano()))}

			for range jobs {
				idx := simulateSeason(localEngine, teamIndex, teamMap, remainingFixtures, currentStandings)
				if idx >= 0 {
					// CPU-level atomic increment — no lock contention on the hot path.
					atomic.AddInt64(&winCounts[idx], 1)
				}
			}
		}()
	}

	wg.Wait()

	results := make([]domain.PredictionResult, len(teams))
	for i, t := range teams {
		results[i] = domain.PredictionResult{
			TeamID:              t.ID,
			TeamName:            t.Name,
			ChampionshipPercent: float64(winCounts[i]) / float64(p.simCount) * 100,
		}
	}
	return results
}

// simulateSeason runs one complete simulation from the current league state.
// Returns the slice index of the winning team, or -1 if indeterminate.
//
// Deep-copies standings before mutating: without this, all goroutines would
// corrupt the same shared slice and produce completely wrong results.
func simulateSeason(
	eng *MatchEngine,
	teamIndex map[int]int,
	teamMap map[int]domain.Team,
	fixtures []domain.Match,
	baseStandings []domain.Standing,
) int {
	// Deep copy: "s := baseStandings[i]" creates a new struct value on the stack.
	standings := make(map[int]*domain.Standing, len(baseStandings))
	for _, s := range baseStandings {
		s := s // create a unique copy for each iteration
		standings[s.TeamID] = &s
	}

	for _, f := range fixtures {
		hs, as := eng.SimulateMatch(teamMap[f.HomeTeamID], teamMap[f.AwayTeamID])
		applySimResult(standings, f.HomeTeamID, f.AwayTeamID, hs, as)
	}

	// Find winner: linear scan is fine — only 4 teams.
	winnerID, bestPts, bestGD, bestGF := -1, -1, -999, -1
	for id, s := range standings {
		if s.Points > bestPts ||
			(s.Points == bestPts && s.GoalDiff > bestGD) ||
			(s.Points == bestPts && s.GoalDiff == bestGD && s.GoalsFor > bestGF) {
			bestPts, bestGD, bestGF = s.Points, s.GoalDiff, s.GoalsFor
			winnerID = id
		}
	}

	if idx, ok := teamIndex[winnerID]; ok {
		return idx
	}
	return -1
}

// applySimResult updates the in-memory standings map for one simulated match.
func applySimResult(standings map[int]*domain.Standing, homeID, awayID, hs, as int) {
	h, a := standings[homeID], standings[awayID]

	h.Played++
	a.Played++
	h.GoalsFor += hs
	h.GoalsAgainst += as
	a.GoalsFor += as
	a.GoalsAgainst += hs
	h.GoalDiff = h.GoalsFor - h.GoalsAgainst
	a.GoalDiff = a.GoalsFor - a.GoalsAgainst

	switch {
	case hs > as:
		h.Won++
		h.Points += 3
		a.Lost++
	case as > hs:
		a.Won++
		a.Points += 3
		h.Lost++
	default:
		h.Drawn++
		a.Drawn++
		h.Points++
		a.Points++
	}
}