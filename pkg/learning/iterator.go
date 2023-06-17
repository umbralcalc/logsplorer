package learning

import (
	"github.com/umbralcalc/learnadex/pkg/likelihood"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// DataIteration
type DataIteration interface {
	GetObjective() float64
	simulator.Iteration
}

// DataIterator
type DataIterator struct {
	cumulativeLogLikelihood float64
	logLikelihood           likelihood.LogLikelihood
	streamer                DataStreamer
}

func (d *DataIterator) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	d.cumulativeLogLikelihood += d.logLikelihood.Evaluate(
		params,
		partitionIndex,
		stateHistories,
		timestepsHistory,
	)
	return d.streamer.NextValue()
}

func (d *DataIterator) GetObjective() float64 {
	return d.cumulativeLogLikelihood
}

// NewDataIterator creates a new DataIterator.
func NewDataIterator(
	logLikelihood likelihood.LogLikelihood,
	streamer DataStreamer,
) *DataIterator {
	return &DataIterator{
		cumulativeLogLikelihood: 0.0,
		logLikelihood:           logLikelihood,
		streamer:                streamer,
	}
}
