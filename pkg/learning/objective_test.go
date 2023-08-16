package learning

import (
	"testing"

	"github.com/umbralcalc/learnadex/pkg/filter"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// dummyConditionalProbability is just used for testing.
type dummyConditionalProbability struct {
	value float64
}

func (d *dummyConditionalProbability) Configure(
	partitionIndex int,
	settings *simulator.LoadSettingsConfig,
) {
	d.value = settings.OtherParams[partitionIndex].FloatParams["dummy_value"][0]
}

func (d *dummyConditionalProbability) SetParams(
	params *simulator.OtherParams,
) {
	d.value = params.FloatParams["dummy_value"][0]
}

func (d *dummyConditionalProbability) Evaluate(
	currentState []float64,
	pastState []float64,
	currentTime float64,
	pastTime float64,
) float64 {
	return d.value
}

func newLearningConfigForTests(settings *simulator.LoadSettingsConfig) *LearningConfig {
	implementations := &simulator.LoadImplementationsConfig{
		Iterations:      make([]simulator.Iteration, 0),
		OutputCondition: &simulator.NilOutputCondition{},
		OutputFunction:  &simulator.NilOutputFunction{},
		TerminationCondition: &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: 100,
		},
		TimestepFunction: NewMemoryTimestepFunctionFromCsv(
			"test_file.csv",
			0,
			true,
		),
	}
	iteration := NewMemoryIterationFromCsv(
		"test_file.csv",
		0,
		[]int{1, 2, 3},
		true,
	)
	implementations.Iterations = append(implementations.Iterations, iteration)
	anotherIteration := NewMemoryIterationFromCsv(
		"test_file.csv",
		0,
		[]int{1, 2, 3},
		true,
	)
	implementations.Iterations = append(implementations.Iterations, anotherIteration)
	objectives := make([]LogLikelihood, 0)
	firstObjective := &filter.ProbabilityFilterLogLikelihood{
		Prob:       &dummyConditionalProbability{},
		DataLink:   &filter.NormalDataLinkingLogLikelihood{},
		Statistics: &filter.Statistics{},
	}
	firstObjective.Configure(0, settings)
	objectives = append(objectives, firstObjective)
	secondObjective := &filter.ProbabilityFilterLogLikelihood{
		Prob:       &dummyConditionalProbability{},
		DataLink:   &filter.NormalDataLinkingLogLikelihood{},
		Statistics: &filter.Statistics{},
	}
	secondObjective.Configure(1, settings)
	objectives = append(objectives, secondObjective)
	return &LearningConfig{
		Streaming:         implementations,
		StreamingSettings: settings,
		Objectives:        objectives,
		ObjectiveOutput:   &NilObjectiveOutputFunction{},
	}
}

func TestLearningObjective(t *testing.T) {
	t.Run(
		"test that the learning objective runs",
		func(t *testing.T) {
			settings := simulator.NewLoadSettingsConfigFromYaml("test_config.yaml")
			config := newLearningConfigForTests(settings)
			learningObjective := NewLearningObjective(config)
			_ = learningObjective.Evaluate(settings.OtherParams)
		},
	)
}
