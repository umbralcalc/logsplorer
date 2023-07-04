package learning

import (
	"sync"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// LearningObjective evaluates the objective function that needs to be optimised
// by running the stochadex simulator as a data iterator and computing the
// cumulative log-likelihood.
type LearningObjective struct {
	config          *LearningConfig
	stochadexConfig *simulator.StochadexConfig
	dataIterators   []*DataIterator
}

func (l *LearningObjective) Evaluate(
	newParams []*simulator.OtherParams,
) float64 {
	var wg sync.WaitGroup

	// set the incoming new params for each state partition
	for index := range l.stochadexConfig.Partitions {
		l.stochadexConfig.Partitions[index].Params.Other = newParams[index]
	}

	// instantiate a new batch of data iterators via the stochadex
	coordinator := simulator.NewPartitionCoordinator(l.stochadexConfig)

	// run the iterations over the data and terminate the for loop
	// when the end of data condition has been met
	for !coordinator.ReadyToTerminate() {
		coordinator.Step(&wg)
	}

	// store the objective value found in this step
	objective := 0.0
	for _, iterator := range l.dataIterators {
		objective += iterator.GetObjective()
	}

	// reset the stateful data iterators to be used again
	l.ResetIterators()

	return objective
}

func (l *LearningObjective) ResetIterators() {
	for i, objective := range l.config.Objectives {
		dataIterator := NewDataIterator(
			objective,
			l.config.Streaming[i].DataStreamer,
			l.config.BurnInSteps,
		)
		l.stochadexConfig.Partitions[i].Iteration = dataIterator
		l.dataIterators[i] = dataIterator
	}
}

// NewLearningObjective creates a new LearningObjective struct given a config
// and loaded settings.
func NewLearningObjective(
	config *LearningConfig,
	settings *simulator.LoadSettingsConfig,
) *LearningObjective {
	// handle some initial typing nonsense
	iterations := make([]simulator.Iteration, 0)
	dataIterators := make([]*DataIterator, 0)
	for i, objective := range config.Objectives {
		dataIterator := NewDataIterator(
			objective,
			config.Streaming[i].DataStreamer,
			config.BurnInSteps,
		)
		iterations = append(iterations, dataIterator)
		dataIterators = append(dataIterators, dataIterator)
	}
	implementations := &simulator.LoadImplementationsConfig{
		Iterations:           iterations,
		OutputCondition:      &simulator.NilOutputCondition{},
		OutputFunction:       &simulator.NilOutputFunction{},
		TerminationCondition: config.Streaming[0].TerminationCondition,
		TimestepFunction:     config.Streaming[0].TimestepFunction,
	}
	stochadexConfig := simulator.NewStochadexConfig(
		settings,
		implementations,
	)
	return &LearningObjective{
		config:          config,
		stochadexConfig: stochadexConfig,
		dataIterators:   dataIterators,
	}
}
