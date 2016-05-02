package game

// Location is a certain position and rotation that a piece can be in
type Location struct {
	Position Position
	Rotation Rotation
}

func (l Location) add(l0 Location) Location {
	return Location{
		l.Position.add(l0.Position),
		l.Rotation.add(l0.Rotation),
	}
}

func (l Location) invert() Location {
	return Location{
		l.Position.invert(),
		l.Rotation.invert(),
	}
}

func (l Location) distance(l0 Location) int {
	return l.Position.distance(l0.Position) + l.Rotation.distance(l0.Rotation)
}
