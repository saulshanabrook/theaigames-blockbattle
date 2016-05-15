package game

// Action represents one way we can act every time the engine asks us for something
type Action struct {
	Moves    []Move
	Location Location
	IsSkip   bool
}
