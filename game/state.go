package game

import (
	"strconv"
	"strings"
)

// Cell is one block in the tetris board
type Cell int

const (
	// Empty means no block
	Empty Cell = iota
	// Shape is part of the moving piece
	Shape
	// Block can be cleared
	Block
	// Solid is at the bottom of the row and cant be changed
	Solid
)

// Field is the whole board, 20 hi and 10 wide
type Field [20][10]Cell

// ParseField parses a string like `[[c,...];...]` into a field
func ParseField(s string) (f *Field, err error) {
	f = &Field{}
	lines := strings.Split(s, ";")
	for row, line := range lines {
		cells := strings.Split(line, ",")
		for column, cell := range cells {
			var c int
			c, err = strconv.Atoi(cell)
			if err != nil {
				return
			}
			f[row][column] = Cell(c)

		}
	}
	return
}

// PlayerState holds all information about one bot
type PlayerState struct {
	RowPoints int
	Combo     int
	Skips     int
	Field     *Field
}

func NewPlayerState() (ps *PlayerState) {
	ps = &PlayerState{}
	ps.Field = &Field{}
	return
}

// Piece is a type of block, one of I, J, L, O, S, T or Z
type Piece string

// Position represents a certain coordinate on the board, from the top left
type Position struct {
	x int
	y int
}

// ParsePosition turns a string like `4,-1` into a position
func ParsePosition(s string) (p *Position, err error) {
	p = &Position{}
	ss := strings.Split(s, ",")
	p.x, err = strconv.Atoi(ss[0])
	if err != nil {
		return
	}
	p.y, err = strconv.Atoi(ss[1])
	return
}

// GameState holds current round information that pertains to both players
type GameState struct {
	Winner            Winner
	ThisPiece         Piece
	NextPiece         Piece
	ThisPiecePosition *Position
}

// State is a representation of whole current state
type State struct {
	Name  string
	Game  *GameState
	Mine  *PlayerState
	Yours *PlayerState
}

func NewState() (st *State) {
	st = &State{}
	st.Game = &GameState{}
	st.Mine = NewPlayerState()
	st.Yours = NewPlayerState()
	return
}

func (ps *PlayerState) processUpdate(type_, value string) (err error) {
	switch type_ {
	case "row_points":
		ps.RowPoints, err = strconv.Atoi(value)
	case "combo":
		ps.Combo, err = strconv.Atoi(value)
	case "skips":
		ps.Skips, err = strconv.Atoi(value)
	case "field":
		ps.Field, err = ParseField(value)
	}
	return
}

// Winner is who won the game
type Winner int

const (
	// none means no one has one yet
	none Winner = iota
	// You means the other bot one
	You
	// Me means this bot has won
	Me
)

func (gs *GameState) processUpdate(name, type_, value string) (winner Winner, err error) {
	switch type_ {
	case "this_piece_type":
		gs.ThisPiece = Piece(value)
	case "next_piece_type":
		gs.NextPiece = Piece(value)
	case "this_piece_position":
		gs.ThisPiecePosition, err = ParsePosition(value)
	case "winner":
		if value == name {
			winner = Me
		} else {
			winner = You
		}
	}
	return
}

func (s *State) processSettings(type_, value string) (err error) {
	switch type_ {
	case "your_bot":
		s.Name = value
	}
	return
}

func (s *State) processUpdate(player, type_, value string) (Winner, error) {
	if player == "game" {
		return s.Game.processUpdate(s.Name, type_, value)
	}
	if player == s.Name {
		return none, s.Mine.processUpdate(type_, value)
	}
	return none, s.Yours.processUpdate(type_, value)
}

func (s *State) processLine(line string) (gotAction bool, winner Winner, err error) {
	parts := strings.Split(line, " ")
	cmd, rest := parts[0], parts[1:]
	switch cmd {
	case "settings":
		return false, none, s.processSettings(rest[0], rest[1])
	case "update":
		winner, err = s.processUpdate(rest[0], rest[1], rest[2])
		return false, winner, err
	case "action":
		return true, none, nil
	}
	panic(line)
}
