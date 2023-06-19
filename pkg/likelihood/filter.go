package likelihood

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

// LogLikelihood
type LogLikelihood interface {
	Evaluate(
		params *simulator.OtherParams,
		partitionIndex int,
		stateHistories []*simulator.StateHistory,
		timestepsHistory *simulator.CumulativeTimestepsHistory,
	) float64
}

// ConditionalProbability
type ConditionalProbability interface {
	SetParams(params *simulator.OtherParams)
	Compute(
		currentState []float64,
		pastState []float64,
		currentTime float64,
		pastTime float64,
	) float64
}

// ProbabilityFilterLogLikelihood
type ProbabilityFilterLogLikelihood struct {
	prob     ConditionalProbability
	dataLink DataLinkingLogLikelihood
}

func (p *ProbabilityFilterLogLikelihood) ComputeStatistics(
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) *Statistics {
	stateHistory := stateHistories[partitionIndex]
	currentTime := timestepsHistory.Values.AtVec(0)
	currentStateValue := stateHistory.Values.RawRowView(0)
	cumulativeWeightSum := 0.0
	mean := make([]float64, stateHistory.StateWidth)
	// i = 1 because we ignore the first (most recent) value in the history
	// as this is the one we want to compare to in the log-likelihood
	for i := 1; i < stateHistory.StateHistoryDepth; i++ {
		weight := p.prob.Compute(
			currentStateValue,
			stateHistory.Values.RawRowView(i),
			currentTime,
			timestepsHistory.Values.AtVec(i),
		)
		cumulativeWeightSum += weight
		if i == 1 {
			for j := 0; j < stateHistory.StateWidth; j++ {
				mean[j] = weight * stateHistory.Values.At(i, j)
			}
			continue
		}
		for j := 0; j < stateHistory.StateWidth; j++ {
			mean[j] += weight * stateHistory.Values.At(i, j)
		}
	}
	meanVec := mat.NewVecDense(stateHistory.StateWidth, mean)
	meanVec.ScaleVec(1.0/cumulativeWeightSum, meanVec)
	statistics := &Statistics{Mean: meanVec}
	// this next chunk to compute the covariance is a tweaked version from
	// https://github.com/gonum/gonum/blob/v0.13.0/stat/statmat.go
	// because we already have computed the mean so this'll be more efficient
	var stateHistoryValuesTrans *mat.Dense
	stateHistoryValuesTrans.CloneFrom(stateHistory.Values.T())
	for j := 0; j < stateHistory.StateWidth; j++ {
		floats.AddConst(-mean[j], stateHistoryValuesTrans.RawRowView(j))
	}
	covMat := mat.NewSymDense(stateHistory.StateWidth*stateHistory.StateWidth, nil)
	covMat.SymOuterK(
		1.0/cumulativeWeightSum,
		stateHistoryValuesTrans.Slice(
			0, stateHistory.StateWidth, 1, stateHistory.StateHistoryDepth,
		),
	)
	statistics.SetCovariance(covMat)
	return statistics
}

func (p *ProbabilityFilterLogLikelihood) Evaluate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) float64 {
	p.prob.SetParams(params)
	statistics := p.ComputeStatistics(
		partitionIndex,
		stateHistories,
		timestepsHistory,
	)
	logLikelihood := p.dataLink.Evaluate(
		statistics,
		stateHistories[partitionIndex].Values.RawRowView(0),
	)
	return logLikelihood
}
