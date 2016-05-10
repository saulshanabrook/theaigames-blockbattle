package learn

import (
	"math/rand"
	"sync"

	"github.com/NOX73/go-neural"
	"github.com/NOX73/go-neural/engine"
	"github.com/NOX73/go-neural/persist"
	"github.com/Sirupsen/logrus"
	"github.com/saulshanabrook/blockbattle/game"
	"github.com/saulshanabrook/blockbattle/player"
	"github.com/saulshanabrook/blockbattle/rl/bot"
	bbEngine "github.com/saulshanabrook/blockbattle/rl/engine"
)

const DiscountFactor = 0.9

// this is the back prop "speed"
// https://github.com/NOX73/go-neural/blob/f327eff30de74b5b2d18236415bd35a2ee5c4e59/learn/learn.go#L47
// copied from learning rate at https://gist.github.com/EderSantana/c7222daa328f0e885093#file-qlearn-py-L124
const nnSpeed = 0.2

// get new weights every n rounds
const switchEvery = 500

// Learner handles all the agent state that we learn
type Learner struct {
	e       engine.Engine
	exps    *experiences
	nTrains int
}

// NewLearner creates a blank agent state
func NewLearner(n *neural.Network) *Learner {
	eng := engine.New(n)
	eng.Start()
	return &Learner{
		eng,
		newExperiences(),
		0,
	}
}

// Persist saves this NN to a file
func (l *Learner) Persist(filename string) {
	persist.DumpToFile(filename, l.e.Dump())
}

// RunEpisodes runs n epsidoes starting epsilon at 1 and bringing it linearly
// to 0.1 over the first 1/2 of the trainin and keeping it at 0.1 for the rest
func (l *Learner) RunEpisodes(n int, nWorkers int) error {
	logrus.WithFields(logrus.Fields{
		"n episodes": n,
		"n workers":  nWorkers,
	}).Info("Starting Episodes")
	ps := make(chan float64)

	var wg sync.WaitGroup
	wg.Add(nWorkers)
	for ; nWorkers > 0; nWorkers-- {
		go func() {
			l.work(ps)
			wg.Done()
		}()
	}
	go l.startTraining()
	go l.fillPs(n, ps)
	wg.Wait()
	return nil
}

// RunEpisodesSingle only uses one worker
func (l *Learner) RunEpisodesSingle(n int) error {
	logrus.WithFields(logrus.Fields{
		"n episodes": n,
	}).Info("Starting Episodes")
	ps := make(chan float64)
	go l.startTraining()
	go l.fillPs(n, ps)
	l.work(ps)
	return nil
}
func (l *Learner) fillPs(n int, ps chan<- float64) {
	for i := 0; i < n; i++ {
		pComplete := float64(i) / float64(n)
		var pRandAct float64
		if pComplete < 0.5 {
			pRandAct = 1 - 0.9*pComplete*2
		} else {
			pRandAct = 0.1
		}
		logrus.WithFields(logrus.Fields{"p": pRandAct, "i": i}).Info("Running episode")
		ps <- pRandAct
	}
	close(ps)
}

func (l *Learner) work(ps <-chan float64) {
	logrus.WithFields(logrus.Fields{}).Info("Starting worker")
	i := 0
	var b *bot.Bot
	for p := range ps {
		if i%switchEvery == 0 {
			logrus.WithFields(logrus.Fields{}).Info("Getting new bot for worker")

			b = &bot.Bot{Calculator: persist.FromDump(l.e.Dump())}
		}
		l.RunEpisode(b, p)
		i++
	}
}

// RunEpisode starts up one game and does learning on both players
func (l *Learner) RunEpisode(b *bot.Bot, pRandAct float64) error {
	logrus.WithFields(logrus.Fields{}).Debug("Starting RunEpisode")
	defer logrus.WithFields(logrus.Fields{}).Debug("Finishing RunEpisode")
	players, err := bbEngine.NewPlayers()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, p := range players {
		wg.Add(1)
		go func(p player.Player) {
			defer wg.Done()
			l.play(b, p, pRandAct)
		}(p)
	}
	wg.Wait()
	return nil
}

// play will play a game, taking the best action it can most of the time
// and a random one `pRandAct`% of the time.
// It will save all experiences it comes accross as well.
func (l *Learner) play(b *bot.Bot, p player.Player, pRandAct float64) {
	logrus.WithFields(logrus.Fields{}).Debug("Starting play")
	defer logrus.WithFields(logrus.Fields{}).Debug("Done play")

	// tell the engine we dont want to send any more moves once we finish
	defer close(p.Moves)
	logrus.WithFields(logrus.Fields{}).Debug("Getting state")
	st := <-p.States
	logrus.WithFields(logrus.Fields{}).Debug("Done getting state")
	for {
		var loc game.Location
		var mvs []game.Move
		if rand.Float64() < pRandAct {
			logrus.WithFields(logrus.Fields{}).Debug("Deciding random action")
			// get a random action pRandAct% of the time
			loc, mvs = randAction(st)
			logrus.WithFields(logrus.Fields{}).Debug("Done deciding random action")
		} else {
			logrus.WithFields(logrus.Fields{}).Debug("Deciding best action")
			// otherwise use our neural network to find the best action
			loc, mvs, _ = b.BestAction(st)
			logrus.WithFields(logrus.Fields{}).Debug("Done deciding best action")
		}
		logrus.WithFields(logrus.Fields{}).Debug("Sending moves")
		p.Moves <- mvs
		logrus.WithFields(logrus.Fields{}).Debug("Done sending moves/getting next state")
		// get our next state so we can store this experience
		nextSt := <-p.States
		// got null value which means the states channel has closed
		// and game ended because of timeout. If this is the case
		// dont record anything and move on
		if (nextSt == game.State{}) {
			logrus.WithFields(logrus.Fields{}).Debug("Game closed from timeout")
			return
		}
		// By recording the current state, the action we took (loc), and the next
		// state, we can later train on this by computing the reward given the
		// two states
		logrus.WithFields(logrus.Fields{}).Debug("Recording experience")
		l.recordExperience(st, loc, nextSt)
		logrus.WithFields(logrus.Fields{}).Debug("Done Recording experience")

		// if the game has ended, stop playing
		if nextSt.IsOver() {
			logrus.WithFields(logrus.Fields{}).Debug("Game over")
			return
		}
		st = nextSt
	}
}

func (l *Learner) startTraining() {
	<-l.exps.readyChan
	for {
		l.train()
	}
}

func (l *Learner) train() {
	logrus.WithFields(logrus.Fields{}).Debug("Starting train")
	defer logrus.WithFields(logrus.Fields{}).Debug("Done train")
	exp := l.exps.pick()
	var nextVal float64
	if exp.nextSt.IsOver() {
		nextVal = 0
	} else {
		_, _, nextVal = (&bot.Bot{Calculator: l.e}).BestAction(exp.nextSt)
	}

	l.e.Learn(
		exp.combFeatures,
		[]float64{exp.reward + DiscountFactor*nextVal},
		nnSpeed,
	)
	l.nTrains++
}

func (l *Learner) recordExperience(st game.State, loc game.Location, nextSt game.State) {
	l.exps.add(&experience{
		combFeatures: append(bot.StateFeatures(st), bot.ActionFeatures(loc)...),
		reward:       reward(nextSt, st),
		nextSt:       nextSt,
	})
}
