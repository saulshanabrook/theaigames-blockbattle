package game

// Move is an action the player can take
type Move int

const (
	MoveDown Move = iota
	MoveLeft
	MoveRight
	MoveTurnLeft
	MoveTurnRight
	MoveSkip
	MoveDrop
)

var AllMoves = []Move{MoveDown, MoveLeft, MoveRight, MoveTurnLeft, MoveTurnRight, MoveSkip, MoveDrop}

var moveLocationDiffs = map[Move]*Location{
	MoveDown:      {Position: Position{Row: 1}},
	MoveLeft:      {Position: Position{Column: -1}},
	MoveRight:     {Position: Position{Column: 1}},
	MoveTurnLeft:  {Rotation: RotatedLeft},
	MoveTurnRight: {Rotation: RotatedRight},
}
