// Package dqn support generic Deep Q Learning, as developed by
// DeepMind
package dqn

import (
	"math/rand"

	"github.com/saulshanabrook/blockbattle/game"
	"github.com/saulshanabrook/blockbattle/player"
)

// GenPRandActs generates the stream of probability of random action constants
// so that this can decrease over time
func GenPRandActs(n int) <-chan float64 {
	ps := make(chan float64)
	go func() {
		defer close(ps)
		for i := 0; i < n; i++ {
			ps <- calcpRandAct(n, i)
		}
	}()
	return ps
}

func calcpRandAct(n, i int) float64 {
	pComplete := float64(i) / float64(n)
	if pComplete < 0.5 {
		return 1 - 0.9*pComplete*2
	}
	return 0.1
}

func randAction(as []game.Action) game.Action {
	return as[rand.Intn(len(as))]
}

// RunPlayer does Q learning on the player
func RunPlayer(p *player.Player, es *Experiences, pRandAct float64, d Decider) {
	defer close(p.Moves)
	var lastSt game.State
	var lastA game.Action
	var haveLast bool
	for st := range p.States {
		var a game.Action
		if rand.Float64() < pRandAct {
			a = randAction(st.Actions())
		} else {
			a, _ = d.Decide(st)
		}
		p.Moves <- a.Moves
		if haveLast {
			es.Push(&Experience{
				State:     lastSt,
				Action:    lastA,
				NextState: st,
			})
		}
		lastSt = st
		lastA = a
		haveLast = true
	}
}
