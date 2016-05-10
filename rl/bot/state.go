package bot

import "github.com/saulshanabrook/blockbattle/game"

const numStateFeatures = 10*20 + 8

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
	// for _, cells := range s.Yours.Field {
	// 	for _, cell := range cells {
	// 		xs = append(xs, float64(cell))
	// 	}
	// }

	pieceToNum := func(p game.Piece) float64 {
		for i, tP := range game.AllPieces {
			if tP == p {
				return float64(i)
			}
		}
		panic("invalid piece")
	}
	pos := s.Game.ThisPiecePosition

	xs = append(
		xs,
		pieceToNum(s.Game.ThisPiece),
		pieceToNum(s.Game.NextPiece),
		float64(s.Mine.Combo),
		float64(s.Mine.RowPoints),
		float64(s.Mine.Skips),
		float64(s.Yours.RowPoints),
		float64(pos.Row),
		float64(pos.Column),
	)
	return xs
}
