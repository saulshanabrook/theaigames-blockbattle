package main

import (
	"bufio"
	"os"

	"github.com/saulshanabrook/blockbattle/player"
)

// NewPlayer returns a player that uses stdin and stdout out to communicate
func NewPlayer() *player.Player {
	return &player.Player{
		Input:  readStdinChan(),
		Output: player.WriteFileChan(os.Stdout),
	}
}

func readStdinChan() <-chan string {
	lines := make(chan string)
	os.Stderr.WriteString("outside go func")
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		os.Stderr.WriteString("In go func")
		for scanner.Scan() {
			txt := scanner.Text()
			os.Stderr.WriteString("got: " + txt)
			lines <- txt
		}
	}()
	return lines
}
