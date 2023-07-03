package models

import (
	"testing"

	"github.com/umbralcalc/learnadex/pkg/filter"
	"github.com/umbralcalc/learnadex/pkg/learning"
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"golang.org/x/exp/rand"
)

// newSimpleLearningConfigForTests creates a learning config with a single
// normally-distributed data linking log-likelihood, standard covariance
// statistics and loads a data streamer from the test_file.csv for testing.
func newSimpleLearningConfigForTests(
	extraSettingsConfigPath string,
	settings *simulator.LoadSettingsConfig,
	conditionalProb filter.ConditionalProbability,
) *learning.LearningConfig {
	extraSettings := learning.NewExtraLoadSettingsConfigFromYaml(extraSettingsConfigPath)
	streamingConfigs := make([]*learning.DataStreamingConfig, 0)
	streamingConfig, _ := learning.NewMemoryDataStreamingConfigFromCsv(
		"test_file.csv",
		0,
		[]int{1, 2, 3},
		true,
	)
	streamingConfigs = append(streamingConfigs, streamingConfig)
	objectives := make([]learning.LogLikelihood, 0)
	objectives = append(
		objectives,
		&filter.ProbabilityFilterLogLikelihood{
			Prob: conditionalProb,
			DataLink: &filter.NormalDataLinkingLogLikelihood{
				Src: rand.NewSource(settings.Seeds[0]),
			},
			Statistics: &filter.StandardCovarianceStatistics{},
		},
	)
	return &learning.LearningConfig{
		BurnInSteps: extraSettings.BurnInSteps,
		Streaming:   streamingConfigs,
		Objectives:  objectives,
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
