package learning

import (
	"sync"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// Learner
type Learner struct {
	config             *LearnerConfig
	stochadexConfig    *simulator.StochadexConfig
	objectiveHistories [][]float64
}

func (l *Learner) Step(wg *sync.WaitGroup) {
	// instantiate a new batch of data iterators via the stochadex
	coordinator := simulator.NewPartitionCoordinator(l.stochadexConfig)

	// terminate the for loop if the condition has been met
	for !coordinator.ReadyToTerminate() {
		coordinator.Step(wg)
	}
	for i, iteration := range l.config.Iterations {
		l.objectiveHistories[i] = append(
			l.objectiveHistories[i],
			iteration.GetObjective(),
		)
	}
}

func (l *Learner) Run() {
	var wg sync.WaitGroup

	l.Step(&wg)
}

func NewLearner(config *LearnerConfig) *Learner {
	settings := &simulator.LoadSettingsConfig{}
	// handle some typing nonsense
	iterations := make([]simulator.Iteration, 0)
	for _, dataIteration := range config.Iterations {
		iterations = append(iterations, dataIteration)
	}
	implementations := &simulator.LoadImplementationsConfig{
		Iterations: iterations,
	}
	stochadexConfig := simulator.NewStochadexConfig(
		settings,
		implementations,
	)
	return &Learner{config: config, stochadexConfig: stochadexConfig}
}
