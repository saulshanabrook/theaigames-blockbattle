package learn

import (
	"math/rand"

	"github.com/saulshanabrook/blockbattle/game"
)

const maxExperiences = 1000000

type experience struct {
	combFeatures []float64
	reward       float64
	nextSt       game.State
}

type experiences struct {
	exps     *[]*experience
	toAdd    chan *experience
	needPick chan interface{}
	picked   chan *experience
}

func newExperiences() *experiences {
	exps := make([]*experience, 0, maxExperiences)
	es := experiences{
		exps:     &exps,
		toAdd:    make(chan *experience),
		needPick: make(chan interface{}),
		picked:   make(chan *experience),
	}
	go es.process()
	return &es
}

func (es *experiences) pick() *experience {
	es.needPick <- nil
	return <-es.picked
}

func (es *experiences) process() {
	for {
		select {
		case e := <-es.toAdd:
			if len(*es.exps) >= cap(*es.exps) {
				(*es.exps)[rand.Intn(len(*es.exps))] = e
			} else {
				*es.exps = append(*es.exps, e)
			}
		case <-es.needPick:
			es.picked <- (*es.exps)[rand.Intn(len(*es.exps))]
		}
	}
}

func (es *experiences) add(e *experience) {
	es.toAdd <- e
	return
}
