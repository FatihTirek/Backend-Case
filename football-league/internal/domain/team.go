package domain

// Team represents a football club with strength metrics used by the match simulation engine.
//
// Attack and Defense are multipliers relative to 1.0 (league average).
//   - Higher Attack  → team scores more goals on average.
//   - Higher Defense → team concedes fewer goals on average.
//
// Example: Attack=1.8, Defense=1.3 represents a dominant team like Manchester City.
type Team struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Attack  float64 `json:"attack"`
	Defense float64 `json:"defense"`
}
