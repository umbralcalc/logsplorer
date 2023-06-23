package filter

import (
	"math"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distmv"
	"gonum.org/v1/gonum/stat/distuv"
)

// DataLinkingLogLikelihood
type DataLinkingLogLikelihood interface {
	Evaluate(
		statistics Statistics,
		data []float64,
	) float64
}

// NormalDataLinkingLogLikelihood
type NormalDataLinkingLogLikelihood struct{}

func (n *NormalDataLinkingLogLikelihood) Evaluate(
	statistics Statistics,
	data []float64,
) float64 {
	dist, _ := distmv.NewNormal(
		statistics.GetMean().RawVector().Data,
		statistics.GetCovariance(),
		rand.NewSource(0),
	)
	return dist.LogProb(data)
}

// GammaDataLinkingLogLikelihood
type GammaDataLinkingLogLikelihood struct {
	dist *distuv.Gamma
}

func (g *GammaDataLinkingLogLikelihood) Evaluate(
	statistics Statistics,
	data []float64,
) float64 {
	if g.dist == nil {
		g.dist = &distuv.Gamma{Alpha: 1.0, Beta: 1.0, Src: rand.NewSource(0)}
	}
	logLike := 0.0
	mean := statistics.GetMean()
	for i := 0; i < mean.Len(); i++ {
		g.dist.Beta = mean.AtVec(i) * statistics.GetCovariance().At(i, i)
		g.dist.Alpha = mean.AtVec(i) *
			mean.AtVec(i) / statistics.GetCovariance().At(i, i)
		logLike += g.dist.LogProb(data[i])
	}
	return logLike
}

// PoissonDataLinkingLogLikelihood
type PoissonDataLinkingLogLikelihood struct {
	dist *distuv.Poisson
}

func (p *PoissonDataLinkingLogLikelihood) Evaluate(
	statistics Statistics,
	data []float64,
) float64 {
	if p.dist == nil {
		p.dist = &distuv.Poisson{Lambda: 1.0, Src: rand.NewSource(0)}
	}
	logLike := 0.0
	mean := statistics.GetMean()
	for i := 0; i < mean.Len(); i++ {
		p.dist.Lambda = mean.AtVec(i)
		logLike += p.dist.LogProb(data[i])
	}
	return logLike
}

// NegativeBinomialDataLinkingLogLikelihood
type NegativeBinomialDataLinkingLogLikelihood struct{}

func (n *NegativeBinomialDataLinkingLogLikelihood) Evaluate(
	statistics Statistics,
	data []float64,
) float64 {
	logLike := 0.0
	mean := statistics.GetMean()
	for i := 0; i < mean.Len(); i++ {
		r := mean.AtVec(i) * mean.AtVec(i) /
			(statistics.GetCovariance().At(i, i) - mean.AtVec(i))
		p := mean.AtVec(i) / mean.At(i, i)
		lg1, _ := math.Lgamma(r + data[i])
		lg2, _ := math.Lgamma(data[i] + 1.0)
		lg3, _ := math.Lgamma(data[i])
		logLike += lg1 + lg2 + lg3 + (r * math.Log(p)) +
			(data[i] * math.Log(1.0-p))
	}
	return logLike
}
