package player

import "github.com/saulshanabrook/blockbattle/game"

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
	return
}

type Player struct {
	States <-chan game.State
	Moves  chan<- []game.Move
}

func Parse(msgs <-chan string) <-chan game.State {
	sts := make(chan game.State)
	go func() {
		defer close(sts)
		st := (state)(game.NewState())
		for msg := range msgs {
			gotAction, err := st.processLine(msg)
			handleErr(err)
			isOver := game.State(st).IsOver()
			if gotAction || isOver {
				sts <- (game.State)(st)
			}
		}
	}()
	return sts
}

func Serialize(mvss <-chan []game.Move) <-chan string {
	msgs := make(chan string)
	go func() {
		for mvs := range mvss {
			msg, err := serializeMoves(mvs)
			handleErr(err)
			msgs <- msg
		}
		close(msgs)
	}()
	return msgs
}

// Process starts a goroutine processing the player. It returns two channels.
// the first you can read from to get the current state. Whenever you recieve
// a state, you should send a list of moves on the second channel and then
// wait for a new state again.
// func (p *Player) Process() (<-chan game.State, chan<- []game.Move) {
// 	sts := make(chan game.State)
// 	mvss := make(chan []game.Move)
// 	go func() {
// 		defer close(sts)
// 		st := (state)(game.NewState())
// 		for msg := range p.Input {
// 			gotAction, err := st.processLine(msg)
// 			if err != nil {
// 				panic(err)
// 			}
// 			if gotAction || game.State(st).IsOver() {
// 				sts <- (game.State)(st)
// 			}
// 		}
// 	}()
//
// 	go func() {
// 		defer close(p.Output)
// 		for mvs := range mvss {
// 			msg, err := serializeMoves(mvs)
// 			if err != nil {
// 				panic(err)
// 			}
// 			p.Output <- msg
// 		}
// 	}()
// 	return sts, mvss
// }
