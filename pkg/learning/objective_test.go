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

func TestLearningObjective(t *testing.T) {
	t.Run(
		"test that the learning objective runs",
		func(t *testing.T) {
			streamingConfigs := make([]*DataStreamingConfig, 0)
			streamingConfigs = append(
				streamingConfigs,
				NewCsvFileDataStreamingConfig(
					"test_file.csv",
					0,
					[]int{1, 2, 3},
					true,
				),
			)
			streamingConfigs = append(
				streamingConfigs,
				NewCsvFileDataStreamingConfig(
					"test_file.csv",
					0,
					[]int{1, 2, 3},
					true,
				),
			)
			objectives := make([]LogLikelihood, 0)
			objectives = append(
				objectives,
				&filter.ProbabilityFilterLogLikelihood{
					Prob: &dummyConditionalProbability{},
					DataLink: &filter.NormalDataLinkingLogLikelihood{
						Src: rand.NewSource(1234),
					},
					Statistics: &filter.StandardCovarianceStatistics{},
				},
			)
			objectives = append(
				objectives,
				&filter.ProbabilityFilterLogLikelihood{
					Prob: &dummyConditionalProbability{},
					DataLink: &filter.NormalDataLinkingLogLikelihood{
						Src: rand.NewSource(1234),
					},
					Statistics: &filter.StandardCovarianceStatistics{},
				},
			)
			config := &LearningConfig{
				Streaming:  streamingConfigs,
				Objectives: objectives,
			}
			settings := simulator.NewLoadSettingsConfigFromYaml("test_config.yaml")
			learningObjective := NewLearningObjective(config, settings)
			_ = learningObjective.Evaluate(settings.OtherParams)
			learningObjective.ResetIterators()
		},
	)
}
