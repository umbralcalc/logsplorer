package models

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/mat"
)

// GaussianProcessCovarianceKernel
type GaussianProcessCovarianceKernel interface {
	SetParams(params *simulator.OtherParams)
	GetCovariance(
		currentState []float64,
		pastState []float64,
		currentTime float64,
		pastTime float64,
	) *mat.SymDense
}

// GaussianProcessConditionalProbability
type GaussianProcessConditionalProbability struct {
	Kernel       GaussianProcessCovarianceKernel
	timesByIndex map[int]float64
	meansInTime  map[float64][]float64
	stateWidth   int
}

func (g *GaussianProcessConditionalProbability) SetParams(
	params *simulator.OtherParams,
) {
	timeIndex := 0
	stateIndex := 0
	for _, mean := range params.FloatParams["flattened_means_in_time"] {
		g.meansInTime[g.timesByIndex[timeIndex]][stateIndex] = mean
		stateIndex += 1
		if stateIndex == g.stateWidth {
			timeIndex += 1
			stateIndex = 0
		}
	}
	g.Kernel.SetParams(params)
}

func (g *GaussianProcessConditionalProbability) Evaluate(
	currentState []float64,
	pastState []float64,
	currentTime float64,
	pastTime float64,
) float64 {
	// var currentDiff []float64
	// var pastDiff []float64
	// currentStateDiffVector := mat.NewVecDense(
	// 	g.stateWidth,
	// 	floats.SubTo(currentDiff, currentState, g.meansInTime[currentTime]),
	// )
	// pastStateDiffVector := mat.NewVecDense(
	// 	g.stateWidth,
	// 	floats.SubTo(pastDiff, pastState, g.meansInTime[pastTime]),
	// )
	// crossCov := g.Kernel.GetCovariance(currentState, pastState, currentTime, pastTime)

	return 0.0
}
