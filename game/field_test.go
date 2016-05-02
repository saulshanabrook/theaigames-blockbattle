package game

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestField(t *testing.T) {

	Convey("restingPositions", t, func() {

		allBottom := func(ps []Position) {
			for col := 0; col < numColumns; col++ {
				So(ps, ShouldContain, Position{Column: col, Row: numRows - 1})
			}
			return
		}
		Convey("When the field is blank", func() {
			f := Field{}
			ps := f.restingPositions()
			allBottom(ps)
			So(ps, ShouldHaveLength, numColumns)
		})

		Convey("When it has some random blocks in it", func() {
			f := Field{}

			// add some block at the top row
			f[0][0] = Block

			// add a block in a middle row
			f[5][5] = Block

			ps := f.restingPositions()
			allBottom(ps)
			So(ps, ShouldContain, Position{Column: 5, Row: 4})
			So(ps, ShouldHaveLength, numColumns+1)
		})
	})

	Convey("validPlacement", t, func() {

		f := Field{}
		f[5][5] = Block
		Convey("the o block", func() {
			Convey("cant be rotated", func() {
				So(f.validPlacement(O, Location{Rotation: RotatedLeft}), ShouldBeFalse)
			})
			Convey("if its not rotated should fit a lot of places", func() {
				for row := -2; row <= numRows; row++ {
					for col := -1; col <= numColumns; col++ {
						var canPlace bool
						if col < 0 || col > (numColumns-2) {
							canPlace = false
						} else if row < -1 || row > (numRows-2) {
							canPlace = false
						} else if (row == 5 || row == 4) && (col == 5 || col == 4) {
							canPlace = false
						} else {
							canPlace = true
						}
						l := Location{Position: Position{Column: col, Row: row}}
						So(
							f.validPlacement(O, l),
							ShouldEqual,
							canPlace,
						)
					}
				}
			})
		})
	})
	Convey("movesFrom", t, func() {
		Convey("O block moves down", func() {
			f := Field{}
			mvs, err := f.movesFrom(
				O,
				Position{Column: 5, Row: -1},
				Location{
					Position: Position{Column: 5, Row: 3},
					Rotation: RotatedUp,
				},
			)
			So(err, ShouldEqual, nil)
			So(mvs, ShouldResemble, []Move{MoveDown, MoveDown, MoveDown, MoveDown})
		})
		Convey("O can't move off field", func() {
			f := Field{}
			_, err := f.movesFrom(
				O,
				Position{Column: 5, Row: -1},
				Location{
					Position: Position{Column: 5, Row: 19},
					Rotation: RotatedUp,
				},
			)
			So(err, ShouldNotEqual, nil)
		})
		Convey("O can't move into block", func() {
			f := Field{}
			f[3][5] = Block
			_, err := f.movesFrom(
				O,
				Position{Column: 5, Row: -1},
				Location{
					Position: Position{Column: 5, Row: 3},
					Rotation: RotatedUp,
				},
			)
			So(err, ShouldNotEqual, nil)
		})
		Convey("O can get around blocks", func() {
			f := Field{}
			f[17] = [10]Cell{Block, Block, Block, Block, Block, Block, Block, Block, Empty, Empty}
			_, err := f.movesFrom(
				O,
				Position{Column: 5, Row: -1},
				Location{
					Position: Position{Column: 5, Row: 18},
					Rotation: RotatedUp,
				},
			)
			So(err, ShouldEqual, nil)
		})
		Convey("O can't move through blocks", func() {
			f := Field{}
			f[17] = [10]Cell{Block, Block, Block, Block, Block, Block, Block, Block, Block, Block}
			_, err := f.movesFrom(
				O,
				Position{Column: 5, Row: -1},
				Location{
					Position: Position{Column: 5, Row: 18},
					Rotation: RotatedUp,
				},
			)
			So(err, ShouldNotEqual, nil)
		})

		Convey("I can rotate", func() {
			f := Field{}
			mvs, err := f.movesFrom(
				L,
				Position{Column: 5, Row: 5},
				Location{
					Position: Position{Column: 5, Row: 5},
					Rotation: RotatedLeft,
				},
			)
			So(err, ShouldEqual, nil)
			So(mvs, ShouldResemble, []Move{MoveTurnLeft})
		})
	})
	Convey("Actions", t, func() {
		Convey("O can end up", func() {
			f := Field{}
			f[5][5] = Block
			as := f.Actions(O, Position{Row: -1, Column: 5})
			Convey("On bottom", func() {
				for col := 0; col < 9; col++ {
					So(as, ShouldContainKey, Location{Position: Position{Row: 18, Column: col}, Rotation: RotatedUp})
				}
			})
			Convey("On floating piece", func() {
				So(as, ShouldContainKey, Location{Position: Position{Row: 3, Column: 5}, Rotation: RotatedUp})
				So(as, ShouldContainKey, Location{Position: Position{Row: 3, Column: 4}, Rotation: RotatedUp})
			})
			Convey("and not anywhere else", func() {
				So(as, ShouldHaveLength, 9+2)
			})
		})
	})
}
