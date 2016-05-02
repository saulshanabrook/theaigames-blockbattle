package game

// Piece is a type of block
type Piece string

const (
	I Piece = "I"
	J       = "J"
	L       = "L"
	O       = "O"
	S       = "S"
	T       = "T"
	Z       = "Z"
)

var grids = map[Piece]map[Rotation]grid{
	I: {
		RotatedUp: {
			{0, 0, 0, 0},
			{1, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		RotatedRight: {
			{0, 0, 1, 0},
			{0, 0, 1, 0},
			{0, 0, 1, 0},
			{0, 0, 1, 0},
		},
		RotatedDown: {
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{1, 1, 1, 1},
			{0, 0, 0, 0},
		},
		RotatedLeft: {
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 0, 0},
		}},
	J: {
		RotatedUp: {
			{1, 0, 0},
			{1, 1, 1},
			{0, 0, 0},
		},
		RotatedRight: {
			{0, 1, 1},
			{0, 1, 0},
			{0, 1, 0},
		},
		RotatedDown: {
			{0, 0, 0},
			{1, 1, 1},
			{0, 0, 1},
		},
		RotatedLeft: {
			{0, 1, 0},
			{0, 1, 0},
			{1, 1, 0},
		}},
	L: {
		RotatedUp: {
			{0, 0, 1},
			{1, 1, 1},
			{0, 0, 0},
		},
		RotatedRight: {
			{0, 1, 0},
			{0, 1, 0},
			{0, 1, 1},
		},
		RotatedDown: {
			{0, 0, 0},
			{1, 1, 1},
			{1, 0, 0},
		},
		RotatedLeft: {
			{1, 1, 0},
			{0, 1, 0},
			{0, 1, 0},
		}},
	O: {
		RotatedUp: {
			{1, 1},
			{1, 1},
		},
	},
	S: {
		RotatedUp: {
			{0, 1, 1},
			{1, 1, 0},
			{0, 0, 0},
		},
		RotatedRight: {
			{0, 1, 0},
			{0, 1, 1},
			{0, 0, 1},
		},
		RotatedDown: {
			{0, 0, 0},
			{0, 1, 1},
			{1, 1, 0},
		},
		RotatedLeft: {
			{1, 0, 0},
			{1, 1, 0},
			{0, 1, 0},
		}},
	T: {
		RotatedUp: {
			{0, 1, 0},
			{1, 1, 1},
			{0, 0, 0},
		},
		RotatedRight: {
			{0, 1, 0},
			{0, 1, 1},
			{0, 1, 0},
		},
		RotatedDown: {
			{0, 0, 0},
			{1, 1, 1},
			{0, 1, 0},
		},
		RotatedLeft: {
			{0, 1, 0},
			{1, 1, 0},
			{0, 1, 0},
		}},
	Z: {
		RotatedUp: {
			{1, 1, 0},
			{0, 1, 1},
			{0, 0, 0},
		},
		RotatedRight: {
			{0, 0, 1},
			{0, 1, 1},
			{0, 1, 0},
		},
		RotatedDown: {
			{0, 0, 0},
			{1, 1, 0},
			{0, 1, 1},
		},
		RotatedLeft: {
			{0, 1, 0},
			{1, 1, 0},
			{1, 0, 0},
		}},
}

func (p Piece) restingOffsetLocations() []Location {
	ls := make([]Location, 0, 10)
	for rot, g := range grids[p] {
		for _, pos := range g.placeOnPositions() {
			ls = append(ls, Location{Position: pos, Rotation: rot})
		}
	}
	return ls
}
