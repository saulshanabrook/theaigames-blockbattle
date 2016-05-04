package main

import (
	"github.com/saulshanabrook/blockbattle/bots"
	"github.com/saulshanabrook/blockbattle/rl/bot"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	bytes, err := Asset("bots/process/nn")
	handleErr(err)
	b, err := bot.NewFromBinary(bytes)
	handleErr(err)
	bots.Play(b, NewPlayer())
}
