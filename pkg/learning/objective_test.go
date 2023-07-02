package learning

import (
	"testing"

	"github.com/umbralcalc/learnadex/pkg/filter"
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"golang.org/x/exp/rand"
)

// dummyConditionalProbability is just used for testing.
type dummyConditionalProbability struct {
	value float64
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
	extraSettings := NewExtraLoadSettingsConfigFromYaml("test_config.yaml")
	streamingConfigs := make([]*DataStreamingConfig, 0)
	streamingConfig, _ := NewMemoryDataStreamingConfigFromCsv(
		"test_file.csv",
		0,
		[]int{1, 2, 3},
		true,
	)
	streamingConfigs = append(streamingConfigs, streamingConfig)
	anotherStreamingConfig, _ := NewMemoryDataStreamingConfigFromCsv(
		"test_file.csv",
		0,
		[]int{1, 2, 3},
		true,
	)
	streamingConfigs = append(streamingConfigs, anotherStreamingConfig)
	objectives := make([]LogLikelihood, 0)
	objectives = append(
		objectives,
		&filter.ProbabilityFilterLogLikelihood{
			Prob: &dummyConditionalProbability{},
			DataLink: &filter.NormalDataLinkingLogLikelihood{
				Src: rand.NewSource(settings.Seeds[0]),
			},
			Statistics: &filter.StandardCovarianceStatistics{},
		},
	)
	objectives = append(
		objectives,
		&filter.ProbabilityFilterLogLikelihood{
			Prob: &dummyConditionalProbability{},
			DataLink: &filter.NormalDataLinkingLogLikelihood{
				Src: rand.NewSource(settings.Seeds[1]),
			},
			Statistics: &filter.StandardCovarianceStatistics{},
		},
	)
	return &LearningConfig{
		BurnInSteps: extraSettings.BurnInSteps,
		Streaming:   streamingConfigs,
		Objectives:  objectives,
	}
}

func TestLearningObjective(t *testing.T) {
	t.Run(
		"test that the learning objective runs",
		func(t *testing.T) {
			settings := simulator.NewLoadSettingsConfigFromYaml("test_config.yaml")
			config := newLearningConfigForTests(settings)
			learningObjective := NewLearningObjective(config, settings)
			_ = learningObjective.Evaluate(settings.OtherParams)
			learningObjective.ResetIterators()
		},
	)
}
