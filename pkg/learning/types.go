package learning

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// DataStreamingConfig
type DataStreamingConfig struct {
	DataStreamer         DataStreamer
	TimestepFunction     simulator.TimestepFunction
	TerminationCondition simulator.TerminationCondition
}

// LearningConfig
type LearningConfig struct {
	Streaming  []*DataStreamingConfig
	Objectives []LogLikelihood
}

// OptimiserConfig
type OptimiserConfig struct {
	Algorithm        OptimisationAlgorithm
	ParamsToOptimise []*simulator.OtherParamsMask
}

// LearnadexConfig
type LearnadexConfig struct {
	Learning  *LearningConfig
	Optimiser *OptimiserConfig
}
