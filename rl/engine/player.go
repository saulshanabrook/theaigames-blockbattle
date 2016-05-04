package engine

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
	"github.com/saulshanabrook/blockbattle/game"
	"github.com/saulshanabrook/blockbattle/player"
)

// NewPlayers returns two players that are play against each other
// it runs the java code to start the game and hooks up the player
// inputs and outputs to the java process through intermediary files
//
// When you get the last state (which should have the winner of the game)
// you should close the `Moves` channel. This will trigget the cleanup.
func NewPlayers() ([2]player.Player, error) {
	ps := [2]player.Player{}
	pCmds := [2]string{}
	for i := range ps {
		inF, outF, err := newFiles()
		if err != nil {
			return ps, err
		}
		pCmds[i] = command(inF, outF)

		mvss := make(chan []game.Move)
		ps[i] = player.Player{
			States: player.Parse(readFileChan(inF)),
			Moves:  mvss,
		}
		cleanup(inF, cleanup(outF, player.WriteFileChan(outF, player.Serialize(mvss))))
	}
	go func() {
		handleErr(startEngine(pCmds))
	}()
	return ps, nil
}

func newFiles() (in, out *os.File, err error) {
	in, err = ioutil.TempFile("", "input")
	if err != nil {
		return
	}
	out, err = ioutil.TempFile("", "output")
	return
}

func command(in, out *os.File) string {
	return fmt.Sprintf(
		"./rl/engine/pipe.bash %v %v",
		out.Name(),
		in.Name(),
	)
}

func cleanup(f *os.File, xs <-chan string) <-chan string {
	xsp := make(chan string)
	go func() {
		for x := range xs {
			xsp <- x
		}
		close(xsp)
		handleErr(cleanupFile(f))
	}()
	return xsp
}

func startEngine(pCmds [2]string) error {
	cmd := exec.Command(
		"java",
		"-cp",
		"rl/engine/javac",
		"com.theaigames.blockbattle.Blockbattle",
		pCmds[0],
		pCmds[1],
	)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	return cmd.Run()
}

func cleanupFile(f *os.File) error {
	err := f.Close()
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
			tail.Config{
				Follow: true,
				Logger: tail.DiscardingLogger,
			},
		)
		if err != nil {
			panic(err)
		}
		for line := range t.Lines {
			lines <- line.Text
		}

	}()
	return lines
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
	return
}
