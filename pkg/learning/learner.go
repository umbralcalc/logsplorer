package learning

import (
	"sync"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// Learner
type Learner struct {
	config             *LearnerConfig
	stochadexConfig    *simulator.StochadexConfig
	dataIterations     []DataIteration
	objectiveHistories [][]float64
}

func (l *Learner) Step(wg *sync.WaitGroup) {
	// instantiate a new batch of data iterators via the stochadex
	coordinator := simulator.NewPartitionCoordinator(l.stochadexConfig)

	// run the iterations over the data and terminate the for loop
	// when the end of data condition has been met
	for !coordinator.ReadyToTerminate() {
		coordinator.Step(wg)
	}

	// store the objective values found in this step
	for i, iteration := range l.dataIterations {
		l.objectiveHistories[i] = append(
			l.objectiveHistories[i],
			iteration.GetObjective(),
		)
	}
}

func (l *Learner) Run(wg *sync.WaitGroup) {
	l.Step(wg)
}

// NewLearner creates a new Learner struct.
func NewLearner(config *LearnerConfig) *Learner {
	settings := &simulator.LoadSettingsConfig{}
	// handle some typing nonsense
	iterations := make([]simulator.Iteration, 0)
	dataIterations := make([]DataIteration, 0)
	for i, objective := range config.Objectives {
		dataIterator := NewDataIterator(objective, config.Streaming[i].DataStreamer)
		iterations = append(iterations, dataIterator)
		dataIterations = append(dataIterations, dataIterator)
	}
	implementations := &simulator.LoadImplementationsConfig{
		Iterations:           iterations,
		TerminationCondition: config.Streaming[0].TerminationCondition,
		TimestepFunction:     config.Streaming[0].TimestepFunction,
	}
	stochadexConfig := simulator.NewStochadexConfig(
		settings,
		implementations,
	)
	return &Learner{
		config:             config,
		stochadexConfig:    stochadexConfig,
		dataIterations:     dataIterations,
		objectiveHistories: make([][]float64, 0),
	}
}
