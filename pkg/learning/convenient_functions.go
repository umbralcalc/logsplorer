package learning

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// RunFilterParamsLearning is a convenient function which runs a
// an optimisation of the probability filter for a particular dataset
// which is provided by the data streamer.
func RunFilterParamsLearning(
	config *LearnadexConfig,
	settings *simulator.LoadSettingsConfig,
) []*simulator.OtherParams {
	return config.Optimiser.Algorithm.Run(
		NewLearningObjective(config.Learning, settings),
		settings.OtherParams,
		config.Optimiser.ParamsToOptimise,
	)
}
