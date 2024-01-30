package learning

import (
	"log"

	"github.com/umbralcalc/learnadex/pkg/params"
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/optimize"
)

// OptimisationAlgorithm defines the interface that must be implemented
// in order to specify an algorithm that can optimise the learning
// objective specified.
type OptimisationAlgorithm interface {
	Run(
		evaluator *ObjectiveEvaluator,
		previousParams []*simulator.OtherParams,
	) []*simulator.OtherParams
}

// GonumOptimisationAlgorithm allows any of the gonum optimisers to be
// used in the learnadex.
type GonumOptimisationAlgorithm struct {
	Method   optimize.Method
	Settings *optimize.Settings
	mappings *params.OptimiserParamsMappings
}

func (g *GonumOptimisationAlgorithm) Run(
	evaluator *ObjectiveEvaluator,
	previousParams []*simulator.OtherParams,
) []*simulator.OtherParams {
	if g.mappings == nil {
		g.mappings = params.NewOptimiserParamsMappings(previousParams)
	}
	problem := optimize.Problem{
		Func: func(x []float64) float64 {
			// this copying ensures thread safety (as required by
			// the gonum optimize package)
			evaluatorCopy := evaluator.Copy()
			paramsCopy := params.NewParamsCopy(previousParams)
			return -evaluatorCopy.Evaluate(
				g.mappings.UpdateParamsFromOptimiser(x, paramsCopy),
			)
		},
	}
	result, err := optimize.Minimize(
		problem,
		g.mappings.GetParamsForOptimiser(previousParams),
		g.Settings,
		g.Method,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = result.Status.Err(); err != nil {
		log.Fatal(err)
	}
	return g.mappings.UpdateParamsFromOptimiser(result.X, previousParams)
}
