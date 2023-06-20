package learning

import (
	"github.com/umbralcalc/learnadex/pkg/likelihood"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// DataStreamingConfig
type DataStreamingConfig struct {
	DataStreamer         DataStreamer
	TimestepFunction     simulator.TimestepFunction
	TerminationCondition simulator.TerminationCondition
}

// LearnerConfig
type LearnerConfig struct {
	Streaming  []DataStreamingConfig
	Objectives []likelihood.LogLikelihood
}

// OptimiserConfig
type OptimiserConfig struct {
}

// LearnadexConfig
type LearnadexConfig struct {
	Learners  []*LearnerConfig
	Optimiser *OptimiserConfig
}

// LearnerInputMessage
type LearnerInputMessage struct {
}

// LearnerOutputMessage
type LearnerOutputMessage struct {
	Objectives []float64
}
