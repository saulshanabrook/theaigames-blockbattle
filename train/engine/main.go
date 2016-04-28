package main

import (
	"sync"

	"github.com/saulshanabrook/blockbattle/bots"
	"github.com/saulshanabrook/blockbattle/player"
)

func main() {
	var b bots.Bot = &bots.Random{}
	players, err := NewPlayers()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, p := range players {
		wg.Add(1)
		go func(p *player.Player) {
			defer wg.Done()
			bots.Play(b, p)
		}(p)
	}
	wg.Wait()
}
