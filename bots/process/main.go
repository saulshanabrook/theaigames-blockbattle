package main

import "github.com/saulshanabrook/blockbattle/bots"

func main() {
	var b = bots.RandomAction{}
	p := NewPlayer()
	bots.Play(b, p)
}
