package learn

import (
	"math/rand"

	"github.com/saulshanabrook/blockbattle/game"
)

const maxExperiences = 100000

type experience struct {
	combFeatures []float64
	reward       float64
	nextSt       game.State
}

type experiences struct {
	exps        [maxExperiences]*experience
	nExps       int
	toAdd       chan *experience
	needPick    chan interface{}
	picked      chan *experience
	readyClosed bool
	readyChan   chan interface{}
}

func newExperiences() *experiences {
	es := experiences{
		exps:        [maxExperiences]*experience{},
		toAdd:       make(chan *experience),
		needPick:    make(chan interface{}),
		picked:      make(chan *experience),
		readyClosed: false,
		readyChan:   make(chan interface{}),
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
			if es.nExps == maxExperiences {
				es.exps[rand.Intn(maxExperiences)] = e
			} else {
				es.exps[es.nExps] = e
			}
			es.nExps++
			if !es.readyClosed {
				es.readyClosed = true
				close(es.readyChan)
			}
		case <-es.needPick:
			es.picked <- es.exps[rand.Intn(es.nExps)]
		}
	}
}

func (es *experiences) add(e *experience) {
	es.toAdd <- e
	return
}
