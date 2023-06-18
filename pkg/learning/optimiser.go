package learning

import (
	"sync"
)

// LearningOptimiser
type LearningOptimiser struct {
	Learners []Learner
}

func (l *LearningOptimiser) Run() {
	for _, learner := range l.Learners {
		learner.Run(&sync.WaitGroup{})
	}
}
