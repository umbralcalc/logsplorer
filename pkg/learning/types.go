package learning

// LearningConfig specifies the learning objectives for each partitioned
// objective function and how the objective values can be output.
type LearningConfig struct {
	Objectives      []LogLikelihood
	ObjectiveOutput ObjectiveOutputFunction
	Optimiser       OptimisationAlgorithm
}
