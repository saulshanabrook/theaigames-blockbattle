package main

import (
	"github.com/saulshanabrook/blockbattle/rl/bot"
	"github.com/saulshanabrook/blockbattle/rl/learn"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// logrus.SetLevel(logrus.DebugLevel)

	l := learn.NewLearner(bot.NewNetwork())
	handleErr(l.RunEpisodesSingle(10000))
	l.Persist("bots/process/nn")

	return
}
