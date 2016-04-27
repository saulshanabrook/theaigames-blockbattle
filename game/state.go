package game

// PlayerState holds all information about one bot
type PlayerState struct {
	RowPoints int
	Combo     int
	Skips     int
	Field     *Field
}

// GameState holds current round information that pertains to both players
type GameState struct {
	Winner            Winner
	ThisPiece         Piece
	NextPiece         Piece
	ThisPiecePosition *Position
}

// State is a representation of whole current state
type State struct {
	Name  string
	Game  *GameState
	Mine  *PlayerState
	Yours *PlayerState
}

// Winner is who won the game
type Winner int

const (
	// None means no one has one yet
	None Winner = iota
	// You means the other bot one
	You
	// Me means this bot has won
	Me
)

// NewState returns a blank state with zero value
func NewState() (st *State) {
	st = &State{}
	st.Game = &GameState{}
	st.Mine = newPlayerState()
	st.Yours = newPlayerState()
	return
}

func newPlayerState() (ps *PlayerState) {
	ps = &PlayerState{}
	ps.Field = &Field{}
	return
}
