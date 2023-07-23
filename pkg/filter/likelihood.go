package filter

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// ConditionalProbability is the interface that must be implemented in order
// to provide a conditionaly probability for the filtering algorithm.
type ConditionalProbability interface {
	Configure(partitionIndex int, settings *simulator.LoadSettingsConfig)
	SetParams(params *simulator.OtherParams)
	Evaluate(
		currentState []float64,
		pastState []float64,
		currentTime float64,
		pastTime float64,
	) float64
}

// UniformConditionalProbability implies a flat rolling window into the past to
// compute empirical statistics with.
type UniformConditionalProbability struct{}

func SetParams(params *simulator.OtherParams) {}

func Evaluate(
	currentState []float64,
	pastState []float64,
	currentTime float64,
	pastTime float64,
) float64 {
	return 1.0
}

// ProbabilityFilterLogLikelihood composes a provided data linking log-likelihood
// together with a provided conditional probability in order to implement the
// empirical probability filter algorithm as a LogLikelihood interface type.
type ProbabilityFilterLogLikelihood struct {
	Prob       ConditionalProbability
	DataLink   DataLinkingLogLikelihood
	Statistics *Statistics
}

func (p *ProbabilityFilterLogLikelihood) Configure(
	partitionIndex int,
	settings *simulator.LoadSettingsConfig,
) {
	p.Prob.Configure(partitionIndex, settings)
	p.DataLink.Configure(partitionIndex, settings)
}

func (p *ProbabilityFilterLogLikelihood) Evaluate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) float64 {
	p.Prob.SetParams(params)
	p.Statistics.Compute(
		p.Prob,
		stateHistories[partitionIndex],
		timestepsHistory,
	)
	logLikelihood := p.DataLink.Evaluate(
		p.Statistics,
		stateHistories[partitionIndex].Values.RawRowView(0),
	)
	return logLikelihood
}
