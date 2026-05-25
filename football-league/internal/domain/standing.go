package domain

import "slices"

// Standing holds a team's accumulated statistics and position in the league table.
type Standing struct {
	TeamID       int    `json:"team_id"`
	TeamName     string `json:"team_name"`
	Played       int    `json:"played"`
	Won          int    `json:"won"`
	Drawn        int    `json:"drawn"`
	Lost         int    `json:"lost"`
	GoalsFor     int    `json:"goals_for"`
	GoalsAgainst int    `json:"goals_against"`
	GoalDiff     int    `json:"goal_diff"`
	Points       int    `json:"points"`
}

// SortStandings sorts a standings slice in-place using Premier League rules.
// It is used by the Monte Carlo engine which operates entirely in memory.
func SortStandings(standings []Standing) {
	slices.SortFunc(standings, func(a, b Standing) int {
		if a.Points != b.Points { return b.Points - a.Points } // Higher points first
		if a.GoalDiff != b.GoalDiff { return b.GoalDiff - a.GoalDiff } // Higher GD first
		return b.GoalsFor - a.GoalsFor // Higher GF first
	})
}
