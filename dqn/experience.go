package dqn

import (
	"math/rand"

	"github.com/saulshanabrook/blockbattle/game"
)

type Experience struct {
	State     game.State
	Action    game.Action
	Reward    float64
	NextState game.State
}

func NewExperience(st game.State, act game.Action, nst game.State) *Experience {
	e := Experience{
		State:     st,
		Action:    act,
		NextState: nst,
	}
	e.calcReward()
	return &e
}

const winningValue = 100

func (e *Experience) calcReward() {
	rpDiff := e.NextState.Mine.RowPoints - e.State.Mine.RowPoints
	var win int
	switch e.NextState.Game.Winner {
	case game.None:
		win = 0
	case game.Me:
		win = 1
	case game.You:
		win = -1
	case game.Tie:
		win = 0
	}
	e.Reward = float64(rpDiff + win*winningValue)
	return
}

type Experiences struct {
	es     []*Experience
	toPush chan *Experience
	toPeek chan *Experience
}

func NewExperiences() *Experiences {
	es := Experiences{
		es:     []*Experience{},
		toPush: make(chan *Experience),
		toPeek: make(chan *Experience),
	}
	go es.work()
	return &es
}

// Push adds an experience to the heap and calculates its reward
func (es *Experiences) Push(e *Experience) {
	es.toPush <- e
	return
}

// Peek grabs a random experience
func (es *Experiences) Peek() *Experience {
	return <-es.toPeek
}

func (es *Experiences) push(e *Experience) {
	es.es = append(es.es, e)
	return
}

func (es *Experiences) peek() *Experience {
	return es.es[rand.Intn(len(es.es))]
}

func (es *Experiences) work() {
	// before we can peek, we have to push at least one
	randE := <-es.toPush
	es.push(randE)
	for {
		select {
		case e := <-es.toPush:
			es.push(e)
		case es.toPeek <- randE:
			randE = es.peek()
		}
	}
}
