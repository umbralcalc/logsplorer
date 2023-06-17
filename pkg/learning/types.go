package learning

import "github.com/umbralcalc/learnadex/pkg/likelihood"

// LearnerConfig
type LearnerConfig struct {
	Iterations []DataIteration
	Streamer   DataStreamer
	Objective  likelihood.LogLikelihood
}

// OptimiserConfig
type OptimiserConfig struct {
}

// LearnadexConfig
type LearnadexConfig struct {
	Learner   *LearnerConfig
	Optimiser *OptimiserConfig
}
