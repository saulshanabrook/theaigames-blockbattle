package main

import (
	"github.com/saulshanabrook/blockbattle/bots"
	"github.com/saulshanabrook/blockbattle/game"
)

// func main() {
// 	var bot bots.Bot = &bots.Random{}
// 	players, err := game.StartLocalGame()
// 	if err != nil {
// 		panic(err)
// 	}
// 	var wg sync.WaitGroup
// 	for _, player := range players {
// 		wg.Add(1)
// 		go func(player *game.Player) {
// 			defer wg.Done()
// 			bots.Play(bot, player)
// 		}(player)
// 	}
// 	wg.Wait()
// }

func main() {
	var bot bots.Bot = &bots.Random{}
	player := game.StartLocalPlayer()
	bots.Play(bot, player)
}
