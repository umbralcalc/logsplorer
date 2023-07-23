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
		paramsToOptimise []*simulator.OtherParamsMask,
	) []*simulator.OtherParams
}

// GetParamsForOptimiser is a convenience function which returns the params
// from the stochadex where the mask has been applied to them in a flattened
// single slice format ready to input into an optimiser.
func GetParamsForOptimiser(
	params []*simulator.OtherParams,
	paramsToOptimise []*simulator.OtherParamsMask,
) []float64 {
	paramsForOptimiser := make([]float64, 0)
	for index, partitionParams := range params {
		for name, paramSlice := range partitionParams.FloatParams {
			_, ok := paramsToOptimise[index].FloatParams[name]
			if !ok {
				continue
			}
			for i, param := range paramSlice {
				if paramsToOptimise[index].FloatParams[name][i] {
					paramsForOptimiser = append(paramsForOptimiser, param)
				}
			}
		}
		for name, paramSlice := range partitionParams.IntParams {
			_, ok := paramsToOptimise[index].IntParams[name]
			if !ok {
				continue
			}
			for i, param := range paramSlice {
				if paramsToOptimise[index].IntParams[name][i] {
					paramsForOptimiser = append(paramsForOptimiser, float64(param))
				}
			}
		}
	}
	return paramsForOptimiser
}

// UpdateParamsFromOptimiser is a convenience function which updates the input params
// from the stochadex which have been retrieved from the flattened slice format that
// is typically used in an optimiser package.
func UpdateParamsFromOptimiser(
	fromOptimiser []float64,
	params []*simulator.OtherParams,
	paramsToOptimise []*simulator.OtherParamsMask,
) []*simulator.OtherParams {
	indexInOptimiser := 0
	for index, partitionParams := range params {
		for name, paramSlice := range partitionParams.FloatParams {
			_, ok := paramsToOptimise[index].FloatParams[name]
			if !ok {
				continue
			}
			for i := range paramSlice {
				if paramsToOptimise[index].FloatParams[name][i] {
					params[index].FloatParams[name][i] = fromOptimiser[i]
					indexInOptimiser += 1
				}
			}
		}
		for name, paramSlice := range partitionParams.IntParams {
			_, ok := paramsToOptimise[index].IntParams[name]
			if !ok {
				continue
			}
			for i := range paramSlice {
				if paramsToOptimise[index].IntParams[name][i] {
					params[index].IntParams[name][i] = int64(fromOptimiser[i])
					indexInOptimiser += 1
				}
			}
		}
	}
	return params
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
	Method optimize.Method
}

func (g *GonumOptimisationAlgorithm) Run(
	learningObj *LearningObjective,
	initialParams []*simulator.OtherParams,
	paramsToOptimise []*simulator.OtherParamsMask,
) []*simulator.OtherParams {
	problem := optimize.Problem{
		Func: func(x []float64) float64 {
			// this copying ensures thread safety (as required by
			// the gonum optimize package)
			learningObjCopy := *learningObj
			learningObjCopy.ResetIterators()
			paramsCopy := NewParamsCopy(initialParams)
			return learningObjCopy.Evaluate(
				UpdateParamsFromOptimiser(x, paramsCopy, paramsToOptimise),
			)
		},
	}
	result, err := optimize.Minimize(
		problem,
		GetParamsForOptimiser(initialParams, paramsToOptimise),
		nil,
		g.Method,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = result.Status.Err(); err != nil {
		log.Fatal(err)
	}
	return UpdateParamsFromOptimiser(result.X, initialParams, paramsToOptimise)
}
