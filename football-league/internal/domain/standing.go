package domain

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