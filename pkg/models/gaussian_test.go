package models

import (
	"testing"

	"github.com/umbralcalc/learnadex/pkg/learning"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

func TestGaussian(t *testing.T) {
	t.Run(
		"test that the Gaussian learning objective evaluates",
		func(t *testing.T) {
			configPath := "gaussian_config.yaml"
			settings := simulator.NewLoadSettingsConfigFromYaml(configPath)
			gaussianProc := &GaussianConditionalProbability{
				Kernel: &ConstantGaussianCovarianceKernel{},
			}
			config := newSimpleLearningConfigForTests(settings, gaussianProc)
			learningObjective := learning.NewLearningObjective(config)
			_ = learningObjective.Evaluate(settings.OtherParams)
		},
	)
}
