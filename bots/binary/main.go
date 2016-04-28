package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	data, err := Asset("bots/process/main")
	checkErr(err)
	filename, err := filepath.Abs("./bot")
	checkErr(err)
	checkErr(ioutil.WriteFile(filename, data, os.ModePerm))

	cmd := exec.Command(filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	checkErr(cmd.Run())
	return
}
