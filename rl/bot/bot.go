package bot

import (
	"encoding/json"
	"math"

	"github.com/NOX73/go-neural"
	"github.com/NOX73/go-neural/persist"
	"github.com/saulshanabrook/blockbattle/game"
)

type Calculator interface {
	Calculate([]float64) []float64
}

// Bot uses a deep Q network to figure out what action to take
type Bot struct {
	Calculator
}

const numInputs = numActionFeatures + numStateFeatures

func NewNetwork() *neural.Network {
	network := neural.NewNetwork(numInputs, []int{numInputs, 20, 20, 20, 1})
	network.RandomizeSynapses()
	return network
}

// NewFromBinary returns a new NN from a binary representation of stuff
func NewFromBinary(b []byte) (*Bot, error) {
	dump := &persist.NetworkDump{}
	err := json.Unmarshal(b, dump)
	if err != nil {
		return nil, err
	}
	return &Bot{persist.FromDump(dump)}, nil
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
