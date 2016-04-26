package game

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

// LocalEngine is on your machine, by starting a server with the java code
type LocalEngine struct {
	pFiles [2]*PlayerFiles
}

// MakeLocalEngine returns a new `LocalEngine`, creating temporary files to
// use to communicate with the game once it starts
func MakeLocalEngine() (g *LocalEngine, err error) {
	g = &LocalEngine{}
	for i := range g.pFiles {
		files := PlayerFiles{}
		g.pFiles[i] = &files
		files.input, err = ioutil.TempFile("", "input")
		if err != nil {
			return
		}
		files.output, err = ioutil.TempFile("", "output")
		if err != nil {
			return
		}
	}
	return
}

func playerFilestoCommand(pf *PlayerFiles) string {
	return fmt.Sprintf(
		"./pipe.bash %v %v",
		pf.output.Name(),
		pf.input.Name(),
	)
}

func (g *LocalEngine) startEngine() {
	go func() {
		err := exec.Command(
			"java",
			"-cp",
			"bin",
			"com.theaigames.blockbattle.Blockbattle",
			playerFilestoCommand(g.pFiles[0]),
			playerFilestoCommand(g.pFiles[1]),
		).Run()
		if err != nil {
			panic(err)
		}
	}()
}

// Start runs the game in the background and returns players to interact with it
//
// it runs the java code to start the game and hooks up the player
// inputs and outputs to the java process through intermediary files
func (g *LocalEngine) Start() [2]*Player {
	ps := [2]*Player{&Player{}, &Player{}}
	for i, pFile := range g.pFiles {
		ps[i] = pFile.ToPlayer()
	}
	g.startEngine()
	return ps
}
