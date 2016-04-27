package player

import (
	"strconv"
	"strings"

	"github.com/saulshanabrook/blockbattle/game"
)

type state game.State

// processLine takes a line from the server and updates the state with it,
// returning if we have a winner and if we got an action
func (s *state) processLine(line string) (gotAction bool, winner game.Winner, err error) {
	parts := strings.Split(line, " ")
	cmd, rest := parts[0], parts[1:]
	switch cmd {
	case "settings":
		return false, game.None, s.processSettings(rest[0], rest[1])
	case "update":
		winner, err = s.processUpdate(rest[0], rest[1], rest[2])
		return false, winner, err
	case "action":
		return true, game.None, nil
	}
	panic(line)
}

func (s *state) processSettings(type_, value string) (err error) {
	switch type_ {
	case "your_bot":
		s.Name = value
	}
	return
}

func (s *state) processUpdate(player, type_, value string) (game.Winner, error) {
	if player == "game" {
		return (*gameState)(s.Game).processUpdate(s.Name, type_, value)
	}
	if player == s.Name {
		return game.None, (*playerState)(s.Mine).processUpdate(type_, value)
	}
	return game.None, (*playerState)(s.Yours).processUpdate(type_, value)
}

// newPosition turns a string like `4,-1` into a position
func newPosition(s string) (p *game.Position, err error) {
	p = &game.Position{}
	ss := strings.Split(s, ",")
	p.X, err = strconv.Atoi(ss[0])
	if err != nil {
		return
	}
	p.Y, err = strconv.Atoi(ss[1])
	return
}

type playerState game.PlayerState

func (ps *playerState) processUpdate(type_, value string) (err error) {
	switch type_ {
	case "row_points":
		ps.RowPoints, err = strconv.Atoi(value)
	case "combo":
		ps.Combo, err = strconv.Atoi(value)
	case "skips":
		ps.Skips, err = strconv.Atoi(value)
	case "field":
		ps.Field, err = newField(value)
	}
	return
}

type gameState game.GameState

func (gs *gameState) processUpdate(name, type_, value string) (winner game.Winner, err error) {
	switch type_ {
	case "this_piece_type":
		gs.ThisPiece = game.Piece(value)
	case "next_piece_type":
		gs.NextPiece = game.Piece(value)
	case "this_piece_position":
		gs.ThisPiecePosition, err = newPosition(value)
	case "winner":
		if value == name {
			winner = game.Me
		} else {
			winner = game.You
		}
	}
	return
}

// serializeMoves returns the textual representation of the moves
func serializeMoves(mvs *[]game.Move) (str string, err error) {
	mvToStr := map[game.Move]string{
		game.Down:      "down",
		game.Left:      "left",
		game.Right:     "right",
		game.TurnLeft:  "turnleft",
		game.TurnRight: "turnright",
		game.Skip:      "skip",
	}
	strs := []string{}
	for _, mv := range *mvs {
		strs = append(strs, mvToStr[mv])
	}
	return strings.Join(strs, ","), nil
}

// newField parses a string like `[[c,...];...]` into a field
func newField(s string) (f *game.Field, err error) {
	f = &game.Field{}
	lines := strings.Split(s, ";")
	for row, line := range lines {
		cells := strings.Split(line, ",")
		for column, cell := range cells {
			var c int
			c, err = strconv.Atoi(cell)
			if err != nil {
				return
			}
			f[row][column] = game.Cell(c)

		}
	}
	return
}
