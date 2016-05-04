package main

import "github.com/saulshanabrook/blockbattle/rl/learn"

func main() {
	l := learn.NewLearner(learn.DefaultLearnerConfig)
	l.RunEpisodes(1)
	l.Persist("bots/process/nn")
	return
}
