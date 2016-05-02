package player

import (
	"reflect"
	"strings"
	"testing"

	"github.com/saulshanabrook/blockbattle/game"
)

// TestProcess runs the begining of the example round here and make sure it works
// properly http://theaigames.com/competitions/ai-block-battle/getting-started
// I have changed it a little to test values that are 0 in the example, to make
// sure they are parsed
func TestProcess(t *testing.T) {
	pInput := make(chan string)
	pOutput := make(chan string)
	defer close(pInput)
	defer close(pOutput)

	p := Player{pInput, pOutput}
	sts, win, mvss := p.Process()

	engineSend := func(msgs string) {
		for _, msg := range strings.Split(msgs, "\n") {
			pInput <- msg
		}
	}

	assertState := func(s *game.State) {
		if !reflect.DeepEqual(<-sts, s) {
			t.Errorf("Got back wrong state: %+v", s)
		}
	}

	assertWinner := func(w game.Winner) {
		actualW := <-win
		if actualW != w {
			t.Errorf("Wrong winner: %v", actualW)
		}
	}

	assertEngineGot := func(expMsg string) {
		msg := <-pOutput
		if msg != expMsg {
			t.Errorf("Sent wrong message to server %v", msg)
		}
	}

	engineSend(`settings timebank 10000
settings time_per_move 500
settings player_names player1,player2
settings your_bot player1
settings field_height 20
settings field_width 10`)
	engineSend(`update game round 1
update game this_piece_type O
update game next_piece_type I
update game this_piece_position 4,-1`)
	engineSend(`update player1 row_points 1
update player1 combo 5
update player1 skips 10
update player1 field 0,0,0,0,1,1,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0`)
	engineSend(`update player2 row_points 0
update player2 combo 0
update player2 field 0,0,0,0,1,1,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0;0,0,0,0,0,0,0,0,0,0`)
	engineSend("action moves 10000")
	playerField := game.Field{}
	playerField[0] = [10]game.Cell{0, 0, 0, 0, 1, 1, 0, 0, 0, 0}
	st := game.State{
		Name: "player1",
		Game: &game.GameState{
			Winner:            game.None,
			ThisPiece:         "O",
			NextPiece:         "I",
			ThisPiecePosition: game.Position{Column: 4, Row: -1},
		},
		Mine: &game.PlayerState{
			RowPoints: 1,
			Combo:     5,
			Skips:     10,
			Field:     playerField,
		},
		Yours: &game.PlayerState{
			RowPoints: 0,
			Combo:     0,
			Skips:     0,
			Field:     playerField,
		},
	}
	assertState(&st)
	mvss <- []game.Move{game.MoveLeft, game.MoveLeft, game.MoveLeft, game.MoveLeft, game.MoveDown}

	assertEngineGot("left,left,left,left,down")
	engineSend(`update game winner player1
action moves 100`)
	assertWinner(game.Me)

}
