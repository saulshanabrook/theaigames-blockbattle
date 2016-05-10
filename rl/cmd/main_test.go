package main

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/saulshanabrook/blockbattle/bots"
	"github.com/saulshanabrook/blockbattle/rl/bot"
	"github.com/saulshanabrook/blockbattle/rl/engine"
)

func BenchmarkRandomGame(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ps, _ := engine.NewPlayers()
		go bots.Play(bots.Random{}, ps[0])
		bots.Play(bots.Random{}, ps[1])
	}
}

func BenchmarkRandomActionGame(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ps, _ := engine.NewPlayers()
		go bots.Play(bots.RandomAction{}, ps[0])
		bots.Play(bots.RandomAction{}, ps[1])
	}
}

func BenchmarkRNNGame(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b := &bot.Bot{Calculator: bot.NewNetwork()}
		ps, _ := engine.NewPlayers()
		go bots.Play(b, ps[0])
		bots.Play(b, ps[1])
	}
}

func init() {
	logrus.SetLevel(logrus.WarnLevel)

	os.Chdir("../..")
}
