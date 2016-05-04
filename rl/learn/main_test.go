package learn

import (
	"os"
	"testing"
)

func BenchmarkTrain(b *testing.B) {
	l := NewLearner(LearnerConfig{MinibatchSize: 0, DiscountFactor: 0.9})
	l.RunEpisodes(1)
	l.c.MinibatchSize = b.N
	b.ResetTimer()
	l.trainMinibatch()
}

func init() {
	os.Chdir("../..")
}
