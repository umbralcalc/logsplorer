package learning

import (
	"sync"
)

// LearningOptimiser
type LearningOptimiser struct {
	Learner Learner
}

func (l *LearningOptimiser) Step(wg *sync.WaitGroup) {

}

func (l *LearningOptimiser) Run() {
	var wg sync.WaitGroup

	l.Step(&wg)
}
