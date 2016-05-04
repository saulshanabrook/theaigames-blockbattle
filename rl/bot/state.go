package bot

import "github.com/saulshanabrook/blockbattle/game"

const numStateFeatures = 204

// StateFeatures turns a state into a flat list of numbers
//
// includes the board state, the current piece, and the next piece, and the position
func StateFeatures(s game.State) []float64 {
	xs := make([]float64, 0, numStateFeatures)

	for _, cells := range s.Mine.Field {
		for _, cell := range cells {
			xs = append(xs, float64(cell))
		}
	}

	pieceToNum := func(p game.Piece) float64 {
		for i, tP := range game.AllPieces {
			if tP == p {
				return float64(i)
			}
		}
		panic("invalid piece")
	}

	xs = append(xs, pieceToNum(s.Game.ThisPiece), pieceToNum(s.Game.NextPiece))

	pos := s.Game.ThisPiecePosition
	xs = append(xs, float64(pos.Row), float64(pos.Column))
	return xs
}
