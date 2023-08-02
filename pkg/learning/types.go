package learning

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// LearningConfig specifies how data is streamed into each partitioned
// objective function and what each objective function is.
type LearningConfig struct {
	Streaming       *simulator.LoadImplementationsConfig
	Objectives      []LogLikelihood
	ObjectiveOutput ObjectiveOutputFunction
}

// LearnadexConfig fully configures a learning problem configured for
// the learnadex.
type LearnadexConfig struct {
	Learning  *LearningConfig
	Optimiser OptimisationAlgorithm
}
