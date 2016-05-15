package dqn

import "github.com/saulshanabrook/blockbattle/game"

type Decider interface {
	Decide(game.State) (game.Action, float64)
}
