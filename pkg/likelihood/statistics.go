package likelihood

import (
	"math"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

// Statistics
type Statistics interface {
	SetMean(meanVector *mat.VecDense)
	GetMean() *mat.VecDense
	GetCovariance() *mat.SymDense
	ComputeOthers(
		weights []float64,
		stateHistory *simulator.StateHistory,
		timestepsHistory *simulator.CumulativeTimestepsHistory,
	)
}

// CovarianceStatistics is a struct which holds the basice mean and covariance
// statistics and makes some useful transformations on this data.
type CovarianceStatistics struct {
	Mean       *mat.VecDense
	Covariance *mat.SymDense
}

func (c *CovarianceStatistics) ComputeOthers(
	weights []float64,
	stateHistory *simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) {
	// this next chunk to compute the covariance is a tweaked version from
	// https://github.com/gonum/gonum/blob/v0.13.0/stat/statmat.go
	// because we already have computed the mean so this'll be more efficient
	var stateHistoryValuesTrans *mat.Dense
	stateHistoryValuesTrans.CloneFrom(stateHistory.Values.T())
	for j := 0; j < stateHistory.StateWidth; j++ {
		v := stateHistoryValuesTrans.RawRowView(j)
		floats.AddConst(-c.Mean.AtVec(j), v)
	}
	covMat := mat.NewSymDense(stateHistory.StateWidth*stateHistory.StateWidth, nil)
	// Multiply by the sqrt of the weights, so that multiplication is symmetric.
	sqrtWeights := make([]float64, stateHistory.StateHistoryDepth)
	cumulativeWeightSum := 0.0
	for i, w := range weights {
		if w < 0 {
			panic("stat: negative covariance matrix weights")
		}
		sqrtWeights[i] = math.Sqrt(w)
		cumulativeWeightSum += w
	}
	// Weight the rows.
	for j := 0; j < stateHistory.StateWidth; j++ {
		v := stateHistoryValuesTrans.RawRowView(j)
		floats.Mul(v, sqrtWeights)
	}
	covMat.SymOuterK(
		1.0/cumulativeWeightSum,
		stateHistoryValuesTrans.Slice(
			0, stateHistory.StateWidth, 1, stateHistory.StateHistoryDepth,
		),
	)
	c.Covariance = covMat
}

func (c *CovarianceStatistics) SetMean(meanVector *mat.VecDense) {
	c.Mean = meanVector
}

func (c *CovarianceStatistics) GetMean() *mat.VecDense {
	return c.Mean
}

func (c *CovarianceStatistics) GetCovariance() *mat.SymDense {
	return c.Covariance
}
