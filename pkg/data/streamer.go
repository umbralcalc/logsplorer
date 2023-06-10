package data

import "github.com/umbralcalc/stochadex/pkg/simulator"

// DataStreamer
type DataStreamer struct {
}

func (d *DataStreamer) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.TimestepsHistory,
) *simulator.State {
	return &simulator.State{}
}
