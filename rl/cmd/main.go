package main

import "github.com/saulshanabrook/blockbattle/rl/learn"

func main() {
	l := learn.NewLearner()
	l.RunEpisodes(1000)
	return
}
