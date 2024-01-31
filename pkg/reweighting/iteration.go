package reweighting

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/mat"
)

// ProbabilisticReweightingIteration iterates the statistics computed
// using probabilistic reweighting forward in time.
type ProbabilisticReweightingIteration struct {
	Prob ConditionalProbability
}

func (p *ProbabilisticReweightingIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	p.Prob.Configure(partitionIndex, settings)
}

func (p *ProbabilisticReweightingIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	p.Prob.SetParams(params)
	stateHistory := stateHistories[partitionIndex]
	currentTime := timestepsHistory.Values.AtVec(0)
	currentStateValue := stateHistory.Values.RawRowView(0)
	cumulativeWeightSum := 0.0
	statisticsVec := mat.NewVecDense(stateHistory.StateWidth, nil)

	// i = 1 because we ignore the first (most recent) value in the history
	// as this is the one we want to compare to in the log-likelihood
	for i := 1; i < stateHistory.StateHistoryDepth; i++ {
		weight := p.Prob.Evaluate(
			currentStateValue,
			stateHistory.Values.RawRowView(i),
			currentTime,
			timestepsHistory.Values.AtVec(i),
		)
		if weight < 0 {
			panic("stat: negative covariance matrix weights")
		}
		cumulativeWeightSum += weight
		if i == 1 {
			for j := 0; j < stateHistory.StateWidth; j++ {
				statisticsVec.SetVec(j, weight*stateHistory.Values.At(i, j))
			}
			continue
		}
		statisticsVec.AddScaledVec(
			statisticsVec,
			weight,
			stateHistory.Values.RowView(i),
		)
	}
	statisticsVec.ScaleVec(1.0/cumulativeWeightSum, statisticsVec)
	return statisticsVec.RawVector().Data
}
