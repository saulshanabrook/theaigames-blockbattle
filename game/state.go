package game

// PlayerState holds all information about one bot
type PlayerState struct {
	RowPoints int
	Combo     int
	Skips     int
	Field     Field
}

// GameState holds current round information that pertains to both players
type GameState struct {
	Winner            Winner `json:"-"`
	ThisPiece         Piece
	NextPiece         Piece
	ThisPiecePosition Position
}

// State is a representation of whole current state
type State struct {
	Name  string `json:"-"`
	Game  GameState
	Mine  PlayerState
	Yours PlayerState
}

// IsOver returns if there is a winner
func (s State) IsOver() bool {
	return s.Game.Winner != None
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
	// Tie means neither of us one
	Tie
)

// NewState returns a blank state with zero value
func NewState() State {
	return State{
		Game:  GameState{},
		Mine:  newPlayerState(),
		Yours: newPlayerState(),
	}
}

// Actions returns a list locations and moves you can make during this game
func (s *State) Actions() []Action {
	if s.IsOver() {
		panic("Game is over")
	}
	fActs := s.Mine.Field.Actions(s.Game.ThisPiece, s.Game.ThisPiecePosition)
	acts := make([]Action, 0, len(fActs)+1)
	for loc, mvs := range fActs {
		acts = append(acts, Action{Moves: mvs, Location: loc})
	}
	if s.Mine.Skips > 0 {
		acts = append(acts, Action{Moves: []Move{MoveSkip}, IsSkip: true})
	}
	if len(acts) == 0 {
		// we cant move because we will lose
		acts = append(
			acts,
			Action{
				Moves:    []Move{},
				Location: Location{Position: s.Game.ThisPiecePosition},
			},
		)
	}
	return acts
}

func newPlayerState() PlayerState {
	ps := PlayerState{}
	ps.Field = Field{}
	return ps
}
