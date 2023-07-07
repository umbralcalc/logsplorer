package learning

import (
	"log"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/optimize"
)

// OptimisationAlgorithm defines the interface that must be implemented
// in order to specify an algorithm that can optimise the learning
// objective specified.
type OptimisationAlgorithm interface {
	Run(
		learningObj *LearningObjective,
		initialParams []*simulator.OtherParams,
	) []*simulator.OtherParams
}

// ParamsTranslator defines the interface that must be implemented in order
// to translate the params which are passed to the probability filter
// algorithm into a form which the optimiser OptimimisationAlgorithm can
// always use. Writing this translator for a specific problem domain also
// enables only translating a subset of all of the parameters if only
// optimising a subset is desired.
type ParamsTranslator interface {
	ToOptimiser(paramsToTranslate []*simulator.OtherParams) []float64
	FromOptimiser(
		fromOptimiser []float64,
		paramsToUpdate []*simulator.OtherParams,
	) []*simulator.OtherParams
}

// NewParamsCopy is a convenience function which copies the input
// []*simulator.OtherParams to help ensure thread safety.
func NewParamsCopy(params []*simulator.OtherParams) []*simulator.OtherParams {
	paramsCopy := make([]*simulator.OtherParams, 0)
	for i := range params {
		p := *params[i]
		paramsCopy = append(paramsCopy, &p)
	}
	return paramsCopy
}

// GonumOptimisationAlgorithm allows any of the gonum optimisers to be
// used in the learnadex.
type GonumOptimisationAlgorithm struct {
	Method     optimize.Method
	Translator ParamsTranslator
}

func (g *GonumOptimisationAlgorithm) Run(
	learningObj *LearningObjective,
	initialParams []*simulator.OtherParams,
) []*simulator.OtherParams {
	problem := optimize.Problem{
		Func: func(x []float64) float64 {
			// this copying ensures thread safety (as required by
			// the gonum optimize package)
			learningObjCopy := *learningObj
			learningObjCopy.ResetIterators()
			paramsCopy := NewParamsCopy(initialParams)
			return learningObjCopy.Evaluate(
				g.Translator.FromOptimiser(x, paramsCopy),
			)
		},
	}
	result, err := optimize.Minimize(
		problem,
		g.Translator.ToOptimiser(initialParams),
		nil,
		g.Method,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = result.Status.Err(); err != nil {
		log.Fatal(err)
	}
	return g.Translator.FromOptimiser(result.X, initialParams)
}
