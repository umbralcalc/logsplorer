package learning

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// LearningConfig
type LearningConfig struct {
	Streaming  *simulator.LoadImplementationsConfig
	Objectives []LogLikelihood
}

// LearnadexConfig
type LearnadexConfig struct {
	Learning  *LearningConfig
	Optimiser OptimisationAlgorithm
}
