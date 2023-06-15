package learning

import "github.com/umbralcalc/stochadex/pkg/simulator"

// LearningCoordinator
type LearningCoordinator struct {
	Learners    []*Learner
	coordinator *simulator.PartitionCoordinator
}

func (l *LearningCoordinator) Step() {

}
