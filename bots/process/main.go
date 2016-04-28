package main

import (
	"github.com/saulshanabrook/blockbattle/bots"
	"github.com/saulshanabrook/blockbattle/player"
)

func main() {
	var bot bots.Bot = &bots.Random{}
	player := player.NewUsingProcess()
	bots.Play(bot, player)
}
