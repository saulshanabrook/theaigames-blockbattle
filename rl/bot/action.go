package bot

import "github.com/saulshanabrook/blockbattle/game"

const numActionFeatures = 3

// ActionFeatures turns a state into a flat list of numbers
//
// includes the board state, the current piece, and the next piece, and the position
func ActionFeatures(loc game.Location) []float64 {
	return []float64{
		float64(loc.Position.Row),
		float64(loc.Position.Column),
		float64(loc.Rotation),
	}
}
