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
	dataIterators   []*DataIterator
	settings        *simulator.LoadSettingsConfig
	implementations *simulator.LoadImplementationsConfig
}

func (l *LearningObjective) Evaluate(
	newParams []*simulator.OtherParams,
) float64 {
	var wg sync.WaitGroup

	// set the incoming new params for each state partition
	for index := range l.settings.OtherParams {
		l.settings.OtherParams[index] = newParams[index]
	}

	// instantiate a new batch of data iterators via the stochadex
	stochadexConfig := simulator.NewStochadexConfig(
		l.settings,
		l.implementations,
	)
	coordinator := simulator.NewPartitionCoordinator(stochadexConfig)

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
		dataIterator := &DataIterator{
			logLikelihood: objective,
			streamer:      l.config.Streaming[i].DataStreamer,
		}
		dataIterator.Configure(i, l.settings)
		l.implementations.Iterations[i] = dataIterator
		l.dataIterators[i] = dataIterator
	}
}

// NewLearningObjective creates a new LearningObjective struct given a config
// and loaded settings.
func NewLearningObjective(
	config *LearningConfig,
	settings *simulator.LoadSettingsConfig,
) *LearningObjective {
	iterations := make([]simulator.Iteration, 0)
	dataIterators := make([]*DataIterator, 0)
	for i, objective := range config.Objectives {
		dataIterator := &DataIterator{
			logLikelihood: objective,
			streamer:      config.Streaming[i].DataStreamer,
		}
		dataIterator.Configure(i, settings)
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
	return &LearningObjective{
		config:          config,
		dataIterators:   dataIterators,
		settings:        settings,
		implementations: implementations,
	}
}
