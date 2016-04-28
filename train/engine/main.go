package main

import (
	"sync"

	"github.com/saulshanabrook/blockbattle/bots"
	"github.com/saulshanabrook/blockbattle/player"
)

func main() {
	var b bots.Bot = &bots.Random{}
	players, err := player.NewUsingEngine()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, player := range players {
		wg.Add(1)
		go func(p *player.Player) {
			defer wg.Done()
			bots.Play(b, p)
		}(player)
	}
	wg.Wait()
}
