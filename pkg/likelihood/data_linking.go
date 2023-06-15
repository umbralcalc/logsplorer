package likelihood

import (
	"math"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distmv"
	"gonum.org/v1/gonum/stat/distuv"
)

// DataLinkingLogLikelihood
type DataLinkingLogLikelihood interface {
	Evaluate(
		statistics *Statistics,
		data []float64,
	) float64
}

// NormalDataLinkingLogLikelihood
type NormalDataLinkingLogLikelihood struct{}

func (n *NormalDataLinkingLogLikelihood) Evaluate(
	statistics *Statistics,
	data []float64,
) float64 {
	return distmv.NormalLogProb(
		data,
		statistics.Mean.RawVector().Data,
		statistics.GetCholeskyCovariance(),
	)
}

// GammaDataLinkingLogLikelihood
type GammaDataLinkingLogLikelihood struct {
	dist *distuv.Gamma
}

func (g *GammaDataLinkingLogLikelihood) Evaluate(
	statistics *Statistics,
	data []float64,
) float64 {
	if g.dist == nil {
		g.dist = &distuv.Gamma{Alpha: 1.0, Beta: 1.0, Src: rand.NewSource(0)}
	}
	logLike := 0.0
	for i := 0; i < statistics.Mean.Len(); i++ {
		g.dist.Beta = statistics.Mean.AtVec(i) / statistics.Covariance.At(i, i)
		g.dist.Alpha = statistics.Mean.AtVec(i) *
			statistics.Mean.AtVec(i) / statistics.Covariance.At(i, i)
		logLike += g.dist.LogProb(data[i])
	}
	return logLike
}

// PoissonDataLinkingLogLikelihood
type PoissonDataLinkingLogLikelihood struct {
	dist *distuv.Poisson
}

func (p *PoissonDataLinkingLogLikelihood) Evaluate(
	statistics *Statistics,
	data []float64,
) float64 {
	if p.dist == nil {
		p.dist = &distuv.Poisson{Lambda: 1.0, Src: rand.NewSource(0)}
	}
	logLike := 0.0
	for i := 0; i < statistics.Mean.Len(); i++ {
		p.dist.Lambda = statistics.Mean.AtVec(i)
		logLike += p.dist.LogProb(data[i])
	}
	return logLike
}

// NegativeBinomialDataLinkingLogLikelihood
type NegativeBinomialDataLinkingLogLikelihood struct{}

func (n *NegativeBinomialDataLinkingLogLikelihood) Evaluate(
	statistics *Statistics,
	data []float64,
) float64 {
	logLike := 0.0
	for i := 0; i < statistics.Mean.Len(); i++ {
		r := statistics.Mean.AtVec(i) * statistics.Mean.AtVec(i) /
			(statistics.Covariance.At(i, i) - statistics.Mean.AtVec(i))
		p := statistics.Mean.AtVec(i) / statistics.Covariance.At(i, i)
		lg1, _ := math.Lgamma(r + data[i])
		lg2, _ := math.Lgamma(data[i] + 1.0)
		lg3, _ := math.Lgamma(data[i])
		logLike += lg1 + lg2 + lg3 + (r * math.Log(p)) +
			(data[i] * math.Log(1.0-p))
	}
	return logLike
}
