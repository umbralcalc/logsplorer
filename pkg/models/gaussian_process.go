package models

import (
	"math"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

const logTwoPi = 1.83788

// GaussianProcessCovarianceKernel is an interface that must be implemented
// in order to create a covariance kernel that can be used in the
// GaussianProcessConditionalProbability.
type GaussianProcessCovarianceKernel interface {
	Configure(partitionIndex int, settings *simulator.LoadSettingsConfig)
	SetParams(params *simulator.OtherParams)
	GetCovariance(
		currentState []float64,
		pastState []float64,
		currentTime float64,
		pastTime float64,
	) *mat.SymDense
}

// GaussianProcessConditionalProbability can be used in the probability
// filter to learn a Gaussian process kernel with a vector of means.
type GaussianProcessConditionalProbability struct {
	Kernel       GaussianProcessCovarianceKernel
	meansInTime  map[float64][]float64
	initialMeans []float64
	stateWidth   int
}

func (g *GaussianProcessConditionalProbability) Configure(
	partitionIndex int,
	settings *simulator.LoadSettingsConfig,
) {
	g.Kernel.Configure(partitionIndex, settings)
	g.meansInTime = make(map[float64][]float64)
	g.initialMeans = settings.OtherParams[partitionIndex].FloatParams["initial_means"]
	g.stateWidth = settings.StateWidths[partitionIndex]
	g.SetParams(settings.OtherParams[partitionIndex])
}

func (g *GaussianProcessConditionalProbability) SetParams(
	params *simulator.OtherParams,
) {
	g.Kernel.SetParams(params)
}

func (g *GaussianProcessConditionalProbability) Evaluate(
	currentState []float64,
	pastState []float64,
	currentTime float64,
	pastTime float64,
) float64 {
	_, ok := g.meansInTime[pastTime]
	if !ok {
		g.meansInTime[pastTime] = g.initialMeans
	}
	_, ok = g.meansInTime[currentTime]
	if !ok {
		g.meansInTime[currentTime] = g.initialMeans
	}
	currentDiff := make([]float64, g.stateWidth)
	pastDiff := make([]float64, g.stateWidth)
	currentStateDiffVector := mat.NewVecDense(
		g.stateWidth,
		floats.SubTo(currentDiff, currentState, g.meansInTime[currentTime]),
	)
	pastStateDiffVector := mat.NewVecDense(
		g.stateWidth,
		floats.SubTo(pastDiff, pastState, g.meansInTime[pastTime]),
	)
	var choleskyDecomp mat.Cholesky
	ok = choleskyDecomp.Factorize(
		g.Kernel.GetCovariance(currentState, pastState, currentTime, pastTime),
	)
	if !ok {
		return math.NaN()
	}
	var vectorInvMat mat.VecDense
	err := choleskyDecomp.SolveVecTo(&vectorInvMat, currentStateDiffVector)
	if err != nil {
		return math.NaN()
	}
	logResult := -0.5 * mat.Dot(&vectorInvMat, pastStateDiffVector)
	logResult -= 0.5 * float64(g.stateWidth) * logTwoPi
	logResult -= 0.5 * choleskyDecomp.LogDet()
	return math.Exp(logResult)
}
