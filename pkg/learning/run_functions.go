package learning

import "github.com/umbralcalc/stochadex/pkg/simulator"

// RunFilterParamsLearning
func RunFilterParamsLearning(
	config *LearnadexConfig,
	settings *simulator.LoadSettingsConfig,
) []*simulator.OtherParams {
	return config.Optimiser.Algorithm.Run(
		NewLearningObjective(config.Learning, settings),
		settings.OtherParams,
	)
}
