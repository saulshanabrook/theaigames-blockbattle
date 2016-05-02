package bots

import (
	"math/rand"

	"github.com/saulshanabrook/blockbattle/game"
)

// Random bot take 3 random moves each turn
type Random struct{}

// Act takes a random move
func (r Random) Act(_ *game.State) []game.Move {
	possibleMoves := game.AllMoves

	randMove := func() game.Move {
		i := rand.Intn(len(possibleMoves))
		return possibleMoves[i]
	}
	return []game.Move{randMove(), randMove(), randMove()}
}
