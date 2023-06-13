package learning

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/mat"
)

// ConditionalProbability
type ConditionalProbability interface {
	Compute(
		currentState []float64,
		pastState []float64,
		currentTime float64,
		pastTime float64,
	) float64
}

// ProbabilityFilterLogLikelihood
type ProbabilityFilterLogLikelihood struct {
	conditionalProbability ConditionalProbability
}

func (p *ProbabilityFilterLogLikelihood) ComputeStatistics(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.TimestepsHistory,
) *Statistics {
	cumulativeWeightsSum := 0.0
	mean := make([]float64, 0)
	flatCov := make([]float64, 0)
	for i, dataVector := range dataHistory {
		// ignore first (most recent) value in the history as this
		// is the one we want to compare to in the log-likelihood
		if i == 0 {
			continue
		}
		weight := ConditionalProbability(
			dataHistory[0],
			dataVector,
			hyperparams,
		)
		cumulativeWeightsSum += weight
		for j := range meanVector {
			meanVector[j] += weight * dataVector[j]
		}
	}
	// normalise the weights derived from the conditional probability
	for j := range meanVector {
		meanVector[j] /= cumulativeWeightsSum
	}
	stats := &Statistics{Mean: mat.NewVecDense(len(mean), mean)}
	stats.SetCovariance(mat.NewSymDense(len(mean), flatCov))
	return meanVector
}

func (p *ProbabilityFilterLogLikelihood) Evaluate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.TimestepsHistory,
) float64 {
	// get the most recent point in the data history
	dataVector := stateHistories[partitionIndex].Values.RowView(0)
	stats := p.ComputeStatistics(dataHistory, hyperparams)
	logLikelihood := DataLinkingLogProbability(
		dataVector,
		mean,
		hyperparams,
	)
	return logLikelihood
}
