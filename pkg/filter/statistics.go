package filter

import (
	"math"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

// Statistics
type Statistics interface {
	GetMean() *mat.VecDense
	GetCovariance() *mat.SymDense
	ComputeAdditional(
		mean *mat.VecDense,
		weights []float64,
		stateHistory *simulator.StateHistory,
		timestepsHistory *simulator.CumulativeTimestepsHistory,
	)
}

// StandardCovarianceStatistics computes the standard covariance matrix.
type StandardCovarianceStatistics struct {
	Mean       *mat.VecDense
	Covariance *mat.SymDense
}

func (s *StandardCovarianceStatistics) ComputeAdditional(
	mean *mat.VecDense,
	weights []float64,
	stateHistory *simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) {
	// this next chunk to compute the covariance is a tweaked version from
	// https://github.com/gonum/gonum/blob/v0.13.0/stat/statmat.go
	// because we already have computed the mean so this'll be less work
	var stateHistoryValuesTrans mat.Dense
	stateHistoryValuesTrans.CloneFrom(stateHistory.Values.T())
	for j := 0; j < stateHistory.StateWidth; j++ {
		v := stateHistoryValuesTrans.RawRowView(j)
		floats.AddConst(-mean.AtVec(j), v)
	}
	covMat := mat.NewSymDense(stateHistory.StateWidth, nil)
	sqrtWeights := make([]float64, stateHistory.StateHistoryDepth)
	cumulativeWeightSum := 0.0
	for i, w := range weights {
		if w < 0 {
			panic("stat: negative covariance matrix weights")
		}
		sqrtWeights[i] = math.Sqrt(w)
		cumulativeWeightSum += w
	}
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
	s.Mean = mean
	s.Covariance = covMat
}

func (s *StandardCovarianceStatistics) GetMean() *mat.VecDense {
	return s.Mean
}

func (s *StandardCovarianceStatistics) GetCovariance() *mat.SymDense {
	return s.Covariance
}
