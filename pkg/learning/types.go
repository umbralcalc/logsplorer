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

// LearnerConfig
type LearnerConfig struct {
	Streaming  []DataStreamingConfig
	Objectives []LogLikelihood
}

// OptimiserConfig
type OptimiserConfig struct {
	Algorithm    OptimisationAlgorithm
	HistoryDepth int
}

// LearnadexConfig
type LearnadexConfig struct {
	Learners  []*LearnerConfig
	Optimiser *OptimiserConfig
}

// LearnerInputMessage
type LearnerInputMessage struct {
	NewParams []*simulator.OtherParams
}

// LearnerOutputMessage
type LearnerOutputMessage struct {
	Objective float64
}
