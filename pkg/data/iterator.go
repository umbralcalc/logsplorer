package data

import (
	"github.com/umbralcalc/learnadex/pkg/learning"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// DataIterator
type DataIterator struct {
	CumulativeLogLike float64
	logLikelihood     learning.LogLikelihood
	streamer          DataStreamer
}

func (d *DataIterator) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	d.CumulativeLogLike += d.logLikelihood.Evaluate(
		params,
		partitionIndex,
		stateHistories,
		timestepsHistory,
	)
	return d.streamer.NextValue()
}

func NewDataIterator(
	logLikelihood learning.LogLikelihood,
	streamer DataStreamer,
) *DataIterator {
	return &DataIterator{
		CumulativeLogLike: 0.0,
		logLikelihood:     logLikelihood,
		streamer:          streamer,
	}
}
