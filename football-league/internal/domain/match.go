package domain

// Match represents a single scheduled fixture between two teams.
//
// HomeScore and AwayScore are pointer types (*int) because there is a semantic
// difference between "the score is 0" and "the match has not been played yet (nil)".
// Using a plain int would make a 0-0 result indistinguishable from an unplayed game.
type Match struct {
	ID         int    `json:"id"`
	Week       int    `json:"week"`
	HomeTeamID int    `json:"home_team_id"`
	AwayTeamID int    `json:"away_team_id"`
	HomeTeam   string `json:"home_team,omitempty"`
	AwayTeam   string `json:"away_team,omitempty"`
	HomeScore  *int   `json:"home_score"`
	AwayScore  *int   `json:"away_score"`
	Played     bool   `json:"played"`
}
