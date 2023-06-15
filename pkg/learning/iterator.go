package learning

import (
	"github.com/umbralcalc/learnadex/pkg/likelihood"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// DataIterator
type DataIterator struct {
	CumulativeLogLikelihood float64
	logLikelihood           likelihood.LogLikelihood
	streamer                DataStreamer
}

func (d *DataIterator) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	d.CumulativeLogLikelihood += d.logLikelihood.Evaluate(
		params,
		partitionIndex,
		stateHistories,
		timestepsHistory,
	)
	return d.streamer.NextValue()
}

// NewDataIterator creates a new DataIterator.
func NewDataIterator(
	logLikelihood likelihood.LogLikelihood,
	streamer DataStreamer,
) *DataIterator {
	return &DataIterator{
		CumulativeLogLikelihood: 0.0,
		logLikelihood:           logLikelihood,
		streamer:                streamer,
	}
}
