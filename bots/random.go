package bots

import (
	"fmt"
	"math/rand"

	"github.com/saulshanabrook/blockbattle/game"
)

// Random bot take 3 random moves each turn
type Random struct{}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func (_ *Random) Act(s *game.State) *[]game.Move {
	possibleMoves := game.AllMoves()

	randMove := func() game.Move {
		i := random(0, len(possibleMoves)-1)
		return possibleMoves[i]
	}
	return &[]game.Move{randMove(), randMove(), randMove()}
}

func (_ *Random) ProcessWinner(w game.Winner) {
	fmt.Println(w)
}
