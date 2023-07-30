package learning

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// LogLikelihood defines the interface that must be implemented in order to
// create an objective function which is applied to the data in an iterative
// fashion by the DataIteration class.
type LogLikelihood interface {
	Configure(partitionIndex int, settings *simulator.LoadSettingsConfig)
	Evaluate(
		params *simulator.OtherParams,
		partitionIndex int,
		stateHistories []*simulator.StateHistory,
		timestepsHistory *simulator.CumulativeTimestepsHistory,
	) float64
}

// IterationWithObjective allows for iteration through a given stochadex simulation
// while evaluating a specified object function in the form of a cumulative
// log-likelihood as it goes.
type IterationWithObjective struct {
	logLikelihood           LogLikelihood
	iteration               simulator.Iteration
	cumulativeLogLikelihood float64
	burnInSteps             int
	stepsTaken              int
}

func (i *IterationWithObjective) Configure(
	partitionIndex int,
	settings *simulator.LoadSettingsConfig,
) {
	i.logLikelihood.Configure(partitionIndex, settings)
	i.burnInSteps = int(
		settings.OtherParams[partitionIndex].IntParams["burn_in_steps"][0],
	)
	i.cumulativeLogLikelihood = 0.0
	i.stepsTaken = 0
}

func (i *IterationWithObjective) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	i.stepsTaken += 1
	if i.stepsTaken <= i.burnInSteps {
		return i.iteration.Iterate(params, partitionIndex, stateHistories, timestepsHistory)
	}
	i.cumulativeLogLikelihood += i.logLikelihood.Evaluate(
		params,
		partitionIndex,
		stateHistories,
		timestepsHistory,
	)
	return i.iteration.Iterate(params, partitionIndex, stateHistories, timestepsHistory)
}

// Get the cumulative log-likelihood that has been calculated by iterating
// through the data and applying the LogLikehood.Evaluate method.
func (i *IterationWithObjective) GetObjective() float64 {
	return i.cumulativeLogLikelihood
}
