package game

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

// Field is the whole board, 20 high and 10 wide
type Field [20][10]Cell

// Piece is a type of block, one of I, J, L, O, S, T or Z
type Piece string

// Position represents a certain coordinate on the board, from the top left
type Position struct {
	X int
	Y int
}
