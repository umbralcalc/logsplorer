package learning

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// LogLikelihood
type LogLikelihood interface {
	Evaluate(
		params *simulator.OtherParams,
		partitionIndex int,
		stateHistories []*simulator.StateHistory,
		timestepsHistory *simulator.CumulativeTimestepsHistory,
	) float64
}

// DataIterator
type DataIterator struct {
	cumulativeLogLikelihood float64
	logLikelihood           LogLikelihood
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
	return d.streamer.NextValue(timestepsHistory)
}

func (d *DataIterator) GetObjective() float64 {
	return d.cumulativeLogLikelihood
}

// NewDataIterator creates a new DataIterator.
func NewDataIterator(
	logLikelihood LogLikelihood,
	streamer DataStreamer,
) *DataIterator {
	return &DataIterator{
		cumulativeLogLikelihood: 0.0,
		logLikelihood:           logLikelihood,
		streamer:                streamer,
	}
}
