package game

import (
	"errors"
	"math"
	"reflect"
)

const numColumns = 10
const numRows = 20

// Field is the whole board, 20 high and 10 wide
type Field [numRows][numColumns]Cell

// Actions gives all the final location you can move a piece to, and the
// moves required to get it there
func (f Field) Actions(p Piece, startPos Position) map[Location][]Move {
	as := make(map[Location][]Move)

	// we also keep a set of the failed actions, so we don't try them twice
	asFailed := make(map[Location]interface{})

	restingOffsetLocs := p.restingOffsetLocations()
	// To figure out all the possible locations this piece could end up,
	// we first look at all positions on the field a block could be "frozen"
	// These are the `restingPos`s, or all the positions above an existing block
	// or the bottom
	for _, restingPos := range f.restingPositions() {
		// Then we go through each way a piece could be offset from a resting
		// position so that it would be sitting on that position
		for _, restingOffsetLoc := range restingOffsetLocs {
			// we then combine the offset with the position to get the final
			// block position where it is at least resting on an existing piece
			loc := restingOffsetLoc.add(Location{restingPos, 0})

			// it is possible we get duplicate locations, so we keep a record
			// of all the ones we have tried
			_, found := as[loc]
			_, failed := asFailed[loc]
			tried := found || failed
			if !tried {
				// then we ask to get the valid moves for this piece and start position
				moves, err := f.movesFrom(p, startPos, loc)
				if err == nil {
					as[loc] = moves
				} else {
					asFailed[loc] = nil
				}
			}
		}
	}
	return as
}

type nodeInfo struct {
	mvs   []Move
	cost  int
	hCost int
}

func (ni nodeInfo) addMv(mv Move) nodeInfo {
	return nodeInfo{
		mvs:  append([]Move{mv}, ni.mvs...),
		cost: ni.cost + 1,
	}
}
func (ni nodeInfo) totalCost() int {
	return ni.cost + ni.hCost
}

func (f Field) movesFrom(p Piece, startPos Position, endLoc Location) (mvs []Move, err error) {
	validL := func(l Location) bool {
		return f.validPlacement(p, l)
	}

	startLocation := Location{Position: startPos}
	if !validL(startLocation) {
		return nil, errors.New("start location not valid")
	}
	if !validL(endLoc) {
		return nil, errors.New("end location not valid")
	}
	// we actually want to search from end to start, instead of start to end
	// because this will hopefully fail faster if there is no possible
	// path

	// these keep track of if we start at the `endLoc` what moves we need to take
	// to get to another location and the A* cost to make those moves
	// to get from the `endLoc` to itself takes 0 moves and costs 0
	currentLs := map[Location]nodeInfo{}

	heuristicCost := func(l Location) int {
		return startLocation.distance(l)
	}
	shouldVisit := func(l Location, ni nodeInfo) []Move {
		if reflect.DeepEqual(l, startLocation) {
			return ni.mvs
		}
		ni.hCost = heuristicCost(l)
		currentLs[l] = ni
		return nil
	}

	if mvs = shouldVisit(endLoc, nodeInfo{mvs: []Move{}, cost: 0}); mvs != nil {
		return
	}

	// a set of the locations we have already visited
	visitedLs := make(map[Location]interface{})

	popSmallestCost := func() (minL Location, minNI nodeInfo) {
		minC := math.MaxInt8
		for l, ni := range currentLs {
			c := ni.totalCost()
			if c < minC {
				minC = c
				minNI = ni
				minL = l
			}
		}
		delete(currentLs, minL)
		return
	}

	neighbors := func(l Location) map[Location]Move {
		lmvs := make(map[Location]Move)
		for mv, offsetL := range moveLocationDiffs {
			// we actually want to do the inverse of the move because we are searching
			// bottom up
			finalL := l.add(offsetL.invert())
			if validL(finalL) {
				lmvs[finalL] = mv
			}
		}
		return lmvs
	}

	for len(currentLs) != 0 {
		currentL, currentN := popSmallestCost()
		visitedLs[currentL] = nil
		for newLoc, mv := range neighbors(currentL) {
			_, visited := visitedLs[newLoc]
			_, current := currentLs[newLoc]
			if visited || current {
				continue
			}
			if mvs = shouldVisit(newLoc, currentN.addMv(mv)); mvs != nil {
				return
			}
		}
	}
	return nil, errors.New("Couldnt find path")
}

func (f Field) validPlacement(p Piece, l Location) bool {
	g, exists := grids[p][l.Rotation]
	if !exists {
		return false
	}
	for row, cells := range g {
		for col, cell := range cells {
			// we want to check all the shape blocks to make sure they
			// are either right above the board (with a row of -1) or they are in
			// a blank or a shape space
			if cell == Shape {
				cellPos := l.Position.add(Position{Row: row, Column: col})
				cellRow, cellCol := cellPos.Row, cellPos.Column
				validCol := 0 <= cellCol && cellCol < numColumns
				if !validCol {
					return false
				}
				if cellRow == -1 {
					continue
				}
				validRow := 0 <= cellRow && cellRow < numRows
				if !validRow {
					return false
				}
				fCell := f[cellRow][cellCol]
				if fCell == Block || fCell == Solid {
					return false
				}
			}
		}
	}
	return true
}

// restingPositions returns all the possible positions where
// a single block could sit on the field.
//
// These are all the Empty or Shape cells that have a Block, Solid, or the
// bottom below them
func (f Field) restingPositions() []Position {
	ps := make([]Position, 0, 10)

	// go through the rows from bottom to top.

	cellsBelow := [numColumns]Cell{Solid, Solid, Solid, Solid, Solid, Solid, Solid, Solid, Solid, Solid}
	for row := len(f) - 1; row >= 0; row-- {
		cells := f[row]
		for col, cell := range cells {
			// If the row below the current row has a block at this index
			// (or it doesnt exist and we are at the bottom row) then we know this position
			// is an OK place to put a block
			if cell == Empty || cell == Shape {
				if cellsBelow[col] == Block || cellsBelow[col] == Solid {
					ps = append(ps, Position{Row: row, Column: col})
				}
			}
		}
		cellsBelow = cells
	}
	return ps
}
