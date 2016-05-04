package learn

import (
	"math/rand"

	"github.com/saulshanabrook/blockbattle/game"
)

const winningValue = 100

func randAction(st game.State) (game.Location, []game.Move) {
	acts := st.Actions()
	chosenI := rand.Intn(len(acts))
	i := 0
	for loc, mvs := range acts {
		if i == chosenI {
			return loc, mvs
		}
		i++
	}
	panic("cant find any actions")
}

func reward(cur, prev game.State) float64 {
	rpDiff := cur.Mine.RowPoints - prev.Mine.RowPoints

	var win int
	switch cur.Game.Winner {
	case game.None:
		win = 0
	case game.Me:
		win = 1
	case game.You:
		win = -1
	case game.Tie:
		win = 0
	}

	return float64(rpDiff + win*winningValue)
}
