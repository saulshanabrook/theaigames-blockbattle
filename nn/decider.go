package nn

import (
	"github.com/gorilla/websocket"
	"github.com/saulshanabrook/blockbattle/game"
)

// Decider communicates to a backend serve, which will decide which action
// to take and what value that action has, given a certain state
//
// It sends
type Decider websocket.Conn

func NewDecider(host string) *Decider {
	return (*Decider)(newConn(host, "/decide"))
}

func (d *Decider) Decide(st game.State) (game.Action, float64) {
	acts := st.Actions()
	handle((*websocket.Conn)(d).WriteJSON(struct {
		State   game.State
		Actions []game.Action
	}{State: st, Actions: st.Actions()}))

	var out struct {
		Index int
		Value float64
	}
	handle((*websocket.Conn)(d).ReadJSON(&out))
	return acts[out.Index], out.Value
}
