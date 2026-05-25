package engine

import (
	"math/rand"

	"github.com/FatihTirek/football-league/internal/domain"
	"github.com/FatihTirek/football-league/pkg/poisson"
)

const (
	leagueAvgGoals = 1.5 // baseline expected goals per team per match (EPL historical average)
	homeAdvantage  = 1.2 // home teams score ~20% more goals than away teams statistically
)

// MatchEngine simulates match scores using the Poisson distribution model.
// Each instance owns its own *rand.Rand source so it is safe to use one
// MatchEngine per goroutine without any mutex overhead.
type MatchEngine struct {
	rng *rand.Rand
}

// NewMatchEngine creates a deterministically seeded engine.
// Pass time.Now().UnixNano() for a random seed each run.
func NewMatchEngine(seed int64) *MatchEngine {
	return &MatchEngine{rng: rand.New(rand.NewSource(seed))}
}

// SimulateMatch computes a scoreline using the Dixon-Coles Poisson model.
//
// Lambda (expected goals) for each team:
//
//	λ = team.Attack × (1 / opponent.Defense) × leagueAvgGoals × homeFactor
//
// A stronger attack raises λ. A stronger opponent defense (higher value) lowers λ.
// The home team receives a 1.2× multiplier; the away team does not.
func (e *MatchEngine) SimulateMatch(home, away domain.Team) (homeScore, awayScore int) {
	homeLambda := home.Attack * (1.0 / away.Defense) * leagueAvgGoals * homeAdvantage
	awayLambda := away.Attack * (1.0 / home.Defense) * leagueAvgGoals

	homeScore = poisson.SimulateGoals(homeLambda, e.rng.Float64())
	awayScore = poisson.SimulateGoals(awayLambda, e.rng.Float64())
	return
}