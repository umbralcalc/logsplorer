package learning

// LearnerConfig
type LearnerConfig struct {
}

// OptimiserConfig
type OptimiserConfig struct {
}

// DataStreamerConfig
type DataStreamerConfig struct {
}

// LearnadexConfig
type LearnadexConfig struct {
	Learners     []*LearnerConfig
	Optimiser    *OptimiserConfig
	DataStreamer *DataStreamerConfig
}
