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

// IterationWithObjective extends the simulator.Iteration interface to
// support generating a cumulative objective function output at the end of
// a given iteration run.
type IterationWithObjective interface {
	simulator.Iteration
	GetObjective() float64
}

// StochadexIterationWithObjective implements the IterationWithObjective
// interface and allows for iteration through a given stochadex simulation while
// evaluating a specified object function in the form of a cumulative log-likelihood
// as it goes.
type StochadexIterationWithObjective struct {
	logLikelihood           LogLikelihood
	iteration               simulator.Iteration
	cumulativeLogLikelihood float64
	burnInSteps             int
	stepsTaken              int
}

func (s *StochadexIterationWithObjective) Configure(
	partitionIndex int,
	settings *simulator.LoadSettingsConfig,
) {
	s.logLikelihood.Configure(partitionIndex, settings)
	s.burnInSteps = int(
		settings.OtherParams[partitionIndex].IntParams["burn_in_steps"][0],
	)
	s.cumulativeLogLikelihood = 0.0
	s.stepsTaken = 0
}

func (s *StochadexIterationWithObjective) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	s.stepsTaken += 1
	if s.stepsTaken <= s.burnInSteps {
		return s.iteration.Iterate(params, partitionIndex, stateHistories, timestepsHistory)
	}
	s.cumulativeLogLikelihood += s.logLikelihood.Evaluate(
		params,
		partitionIndex,
		stateHistories,
		timestepsHistory,
	)
	return s.iteration.Iterate(params, partitionIndex, stateHistories, timestepsHistory)
}

// Get the cumulative log-likelihood that has been calculated by iterating
// through the data and applying the LogLikehood.Evaluate method.
func (s *StochadexIterationWithObjective) GetObjective() float64 {
	return s.cumulativeLogLikelihood
}

// DataIterationWithObjective implements the IterationWithObjective interface
// and allows for iteration through streamed data while evaluating the objective
// function in the form of a cumulative log-likelihood as it goes.
type DataIterationWithObjective struct {
	logLikelihood           LogLikelihood
	streamer                DataStreamer
	cumulativeLogLikelihood float64
	burnInSteps             int
	stepsTaken              int
}

func (d *DataIterationWithObjective) Configure(
	partitionIndex int,
	settings *simulator.LoadSettingsConfig,
) {
	d.logLikelihood.Configure(partitionIndex, settings)
	d.burnInSteps = int(
		settings.OtherParams[partitionIndex].IntParams["burn_in_steps"][0],
	)
	d.cumulativeLogLikelihood = 0.0
	d.stepsTaken = 0
}

func (d *DataIterationWithObjective) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	d.stepsTaken += 1
	if d.stepsTaken <= d.burnInSteps {
		return d.streamer.NextValue(timestepsHistory)
	}
	d.cumulativeLogLikelihood += d.logLikelihood.Evaluate(
		params,
		partitionIndex,
		stateHistories,
		timestepsHistory,
	)
	return d.streamer.NextValue(timestepsHistory)
}

// Get the cumulative log-likelihood that has been calculated by iterating
// through the data and applying the LogLikehood.Evaluate method.
func (d *DataIterationWithObjective) GetObjective() float64 {
	return d.cumulativeLogLikelihood
}
