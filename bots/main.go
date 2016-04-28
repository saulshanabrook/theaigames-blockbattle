package bots

import (
	"github.com/saulshanabrook/blockbattle/game"
	"github.com/saulshanabrook/blockbattle/player"
)

// Bot is anything that can play a game. All it has to do is take in a game
// state and return a list of moves to perform
// It also has to handle when the game finishes, if it wants to do something with
// this state
type Bot interface {
	Act(*game.State) *[]game.Move
}

// Play starts using the bot to play a player
func Play(b Bot, p *player.Player) {
	sts, _, mvss := p.Process()
	for st := range sts {
		mvss <- b.Act(st)
	}
}
