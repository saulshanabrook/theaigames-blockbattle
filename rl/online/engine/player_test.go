package engine

import (
	"os"
	"testing"

	"github.com/saulshanabrook/blockbattle/game"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewPlayers(t *testing.T) {
	Convey("When making new players", t, func() {
		os.Chdir("../../..")

		ps, err := newPlayers()
		So(err, ShouldEqual, nil)
		Convey("We should be able to send stuff for each", func() {

			go func() {
				for st := range ps[0].States {
					ps[0].Moves <- []game.Move{game.MoveDown}
					if st.IsOver() {
						close(ps[0].Moves)
					}
				}
			}()

			sts, mvss := ps[1].States, ps[1].Moves

			for st := range sts {
				So(st, ShouldNotEqual, nil)
				mvss <- []game.Move{game.MoveDown}
				if st.IsOver() {
					close(mvss)
				}
			}
		})
	})
}
