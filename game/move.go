package game

// Move is an action the player can take
type Move int

const (
	Down Move = iota
	Left
	Right
	TurnLeft
	TurnRight
	Skip
)

func AllMoves() []Move {
	return []Move{Down, Left, Right, TurnLeft, TurnRight, Skip}
}
