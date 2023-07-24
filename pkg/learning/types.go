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

// ExtraLoadSettings is a yaml-loadable config for the extra settings which exist
// for the learnadex (beyond what is already loaded from stochadex settings).
type ExtraLoadSettings struct {
	ParamsToOptimise []*simulator.OtherParamsMask `yaml:"params_to_optimise"`
}
