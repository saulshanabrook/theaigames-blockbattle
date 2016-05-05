package main

import "github.com/saulshanabrook/blockbattle/rl/learn"

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	l := learn.NewLearner(learn.DefaultLearnerConfig)
	handleErr(l.RunEpisodes(100))
	l.Persist("bots/process/nn")
	return
}
