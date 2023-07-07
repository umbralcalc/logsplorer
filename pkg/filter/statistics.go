package filter

import (
	"math"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

// Statistics estimates the statistics of the empirical distribution in the
// probability filtering algorithm.
// TODO: Support banded matrices to enable spatial statistics.
type Statistics struct {
	Mean       *mat.VecDense
	Covariance mat.Symmetric
}

func (s *Statistics) Compute(
	prob ConditionalProbability,
	stateHistory *simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) {
	currentTime := timestepsHistory.Values.AtVec(0)
	currentStateValue := stateHistory.Values.RawRowView(0)
	cumulativeWeightSum := 0.0
	mean := make([]float64, stateHistory.StateWidth)
	weights := make([]float64, stateHistory.StateHistoryDepth)
	sqrtWeights := make([]float64, stateHistory.StateHistoryDepth)

	// i = 1 because we ignore the first (most recent) value in the history
	// as this is the one we want to compare to in the log-likelihood
	for i := 1; i < stateHistory.StateHistoryDepth; i++ {
		weights[i] = prob.Evaluate(
			currentStateValue,
			stateHistory.Values.RawRowView(i),
			currentTime,
			timestepsHistory.Values.AtVec(i),
		)
		if weights[i] < 0 {
			panic("stat: negative covariance matrix weights")
		}
		sqrtWeights[i] = math.Sqrt(weights[i])
		cumulativeWeightSum += weights[i]
		if i == 1 {
			for j := 0; j < stateHistory.StateWidth; j++ {
				mean[j] = weights[i] * stateHistory.Values.At(i, j)
			}
			continue
		}
		for j := 0; j < stateHistory.StateWidth; j++ {
			mean[j] += weights[i] * stateHistory.Values.At(i, j)
		}
	}
	meanVec := mat.NewVecDense(stateHistory.StateWidth, mean)
	meanVec.ScaleVec(1.0/cumulativeWeightSum, meanVec)

	// this next chunk to compute the covariance is a tweaked version from
	// https://github.com/gonum/gonum/blob/v0.13.0/stat/statmat.go
	// because we already have computed the mean so this'll be less work
	var stateHistoryValuesTrans mat.Dense
	stateHistoryValuesTrans.CloneFrom(stateHistory.Values.T())
	for j := 0; j < stateHistory.StateWidth; j++ {
		v := stateHistoryValuesTrans.RawRowView(j)
		floats.AddConst(-meanVec.AtVec(j), v)
	}
	covMat := mat.NewSymDense(stateHistory.StateWidth, nil)
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
	s.Mean = meanVec
	s.Covariance = covMat
}
