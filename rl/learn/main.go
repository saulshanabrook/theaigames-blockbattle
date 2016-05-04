package learn

import (
	"math/rand"
	"sync"

	"github.com/saulshanabrook/blockbattle/game"
	"github.com/saulshanabrook/blockbattle/player"
	"github.com/saulshanabrook/blockbattle/rl/bot"
	bbEngine "github.com/saulshanabrook/blockbattle/rl/engine"
)

// this is the back prop "speed"
// https://github.com/NOX73/go-neural/blob/f327eff30de74b5b2d18236415bd35a2ee5c4e59/learn/learn.go#L47
// 0.9 just seemed like an OK constant
const nnSpeed = 0.9

// Learner handles all the agent state that we learn
type Learner struct {
	b    *bot.Bot
	exps *experiences
	c    LearnerConfig
}

type LearnerConfig struct {
	MinibatchSize  int
	DiscountFactor float64
}

// copied from deepmind paper
var DefaultLearnerConfig = LearnerConfig{
	32,
	0.9,
}

// NewLearner creates a blank agent state
func NewLearner(c LearnerConfig) *Learner {
	return &Learner{
		bot.New(),
		newExperiences(),
		c,
	}
}

// RunEpisodes runs n epsidoes starting epsilon at 1 and bringing it linearly
// to 0.1 over the first 1/2 of the trainin and keeping it at 0.1 for the rest
func (l *Learner) RunEpisodes(n int) error {
	for i := 0; i < n; i++ {
		pComplete := float64(i) / float64(n)
		var pRandAct float64
		if pComplete < 0.5 {
			pRandAct = 1 - 0.9*pComplete*2
		} else {
			pRandAct = 0.1
		}
		if err := l.RunEpisode(pRandAct); err != nil {
			return err
		}
	}
	return nil

}

// RunEpisode starts up one game and does learning on both players
func (l *Learner) RunEpisode(pRandAct float64) error {
	players, err := bbEngine.NewPlayers()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, p := range players {
		wg.Add(1)
		go func(p player.Player) {
			defer wg.Done()
			l.play(p, pRandAct)
		}(p)
	}
	wg.Wait()
	return nil
}

func (l *Learner) play(p player.Player, pRandAct float64) {
	defer close(p.Moves)
	st := <-p.States
	for {
		var loc game.Location
		var mvs []game.Move
		if rand.Float64() < pRandAct {
			loc, mvs = randAction(st)
		} else {
			loc, mvs, _ = l.b.BestAction(st)
		}
		p.Moves <- mvs
		nextSt := <-p.States

		// got null value which means channel is closed
		// and game ended because of timeout
		if (nextSt == game.State{}) {
			return
		}
		l.recordExperience(st, loc, nextSt)
		go l.trainMinibatch()
		if nextSt.IsOver() {
			return
		}
		st = nextSt
	}
}

func (l *Learner) trainMinibatch() {
	for i := 0; i < l.c.MinibatchSize; i++ {
		exp := l.exps.pick()
		var nextVal float64
		if exp.nextSt.IsOver() {
			nextVal = 0
		} else {
			_, _, nextVal = l.b.BestAction(exp.nextSt)
		}

		l.b.Learn(
			exp.combFeatures,
			[]float64{exp.reward + l.c.DiscountFactor*nextVal},
			nnSpeed,
		)
	}
}

func (l *Learner) recordExperience(st game.State, loc game.Location, nextSt game.State) {
	l.exps.add(&experience{
		combFeatures: append(bot.StateFeatures(st), bot.ActionFeatures(loc)...),
		reward:       reward(nextSt, st),
		nextSt:       nextSt,
	})
}
