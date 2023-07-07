package filter

import (
	"math"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distmv"
	"gonum.org/v1/gonum/stat/distuv"
)

// DataLinkingLogLikelihood is the interface that must be implemented in
// order to create a likelihood that connects derived statistics from the
// probability filter to observed actual data values.
type DataLinkingLogLikelihood interface {
	Evaluate(statistics *Statistics, data []float64) float64
	GenerateNewSamples(statistics *Statistics) []float64
}

// NormalDataLinkingLogLikelihood assumes the real data are well described
// by a normal distribution, given the statistics provided by the filter.
type NormalDataLinkingLogLikelihood struct {
	Src rand.Source
}

func (n *NormalDataLinkingLogLikelihood) Evaluate(
	statistics *Statistics,
	data []float64,
) float64 {
	dist, ok := distmv.NewNormal(
		statistics.Mean.RawVector().Data,
		statistics.Covariance,
		n.Src,
	)
	if !ok {
		return math.NaN()
	}
	return dist.LogProb(data)
}

func (n *NormalDataLinkingLogLikelihood) GenerateNewSamples(
	statistics *Statistics,
) []float64 {
	dist, ok := distmv.NewNormal(
		statistics.Mean.RawVector().Data,
		statistics.Covariance,
		n.Src,
	)
	if !ok {
		values := make([]float64, 0)
		for i := 0; i < statistics.Mean.Len(); i++ {
			values = append(values, math.NaN())
		}
		return values
	}
	return dist.Rand(nil)
}

// GammaDataLinkingLogLikelihood assumes the real data are well described
// by a gamma distribution, given the statistics provided by the filter.
type GammaDataLinkingLogLikelihood struct {
	Src rand.Source
}

func (g *GammaDataLinkingLogLikelihood) Evaluate(
	statistics *Statistics,
	data []float64,
) float64 {
	dist := &distuv.Gamma{Alpha: 1.0, Beta: 1.0, Src: g.Src}
	logLike := 0.0
	mean := statistics.Mean
	for i := 0; i < mean.Len(); i++ {
		dist.Beta = mean.AtVec(i) * statistics.Covariance.At(i, i)
		dist.Alpha = mean.AtVec(i) *
			mean.AtVec(i) / statistics.Covariance.At(i, i)
		logLike += dist.LogProb(data[i])
	}
	return logLike
}

func (g *GammaDataLinkingLogLikelihood) GenerateNewSamples(
	statistics *Statistics,
) []float64 {
	samples := make([]float64, 0)
	dist := &distuv.Gamma{Alpha: 1.0, Beta: 1.0, Src: g.Src}
	mean := statistics.Mean
	for i := 0; i < mean.Len(); i++ {
		dist.Beta = mean.AtVec(i) * statistics.Covariance.At(i, i)
		dist.Alpha = mean.AtVec(i) *
			mean.AtVec(i) / statistics.Covariance.At(i, i)
		samples = append(samples, dist.Rand())
	}
	return samples
}

// PoissonDataLinkingLogLikelihood assumes the real data are well described
// by a Poisson distribution, given the statistics provided by the filter.
type PoissonDataLinkingLogLikelihood struct {
	Src rand.Source
}

func (p *PoissonDataLinkingLogLikelihood) Evaluate(
	statistics *Statistics,
	data []float64,
) float64 {
	dist := &distuv.Poisson{Lambda: 1.0, Src: p.Src}
	logLike := 0.0
	mean := statistics.Mean
	for i := 0; i < mean.Len(); i++ {
		dist.Lambda = mean.AtVec(i)
		logLike += dist.LogProb(data[i])
	}
	return logLike
}

func (p *PoissonDataLinkingLogLikelihood) GenerateNewSamples(
	statistics *Statistics,
) []float64 {
	samples := make([]float64, 0)
	dist := &distuv.Poisson{Lambda: 1.0, Src: p.Src}
	mean := statistics.Mean
	for i := 0; i < mean.Len(); i++ {
		dist.Lambda = mean.AtVec(i)
		samples = append(samples, dist.Rand())
	}
	return samples
}

// NegativeBinomialDataLinkingLogLikelihood assumes the real data are well
// described by a negative binomial distribution, given the statistics
// provided by the filter.
type NegativeBinomialDataLinkingLogLikelihood struct {
	Src rand.Source
}

func (n *NegativeBinomialDataLinkingLogLikelihood) Evaluate(
	statistics *Statistics,
	data []float64,
) float64 {
	logLike := 0.0
	mean := statistics.Mean
	for i := 0; i < mean.Len(); i++ {
		r := mean.AtVec(i) * mean.AtVec(i) /
			(statistics.Covariance.At(i, i) - mean.AtVec(i))
		p := mean.AtVec(i) / statistics.Covariance.At(i, i)
		lg1, _ := math.Lgamma(r + data[i])
		lg2, _ := math.Lgamma(data[i] + 1.0)
		lg3, _ := math.Lgamma(data[i])
		logLike += lg1 + lg2 + lg3 + (r * math.Log(p)) +
			(data[i] * math.Log(1.0-p))
	}
	return logLike
}

func (n *NegativeBinomialDataLinkingLogLikelihood) GenerateNewSamples(
	statistics *Statistics,
) []float64 {
	samples := make([]float64, 0)
	distPoisson := &distuv.Poisson{Lambda: 1.0, Src: n.Src}
	distGamma := &distuv.Gamma{Alpha: 1.0, Beta: 1.0, Src: n.Src}
	mean := statistics.Mean
	for i := 0; i < mean.Len(); i++ {
		distGamma.Beta = 1.0 /
			((statistics.Covariance.At(i, i) / mean.AtVec(i)) - 1.0)
		distGamma.Alpha = mean.AtVec(i) * distGamma.Beta
		distPoisson.Lambda = distGamma.Rand()
		samples = append(samples, distPoisson.Rand())
	}
	return samples
}
