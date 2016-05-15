package bots

import (
	"github.com/saulshanabrook/blockbattle/dqn"
	"github.com/saulshanabrook/blockbattle/game"
	"github.com/saulshanabrook/blockbattle/nn"
)

type DQN struct {
	dqn.Decider
}

// Act takes a random move
func (d *DQN) Act(st game.State) []game.Move {
	a, _ := d.Decide(st)
	return a.Moves
}

func NewDQN(nnHOST string) *DQN {
	return &DQN{nn.NewDecider(nnHOST)}
}
