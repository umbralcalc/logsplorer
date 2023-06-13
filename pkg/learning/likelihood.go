package learning

import (
	"math"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distmv"
	"gonum.org/v1/gonum/stat/distuv"
)

// LogLikelihood
type LogLikelihood interface {
	Evaluate(
		params *simulator.OtherParams,
		partitionIndex int,
		stateHistories []*simulator.StateHistory,
		timestepsHistory *simulator.TimestepsHistory,
	) float64
}

// DataLinkingLogLikelihood
type DataLinkingLogLikelihood interface {
	Evaluate(
		params *simulator.OtherParams,
		data []float64,
	) float64
}

// NormalDataLinkingLogLikelihood
type NormalDataLinkingLogLikelihood struct{}

func (n *NormalDataLinkingLogLikelihood) Evaluate(
	stats Statistics,
	data []float64,
) float64 {
	return distmv.NormalLogProb(
		data,
		stats.Mean.RawVector().Data,
		stats.GetCholeskyCovariance(),
	)
}

// GammaDataLinkingLogLikelihood
type GammaDataLinkingLogLikelihood struct {
	dist *distuv.Gamma
}

func (g *GammaDataLinkingLogLikelihood) Evaluate(
	stats Statistics,
	data []float64,
) float64 {
	if g.dist == nil {
		g.dist = &distuv.Gamma{Alpha: 1.0, Beta: 1.0, Src: rand.NewSource(0)}
	}
	logLike := 0.0
	for i := 0; i < stats.Mean.Len(); i++ {
		g.dist.Beta = stats.Mean.AtVec(i) / stats.Covariance.At(i, i)
		g.dist.Alpha = stats.Mean.AtVec(i) *
			stats.Mean.AtVec(i) / stats.Covariance.At(i, i)
		logLike += g.dist.LogProb(data[i])
	}
	return logLike
}

// PoissonDataLinkingLogLikelihood
type PoissonDataLinkingLogLikelihood struct {
	dist *distuv.Poisson
}

func (p *PoissonDataLinkingLogLikelihood) Evaluate(
	stats Statistics,
	data []float64,
) float64 {
	if p.dist == nil {
		p.dist = &distuv.Poisson{Lambda: 1.0, Src: rand.NewSource(0)}
	}
	logLike := 0.0
	for i := 0; i < stats.Mean.Len(); i++ {
		p.dist.Lambda = stats.Mean.AtVec(i)
		logLike += p.dist.LogProb(data[i])
	}
	return logLike
}

// NegativeBinomialDataLinkingLogLikelihood
type NegativeBinomialDataLinkingLogLikelihood struct{}

func (n *NegativeBinomialDataLinkingLogLikelihood) Evaluate(
	stats Statistics,
	data []float64,
) float64 {
	logLike := 0.0
	for i := 0; i < stats.Mean.Len(); i++ {
		r := stats.Mean.AtVec(i) * stats.Mean.AtVec(i) /
			(stats.Covariance.At(i, i) - stats.Mean.AtVec(i))
		p := stats.Mean.AtVec(i) / stats.Covariance.At(i, i)
		lg1, _ := math.Lgamma(r + data[i])
		lg2, _ := math.Lgamma(data[i] + 1.0)
		lg3, _ := math.Lgamma(data[i])
		logLike += lg1 + lg2 + lg3 + (r * math.Log(p)) +
			(data[i] * math.Log(1.0-p))
	}
	return logLike
}
