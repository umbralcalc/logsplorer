package learning

import (
	"testing"

	"github.com/umbralcalc/learnadex/pkg/models"
	"github.com/umbralcalc/learnadex/pkg/reweighting"
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/optimize"
)

func newOnlineLearningImplementationsForTests(
	settings *simulator.Settings,
) *simulator.Implementations {
	implementations := &simulator.Implementations{
		Iterations:      make([]simulator.Iteration, 0),
		OutputCondition: &simulator.NilOutputCondition{},
		OutputFunction:  &simulator.NilOutputFunction{},
		TerminationCondition: &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: 99,
		},
		TimestepFunction: NewMemoryTimestepFunctionFromCsv(
			"test_file.csv",
			0,
			true,
		),
	}
	iteration := NewMemoryIterationFromCsv(
		"test_file.csv",
		[]int{1, 2, 3},
		true,
	)
	objective := &reweighting.ProbabilisticReweightingLogLikelihood{
		Prob:       &models.ExponentialTimeWeightingConditionalProbability{},
		DataLink:   &reweighting.NormalDataLinkingLogLikelihood{},
		Statistics: &reweighting.Statistics{},
	}
	implementations.Iterations = append(implementations.Iterations, iteration)
	learningConfig := &LearningConfig{
		Objectives:      []LogLikelihood{objective},
		ObjectiveOutput: &NilObjectiveOutputFunction{},
		Optimiser: &GonumOptimisationAlgorithm{
			Method:   &optimize.NelderMead{},
			Settings: &optimize.Settings{Concurrent: 10},
		},
	}
	implementations.Iterations = append(
		implementations.Iterations,
		NewOnlineLearningIteration(learningConfig),
	)
	for i, iteration := range implementations.Iterations {
		iteration.Configure(i, settings)
	}
	return implementations
}

func TestOnlineLearningIteration(t *testing.T) {
	t.Run(
		"test that the online learning iterator runs",
		func(t *testing.T) {
			settings := simulator.LoadSettingsFromYaml("test_online_config.yaml")
			implementations := newOnlineLearningImplementationsForTests(
				settings,
			)
			coordinator := simulator.NewPartitionCoordinator(
				settings,
				implementations,
			)
			coordinator.Run()
		},
	)
}
