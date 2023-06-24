package learning

import (
	"log"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/optimize"
)

// ParamsTranslator
type ParamsTranslator interface {
	ToOptimiser(paramsToTranslate []*simulator.OtherParams) []float64
	FromOptimiser(
		fromOptimiser []float64,
		paramsToUpdate []*simulator.OtherParams,
	) []*simulator.OtherParams
}

// OptimisationAlgorithm
type OptimisationAlgorithm interface {
	Run(
		learningObj LearningObjective,
		paramsTranslator ParamsTranslator,
		initialParams []*simulator.OtherParams,
	) []*simulator.OtherParams
}

// GonumOptimisationAlgorithm
type GonumOptimisationAlgorithm struct {
	Method optimize.Method
}

func (g *GonumOptimisationAlgorithm) Run(
	learningObj LearningObjective,
	paramsTranslator ParamsTranslator,
	initialParams []*simulator.OtherParams,
) []*simulator.OtherParams {
	problem := optimize.Problem{
		Func: func(x []float64) float64 {
			return learningObj.Evaluate(
				paramsTranslator.FromOptimiser(x, initialParams),
			)
		},
	}
	result, err := optimize.Minimize(
		problem,
		paramsTranslator.ToOptimiser(initialParams),
		nil,
		g.Method,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = result.Status.Err(); err != nil {
		log.Fatal(err)
	}
	return initialParams
}
