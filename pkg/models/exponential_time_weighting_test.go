package models

import (
	"testing"

	"github.com/umbralcalc/learnadex/pkg/filter"
	"github.com/umbralcalc/learnadex/pkg/learning"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// newSimpleLearningConfigForTests creates a learning config with a single
// normally-distributed data linking log-likelihood, standard covariance
// statistics and loads a data streamer from the test_file.csv for testing.
func newSimpleLearningConfigForTests(
	extraSettingsConfigPath string,
	settings *simulator.LoadSettingsConfig,
	conditionalProb filter.ConditionalProbability,
) *learning.LearningConfig {
	implementations := &simulator.LoadImplementationsConfig{
		Iterations:      make([]simulator.Iteration, 0),
		OutputCondition: &simulator.NilOutputCondition{},
		OutputFunction:  &simulator.NilOutputFunction{},
		TerminationCondition: &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: 100,
		},
		TimestepFunction: learning.NewMemoryTimestepFunctionFromCsv(
			"test_file.csv",
			0,
			true,
		),
	}
	iteration := learning.NewMemoryIterationFromCsv(
		"test_file.csv",
		0,
		[]int{1, 2, 3},
		true,
	)
	implementations.Iterations = append(implementations.Iterations, iteration)
	objectives := make([]learning.LogLikelihood, 0)
	logLike := &filter.ProbabilityFilterLogLikelihood{
		Prob:       conditionalProb,
		DataLink:   &filter.NormalDataLinkingLogLikelihood{},
		Statistics: &filter.Statistics{},
	}
	objectives = append(objectives, logLike)
	return &learning.LearningConfig{
		Streaming:       implementations,
		Objectives:      objectives,
		ObjectiveOutput: &learning.NilObjectiveOutputFunction{},
	}
}

func TestExponentialTimeWeighting(t *testing.T) {
	t.Run(
		"test that the exponential time weighting learning objective evaluates",
		func(t *testing.T) {
			configPath := "exponential_time_weighting_config.yaml"
			settings := simulator.NewLoadSettingsConfigFromYaml(configPath)
			config := newSimpleLearningConfigForTests(
				configPath,
				settings,
				&ExponentialTimeWeightingConditionalProbability{},
			)
			learningObjective := learning.NewLearningObjective(config, settings)
			_ = learningObjective.Evaluate(settings.OtherParams)
		},
	)
}
