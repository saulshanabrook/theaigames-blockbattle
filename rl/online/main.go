package main

import (
	"sync"

	"github.com/saulshanabrook/blockbattle/dqn"
	"github.com/saulshanabrook/blockbattle/nn"
	"github.com/saulshanabrook/blockbattle/player"
	"github.com/saulshanabrook/blockbattle/rl/online/engine"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

const numWorkers = 2
const domain = "50.133.191.147"
const n = 10000

func learnExperience(es *dqn.Experiences) {
	l := nn.NewLearner(domain)
	esD := nn.NewDecider(domain)
	dqn.LearnExperiences(l, esD, es)
}

func worker(ps <-chan *player.Player, es *dqn.Experiences, pRandActs <-chan float64) {
	d := nn.NewDecider(domain)
	for pRandAct := range pRandActs {
		dqn.RunPlayer(<-ps, es, pRandAct, d)
	}
}

func workAll(es *dqn.Experiences) {
	pRandActs := dqn.GenPRandActs(n)
	ps := engine.NewPlayers()
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			worker(ps, es, pRandActs)
		}()
	}
	wg.Wait()
}

func main() {
	// logrus.SetLevel(logrus.DebugLevel)
	es := dqn.NewExperiences()
	go learnExperience(es)
	workAll(es)
	return
}
