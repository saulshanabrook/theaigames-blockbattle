package player

import "os"

// NewPlayerUsingProcess returns a player that uses stdin and stdout out to communicate
func NewPlayerUsingProcess() *Player {
	return &Player{
		input:  readFileChan(os.Stdin),
		output: writeFileChan(os.Stdout),
	}
}
