package game

type row []Cell
type grid []row

func newGrid(s int) grid {
	g := make(grid, s)
	for r := range g {
		g[r] = make(row, s)
	}
	return g
}

// rotate once clockwise
// http://stackoverflow.com/a/42535/907060
// func (g grid) rotate() grid {
// 	newG := newGrid(len(g))
// 	for i, row := range g {
// 		for j := range row {
// 			newG[i][j] = g[len(g)-j-1][i]
// 		}
// 	}
// 	return newG
// }

// placeOnPosition take in a grid of Empty and Shapes.
// If we call the top left corner of the shape (0, 0),
// then this returns all the relative positions we can put
// put this in so that is resting right above a shape.
//
// It returns all changes in position that cause this grid to move to sit on a
// block
func (g grid) placeOnPositions() []Position {
	ls := make([]Position, 0, 25)
	var prevCells []Cell
	for row, cells := range g {
		for col, cell := range cells {
			if cell == Empty {
				if prevCells != nil && prevCells[col] == Shape {
					ls = append(ls, Position{Column: -col, Row: -row})
				}
			}
		}
		prevCells = cells
	}
	for col, cell := range prevCells {
		if cell == Shape {
			ls = append(ls, Position{Column: -col, Row: -(len(g) - 1)})
		}
	}
	return ls
}
