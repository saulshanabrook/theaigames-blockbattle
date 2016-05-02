package game

// Position represents a certain coordinate on the board, as an Column, Row pair
type Position struct {
	Row    int
	Column int
}

func (p Position) isValid() bool {
	return (0 <= p.Row) && (p.Row < numRows) &&
		(0 <= p.Column) && (p.Column < numColumns)
}

func (p Position) add(p0 Position) Position {
	return Position{p.Row + p0.Row, p.Column + p0.Column}
}

func (p Position) invert() Position {
	return Position{-p.Row, -p.Column}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
func (p Position) distance(p0 Position) int {
	return abs(p.Row-p0.Row) + abs(p.Column-p0.Column)
}
