package player

import "github.com/saulshanabrook/blockbattle/game"

// Player can send instructions to the game and get instructions back
type Player struct {
	Input  <-chan string
	Output chan<- string
}

// Process starts a goroutine processing the player. It returns two channels.
// the first you can read from to get the current state. Whenever you recieve
// a state, you should send a list of moves on the second channel and then
// wait for a new state again.
func (p *Player) Process() (<-chan *game.State, <-chan game.Winner, chan<- []game.Move) {
	sts := make(chan *game.State)
	mvss := make(chan []game.Move)
	win := make(chan game.Winner, 1)
	go func() {
		defer close(sts)
		st := (*state)(game.NewState())
		for {
			// wait for a message from the server or some moves to send to it
			// from the user
			select {
			case msg := <-p.Input:
				// we have closed the channel so the file has been deleted so we are
				// done and we can exit
				if msg == "" {
					return
				}
				gotAction, winner, err := st.processLine(msg)
				if err != nil {
					panic(err)
				}
				if gotAction {
					sts <- (*game.State)(st)
				}
				if winner != game.None {
					win <- winner
				}
			case mvs := <-mvss:

				msg, err := serializeMoves(mvs)
				if err != nil {
					panic(err)
				}
				p.Output <- msg
			}
		}
	}()
	return sts, win, mvss
}
