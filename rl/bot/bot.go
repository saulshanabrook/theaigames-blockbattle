package bot

import (
	"math"

	"github.com/NOX73/go-neural"
	"github.com/NOX73/go-neural/engine"
	"github.com/saulshanabrook/blockbattle/game"
)

// Bot uses a deep Q network to figure out what action to take
type Bot struct {
	engine.Engine
}

const numInputs = numActionFeatures + numStateFeatures

// New creates a bot with a randomly intialized NN for rewards
func New() *Bot {
	network := neural.NewNetwork(numInputs, []int{numInputs, numInputs, 1})
	network.RandomizeSynapses()
	eng := engine.New(network)
	eng.Start()
	return &Bot{eng}
}

// Act returns the best set of moves for a state
func (b *Bot) Act(st game.State) []game.Move {
	_, mvs, _ := b.BestAction(st)
	return mvs
}

// BestAction returns the best location, moves, and value for that new location for any state
func (b *Bot) BestAction(st game.State) (bestLoc game.Location, bestMvs []game.Move, bestV float64) {
	bestV = -math.MaxFloat64
	stateFts := StateFeatures(st)
	for loc, mvs := range st.Actions() {
		v := b.estValue(stateFts, ActionFeatures(loc))
		if v > bestV {
			bestV = v
			bestLoc = loc
			bestMvs = mvs
		}
	}
	return bestLoc, bestMvs, bestV
}

func (b *Bot) estValue(stateFeatures, actionFeatures []float64) float64 {
	return b.Calculate(append(stateFeatures, actionFeatures...))[0]
}