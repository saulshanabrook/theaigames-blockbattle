package game

import (
	"bufio"
	"os"
)

// PlayerFiles holds input and output files for a certain player in a game
type PlayerFiles struct {
	output *os.File // engine will read from this file, we write to it
	input  *os.File // engine will write to this file, we read from it
}

// ToPlayer creates a Player based on these files and starts up goroutines
// to keep the channels on that player in sync with these files
func (pf *PlayerFiles) ToPlayer() *Player {
	input := make(chan string)
	go func() {
		scanner := bufio.NewScanner(pf.input)
		for scanner.Scan() {
			input <- scanner.Text()
		}
	}()

	output := make(chan string)
	go func() {
		for msg := range output {
			_, err := pf.output.WriteString(msg + "\n")
			if err != nil {
				panic(err)
			}
		}
	}()
	return &Player{input: input, output: output}
}
