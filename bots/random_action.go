package bots

import "github.com/saulshanabrook/blockbattle/game"

// RandomAction bot choses a random possible action (ending stat for block)
// and takes it
type RandomAction struct{}

// Act takes a random move
func (ra RandomAction) Act(s game.State) []game.Move {
	possibleActions := s.Mine.Field.Actions(s.Game.ThisPiece, s.Game.ThisPiecePosition)
	for _, mvs := range possibleActions {
		return mvs
	}
	panic("no possible actions")
}
