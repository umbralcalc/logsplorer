package learning

import "github.com/umbralcalc/stochadex/pkg/simulator"

// LikelihoodGenerator
type LikelihoodGenerator interface {
	Generate(
		params *simulator.OtherParams,
		stateHistory *simulator.StateHistory,
	) float64
}
