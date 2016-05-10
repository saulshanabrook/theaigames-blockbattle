package learn

import (
	"fmt"
	"os"
	"testing"

	"github.com/saulshanabrook/blockbattle/rl/bot"
)

func benchRunEpisode(b *testing.B, p float64) {
	net := bot.NewNetwork()
	l := NewLearner(net)
	bot_ := bot.Bot{Calculator: net}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		l.RunEpisode(&bot_, p)
	}
}
func BenchmarkRunEpisodeLowP(b *testing.B) {
	benchRunEpisode(b, 0)
}

func BenchmarkRunEpisodeHighP(b *testing.B) {
	benchRunEpisode(b, 1)
}

func BenchmarkRunEpisodesSingle(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		l := NewLearner(bot.NewNetwork())
		b.StartTimer()
		l.RunEpisodesSingle(16)
		fmt.Println(l.nTrains)
	}
}

func benchRunEpisodesWorkers(b *testing.B, nWorkers int) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		l := NewLearner(bot.NewNetwork())
		b.StartTimer()
		l.RunEpisodes(16, nWorkers)
		fmt.Println(l.nTrains)
	}
}

func BenchmarkRunEpisodes1Worker(b *testing.B) {
	benchRunEpisodesWorkers(b, 1)
}

func BenchmarkRunEpisodes2Worker(b *testing.B) {
	benchRunEpisodesWorkers(b, 2)
}

func BenchmarkRunEpisodes3Worker(b *testing.B) {
	benchRunEpisodesWorkers(b, 3)
}

func BenchmarkRunEpisodes4Worker(b *testing.B) {
	benchRunEpisodesWorkers(b, 4)
}

func init() {
	// logrus.SetLevel(logrus.WarnLevel)

	os.Chdir("../..")
}
