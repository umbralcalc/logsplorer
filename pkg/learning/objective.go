package learning

import (
	"sync"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// LearningObjective evaluates the objective function that needs to be optimised
// by running the stochadex simulator as a data iterator and computing the
// cumulative log-likelihood.
type LearningObjective struct {
	Iterations      []*IterationWithObjective
	OutputFunction  ObjectiveOutputFunction
	config          *LearningConfig
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
	coordinator := simulator.NewPartitionCoordinator(
		simulator.NewStochadexConfig(l.settings, l.implementations),
	)

	// run the iterations over the data and terminate the for loop
	// when the end of data condition has been met
	for !coordinator.ReadyToTerminate() {
		coordinator.Step(&wg)
	}

	// store the objective value found in this step and output its values for
	// each of the state partitions
	objective := 0.0
	for partitionIndex, iteration := range l.Iterations {
		partitionObjective := iteration.GetObjective()
		l.OutputFunction.Output(
			partitionIndex,
			partitionObjective,
			newParams[partitionIndex],
		)
		objective += partitionObjective
	}

	// reset the stateful data iterators to be used again
	l.ResetIterators()

	return objective
}

func (l *LearningObjective) ResetIterators() {
	for i, iteration := range l.Iterations {
		iteration.Configure(i, l.settings)
	}
}

// NewLearningObjective creates a new LearningObjective struct given a config
// and loaded settings.
func NewLearningObjective(config *LearningConfig) *LearningObjective {
	dataIterations := make([]*IterationWithObjective, 0)
	for i, objective := range config.Objectives {
		iteration := &IterationWithObjective{
			logLikelihood: objective,
			iteration:     config.Streaming.Iterations[i],
		}
		iteration.Configure(i, config.StreamingSettings)
		dataIterations = append(dataIterations, iteration)
		config.Streaming.Iterations[i] = iteration
	}
	return &LearningObjective{
		Iterations:      dataIterations,
		OutputFunction:  config.ObjectiveOutput,
		config:          config,
		settings:        config.StreamingSettings,
		implementations: config.Streaming,
	}
}
