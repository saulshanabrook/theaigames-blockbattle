// Package dqn support generic Deep Q Learning, as developed by
// DeepMind
package dqn

import "github.com/saulshanabrook/blockbattle/game"

type Learner interface {
	Learn(st game.State, a game.Action, value float64)
}

func LearnExperience(l Learner, d Decider, e *Experience) {
	_, bestNextVal := d.Decide(e.NextState)
	val := bestNextVal + e.Reward
	l.Learn(e.State, e.Action, val)
}

func LearnExperiences(l Learner, d Decider, e *Experiences) {
	for {
		LearnExperience(l, d, e.Peek())
	}
}
