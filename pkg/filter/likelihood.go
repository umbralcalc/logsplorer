package filter

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/mat"
)

// ConditionalProbability
type ConditionalProbability interface {
	SetParams(params *simulator.OtherParams)
	Evaluate(
		currentState []float64,
		pastState []float64,
		currentTime float64,
		pastTime float64,
	) float64
}

// ProbabilityFilterLogLikelihood
type ProbabilityFilterLogLikelihood struct {
	prob       ConditionalProbability
	dataLink   DataLinkingLogLikelihood
	statistics Statistics
}

func (p *ProbabilityFilterLogLikelihood) ComputeStatistics(
	stateHistory *simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) {
	currentTime := timestepsHistory.Values.AtVec(0)
	currentStateValue := stateHistory.Values.RawRowView(0)
	cumulativeWeightSum := 0.0
	mean := make([]float64, stateHistory.StateWidth)
	weights := make([]float64, stateHistory.StateHistoryDepth)
	// i = 1 because we ignore the first (most recent) value in the history
	// as this is the one we want to compare to in the log-likelihood
	for i := 1; i < stateHistory.StateHistoryDepth; i++ {
		weights[i] = p.prob.Evaluate(
			currentStateValue,
			stateHistory.Values.RawRowView(i),
			currentTime,
			timestepsHistory.Values.AtVec(i),
		)
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
	p.statistics.ComputeAdditional(
		meanVec,
		weights,
		stateHistory,
		timestepsHistory,
	)
}

func (p *ProbabilityFilterLogLikelihood) Evaluate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) float64 {
	p.prob.SetParams(params)
	p.ComputeStatistics(stateHistories[partitionIndex], timestepsHistory)
	logLikelihood := p.dataLink.Evaluate(
		p.statistics,
		stateHistories[partitionIndex].Values.RawRowView(0),
	)
	return logLikelihood
}
