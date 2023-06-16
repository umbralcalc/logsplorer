package learning

import "github.com/umbralcalc/learnadex/pkg/likelihood"

// LearnerConfig
type LearnerConfig struct {
	LogLikelihood likelihood.LogLikelihood
	DataStreamer  DataStreamer
}

// OptimiserConfig
type OptimiserConfig struct {
}

// LearnadexConfig
type LearnadexConfig struct {
	Learners  []*LearnerConfig
	Optimiser *OptimiserConfig
}
