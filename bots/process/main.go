package main

import "github.com/saulshanabrook/blockbattle/bots"

func main() {
	var b bots.Bot = &bots.Random{}
	p := NewPlayer()
	bots.Play(b, p)
}
