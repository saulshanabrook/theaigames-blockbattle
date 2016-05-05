package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/saulshanabrook/blockbattle/game"
	"github.com/saulshanabrook/blockbattle/player"
)

// NewPlayer returns a player that uses stdin and stdout out to communicate
func NewPlayer() player.Player {
	mvs := make(chan []game.Move)
	go writeStdinChan(player.Serialize(mvs))
	return player.Player{
		States: player.Parse(readStdinChan()),
		Moves:  mvs,
		Done:   writeStdinChan(player.Serialize(mvs)),
	}
}

func readStdinChan() <-chan string {
	lines := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()
	return lines
}

func writeStdinChan(lines <-chan string) <-chan interface{} {
	done := make(chan interface{})
	go func() {
		for line := range lines {
			fmt.Println(line)
		}
	}()
	return done
}
