package game

// Rotation is how a piece is rotated
type Rotation int

const (
	// RotatedUp is the default orientation
	RotatedUp Rotation = iota
	// RotatedRight is it rotated once to the right
	RotatedRight
	// RotatedDown is it rotated twice to the left or right
	RotatedDown
	// RotatedLeft is it rotated once to the left
	RotatedLeft
)

const numRotations = 4

var allRotations = []Rotation{RotatedUp, RotatedRight, RotatedDown, RotatedLeft}

func (r Rotation) add(r0 Rotation) Rotation {
	return Rotation((int(r) + int(r0)) % numRotations)
}

// invert should return an additive inverse of the rotation
func (r Rotation) invert() Rotation {
	return Rotation((int(r) * -1) % numRotations)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// gets the distance modulo numRotations
func (r Rotation) distance(r0 Rotation) int {
	return min(
		(int(r)-int(r0)+numRotations)%numRotations,
		(int(r0)-int(r)+numRotations)%numRotations,
	)
}
