package learning

import (
	"sync"
)

// LearningOptimiser
type LearningOptimiser struct {
	Learners             []Learner
	objectiveHistories   [][]float64
	newWorkChannels      [](chan *LearnerInputMessage)
	workReturnedChannels [](chan *LearnerOutputMessage)
}

func (l *LearningOptimiser) StepBatch() {
	var wg sync.WaitGroup

	for index := range l.Learners {
		wg.Add(1)
		i := index
		go func() {
			defer wg.Done()
			l.Learners[i].ReceiveAndSendObjectives(
				l.newWorkChannels[i],
				l.workReturnedChannels[i],
			)
		}()
	}
}
