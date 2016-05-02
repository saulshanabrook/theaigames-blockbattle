package main

// This file starts a local engine.
//
// It creates two files for each player:
// input: the engine writes to this and we read from it
// output: the engine reads from this and we write to it
//
// It then tells the engine to start with a bash script
// that will take its `stdin` and write it to `input`
// and reads `output` to its stdout

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/hpcloud/tail"
	"github.com/saulshanabrook/blockbattle/player"
)

// pFiles holds input and output files for a certain player in a game
type playerFiles struct {
	input  *os.File
	output *os.File
}

func newPlayerFiles() (pf *playerFiles, err error) {
	pf = &playerFiles{}
	pf.input, err = ioutil.TempFile("", "input")
	if err != nil {
		return
	}
	pf.output, err = ioutil.TempFile("", "output")
	return
}

// ToPlayer creates a Player based on these files and starts up goroutines
// to keep the channels on that player in sync with these files
func (pf *playerFiles) toPlayer() *player.Player {
	return &player.Player{
		Input:  readFileChan(pf.input),
		Output: player.WriteFileChan(pf.output),
	}
}

func (pf *playerFiles) command() string {
	return fmt.Sprintf(
		"./train/engine/pipe.bash %v %v",
		pf.output.Name(),
		pf.input.Name(),
	)
}

func (pf *playerFiles) cleanup() {
	err := cleanupFile(pf.input)
	if err != nil {
		panic(err)
	}
	err = cleanupFile(pf.output)
	if err != nil {
		panic(err)
	}
	return
}

// NewPlayers returns two players that are play against each other
// it runs the java code to start the game and hooks up the player
// inputs and outputs to the java process through intermediary files
func NewPlayers() (ps [2]*player.Player, err error) {
	ps = [2]*player.Player{}
	pfss := [2]*playerFiles{}
	pCmds := [2]string{}
	for i := range pfss {
		pfss[i], err = newPlayerFiles()
		if err != nil {
			return
		}

		ps[i] = pfss[i].toPlayer()

		pCmds[i] = pfss[i].command()
	}
	go func() {
		defer pfss[0].cleanup()
		defer pfss[1].cleanup()
		err = startEngine(pCmds)
		if err != nil {
			panic(err)
		}
	}()
	return
}

func startEngine(pCmds [2]string) error {
	cmd := exec.Command(
		"java",
		"-cp",
		"train/engine/javac",
		"com.theaigames.blockbattle.Blockbattle",
		pCmds[0],
		pCmds[1],
	)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	return cmd.Run()
}

func cleanupFile(f *os.File) (err error) {
	err = f.Close()
	if err != nil {
		return err
	}
	return os.Remove(f.Name())
}

// readFileChan takes in a file and returns a channel that sends each
// line in the file
func readFileChan(file *os.File) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		t, err := tail.TailFile(
			file.Name(),
			tail.Config{Follow: true})
		if err != nil {
			panic(err)
		}
		for line := range t.Lines {
			lineT := line.Text
			lines <- lineT
		}
	}()
	return lines
}