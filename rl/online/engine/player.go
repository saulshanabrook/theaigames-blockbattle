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

func NewPlayers() <-chan *player.Player {
	ps := make(chan *player.Player, 2)
	go func() {
		for {
			twoP, err := newPlayers()
			handleErr(err)
			for _, p := range twoP {
				ps <- &p
			}
		}
	}()
	return ps
}

// NewPlayers returns two players that are play against each other
// it runs the java code to start the game and hooks up the player
// inputs and outputs to the java process through intermediary files
//
// When you get the last state (which should have the winner of the game)
// you should close the `Moves` channel. This will trigget the cleanup.
func newPlayers() ([2]player.Player, error) {
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
			States: player.Parse(closeOnEnd(lineStrings(tailLines(inF)))),
			Moves:  mvss,
			Done:   cleanup(inF, cleanup(outF, writeFileChan(outF, player.Serialize(mvss)))),
		}
	}
	go func() {
		handleErr(engineCmd(pCmds).Run())
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
		"./rl/online/engine/pipe.bash %v %v",
		out.Name(),
		in.Name(),
	)
}

func closeOnEnd(xs <-chan string) <-chan string {
	xsp := make(chan string)
	go func() {
		defer close(xsp)
		for x := range xs {
			if x == "END" {
				return
			}
			xsp <- x
		}

	}()
	return xsp
}

func cleanup(f *os.File, done <-chan interface{}) <-chan interface{} {
	doneN := make(chan interface{})
	go func() {
		defer close(doneN)
		<-done
		handleErr(cleanupFile(f))
	}()
	return doneN
}

func engineCmd(pCmds [2]string) *exec.Cmd {
	cmd := exec.Command(
		"java",
		"-cp",
		"rl/online/engine/javac",
		"com.theaigames.blockbattle.Blockbattle",
		pCmds[0],
		pCmds[1],
	)
	// cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func cleanupFile(f *os.File) error {
	err := f.Close()
	if err != nil {
		return err
	}
	return os.Remove(f.Name())
}

func tailLines(file *os.File) <-chan *tail.Line {
	t, err := tail.TailFile(
		file.Name(),
		tail.Config{
			Follow: true,
			Logger: tail.DiscardingLogger,
		},
	)
	handleErr(err)
	return t.Lines
}

func lineStrings(lines <-chan *tail.Line) <-chan string {
	strings := make(chan string)
	go func() {
		defer close(strings)
		for line := range lines {
			handleErr(line.Err)
			strings <- line.Text
		}
	}()
	return strings
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
	return
}

func writeFileChan(file *os.File, lines <-chan string) (done <-chan interface{}) {
	doneRW := make(chan interface{})
	go func() {
		defer close(doneRW)
		for line := range lines {
			_, err := file.WriteString(line + "\n")
			if err != nil {
				panic(err)
			}
		}
	}()
	return doneRW
}
