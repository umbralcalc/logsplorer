package learning

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// ObjectiveEvaluator evaluates the objective function that needs to
// be optimised by running the stochadex simulator as a data iterator and
// computing the cumulative log-likelihood.
type ObjectiveEvaluator struct {
	Iterations      []*IterationWithObjective
	OutputFunction  ObjectiveOutputFunction
	config          *LearningConfig
	settings        *simulator.LoadSettingsConfig
	implementations *simulator.LoadImplementationsConfig
}

func (o *ObjectiveEvaluator) Evaluate(
	newParams []*simulator.OtherParams,
) float64 {
	// set the incoming new params for each state partition
	for index := range o.settings.OtherParams {
		o.settings.OtherParams[index] = newParams[index]
	}

	// instantiate a new batch of data iterators via the stochadex
	coordinator := simulator.NewPartitionCoordinator(
		simulator.NewStochadexConfig(o.settings, o.implementations),
	)

	// run the iterations over the data
	coordinator.Run()

	// store the objective value found in this step and output its values for
	// each of the state partitions
	objective := 0.0
	for partitionIndex, iteration := range o.Iterations {
		partitionObjective := iteration.GetObjective()
		o.OutputFunction.Output(
			partitionIndex,
			partitionObjective,
			newParams[partitionIndex],
		)
		objective += partitionObjective
	}

	// reset the stateful data iterators to be used again
	o.ResetIterations(o.settings)

	return objective
}

func (o *ObjectiveEvaluator) Copy() *ObjectiveEvaluator {
	evaluatorCopy := *o
	evaluatorCopy.ResetIterations(o.settings)
	return &evaluatorCopy
}

func (o *ObjectiveEvaluator) ResetIterations(
	settings *simulator.LoadSettingsConfig,
) {
	for i, iteration := range o.Iterations {
		iteration.Configure(i, settings)
	}
}

// NewObjectiveEvaluator creates a new ObjectiveEvaluator struct
// given a config and loaded settings.
func NewObjectiveEvaluator(config *LearningConfig) *ObjectiveEvaluator {
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
	return &ObjectiveEvaluator{
		Iterations:      dataIterations,
		OutputFunction:  config.ObjectiveOutput,
		config:          config,
		settings:        config.StreamingSettings,
		implementations: config.Streaming,
	}
}
