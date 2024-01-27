package models

import (
	"testing"

	"github.com/umbralcalc/learnadex/pkg/learning"
	"github.com/umbralcalc/learnadex/pkg/reweighting"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// newImplementationsAndSimpleLearningConfigForTests creates an implementations
// config which loads a data streamer from the test_file.csv for testing and a
// simple learning config with a single normally-distributed data linking
// log-likelihooda and standard covariance statistics.
func newImplementationsAndSimpleLearningConfigForTests(
	settings *simulator.Settings,
	conditionalProb reweighting.ConditionalProbability,
) (*simulator.Implementations, *learning.LearningConfig) {
	implementations := &simulator.Implementations{
		Iterations:      make([]simulator.Iteration, 0),
		OutputCondition: &simulator.NilOutputCondition{},
		OutputFunction:  &simulator.NilOutputFunction{},
		TerminationCondition: &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: 99,
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
	return implementations, &learning.LearningConfig{
		Objectives:      objectives,
		ObjectiveOutput: &learning.NilObjectiveOutputFunction{},
	}
}

func TestExponentialTimeWeighting(t *testing.T) {
	t.Run(
		"test that the exponential time weighting learning objective evaluates",
		func(t *testing.T) {
			configPath := "exponential_time_weighting_config.yaml"
			settings := simulator.LoadSettingsFromYaml(configPath)
			implementations, config :=
				newImplementationsAndSimpleLearningConfigForTests(
					settings,
					&ExponentialTimeWeightingConditionalProbability{},
				)
			learningObjective := learning.NewObjectiveEvaluator(
				implementations,
				settings,
				config,
			)
			_ = learningObjective.Evaluate(settings.OtherParams)
		},
	)
}
