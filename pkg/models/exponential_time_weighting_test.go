package models

import (
	"testing"

	"github.com/umbralcalc/learnadex/pkg/learning"
	"github.com/umbralcalc/learnadex/pkg/reweighting"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// newSimpleLearningConfigForTests creates a learning config with a single
// normally-distributed data linking log-likelihood, standard covariance
// statistics and loads a data streamer from the test_file.csv for testing.
func newSimpleLearningConfigForTests(
	settings *simulator.LoadSettingsConfig,
	conditionalProb reweighting.ConditionalProbability,
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
	logLike := &reweighting.ProbabilisticReweightingLogLikelihood{
		Prob:       conditionalProb,
		DataLink:   &reweighting.NormalDataLinkingLogLikelihood{},
		Statistics: &reweighting.Statistics{},
	}
	objectives = append(objectives, logLike)
	return &learning.LearningConfig{
		Streaming:         implementations,
		StreamingSettings: settings,
		Objectives:        objectives,
		ObjectiveOutput:   &learning.NilObjectiveOutputFunction{},
	}
}

func TestExponentialTimeWeighting(t *testing.T) {
	t.Run(
		"test that the exponential time weighting learning objective evaluates",
		func(t *testing.T) {
			configPath := "exponential_time_weighting_config.yaml"
			settings := simulator.NewLoadSettingsConfigFromYaml(configPath)
			config := newSimpleLearningConfigForTests(
				settings,
				&ExponentialTimeWeightingConditionalProbability{},
			)
			learningObjective := learning.NewLearningObjective(config)
			_ = learningObjective.Evaluate(settings.OtherParams)
		},
	)
}
