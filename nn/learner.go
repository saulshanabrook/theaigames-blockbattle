package nn

import (
	"github.com/gorilla/websocket"
	"github.com/saulshanabrook/blockbattle/game"
)

type Learner websocket.Conn

func NewLearner(host string) *Learner {
	return (*Learner)(newConn(host, "/learn"))
}

func (l *Learner) Learn(st game.State, act game.Action, val float64) {
	handle((*websocket.Conn)(l).WriteJSON(struct {
		State  game.State
		Action game.Action
		Value  float64
	}{State: st, Action: act, Value: val}))
}
