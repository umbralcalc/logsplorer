package learning

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// LogLikelihood defines the interface that must be implemented in order to
// create an objective function which is applied to the data in an iterative
// fashion by the DataIterator class.
type LogLikelihood interface {
	Evaluate(
		params *simulator.OtherParams,
		partitionIndex int,
		stateHistories []*simulator.StateHistory,
		timestepsHistory *simulator.CumulativeTimestepsHistory,
	) float64
}

// DataIterator implements the stochadex simulator.Iteration interface to
// allow for iteration through streamed data while evaluating the objective
// function in the form of a cumulative log-likelihood as it goes. It also
// extends this interface to include the ability to output an objective
// function via the .GetObjective() method call.
type DataIterator struct {
	logLikelihood           LogLikelihood
	streamer                DataStreamer
	cumulativeLogLikelihood float64
	burnInSteps             int
	stepsTaken              int
}

func (d *DataIterator) Configure(
	partitionIndex int,
	settings *simulator.LoadSettingsConfig,
) {
	d.burnInSteps = int(
		settings.OtherParams[partitionIndex].IntParams["burn_in_steps"][0],
	)
	d.cumulativeLogLikelihood = 0.0
	d.stepsTaken = 0
}

func (d *DataIterator) Iterate(
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
func (d *DataIterator) GetObjective() float64 {
	return d.cumulativeLogLikelihood
}
