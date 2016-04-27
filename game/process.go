package game

import "os"

// StartLocalPlayer returns a player that uses stdin and stdout out to communicate
func StartLocalPlayer() *Player {
	return &Player{
		input:  readFileChan(os.Stdin),
		output: writeFileChan(os.Stdout),
	}
}
