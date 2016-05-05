package player

import "github.com/saulshanabrook/blockbattle/game"

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
	return
}

// Player holds everything you need to interact as a blockbot player
type Player struct {
	// States holds each of the states as they are recieved from the engine
	// you only get a state if an action is required or if the game is over
	States <-chan game.State
	// Moves contains the moves you want to send, after you get a non final state
	Moves chan<- []game.Move
	// Done is closed after all the moves have been sent
	Done <-chan interface{}
}

// Parse converts a channels of messages from the game engine into one of states
// by combing the messages as they are read
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

// Serialize takes a list of moves and returns the string for those
// moves to send to the game engine
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
