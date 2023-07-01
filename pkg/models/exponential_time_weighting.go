package models

import (
	"math"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// ExponentialTimeWeightingConditionalProbability can be used in the probability
// filter to learn a simple exponentially-weighted historical average value.
type ExponentialTimeWeightingConditionalProbability struct {
	timescale float64
}

func (e *ExponentialTimeWeightingConditionalProbability) SetParams(
	params *simulator.OtherParams,
) {
	e.timescale = params.FloatParams["exponential_weighting_timescale"][0]
}

func (e *ExponentialTimeWeightingConditionalProbability) Evaluate(
	currentState []float64,
	pastState []float64,
	currentTime float64,
	pastTime float64,
) float64 {
	return math.Exp((pastTime - currentTime) / e.timescale)
}
