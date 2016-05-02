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
