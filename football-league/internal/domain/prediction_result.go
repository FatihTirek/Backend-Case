package domain

// PredictionResult holds the Monte Carlo championship probability for one team.
type PredictionResult struct {
	TeamID              int     `json:"team_id"`
	TeamName            string  `json:"team_name"`
	ChampionshipPercent float64 `json:"championship_probability"`
}
