package models

import (
	"testing"

	"github.com/umbralcalc/learnadex/pkg/learning"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

func TestGaussianProcess(t *testing.T) {
	t.Run(
		"test that the Gaussian process learning objective evaluates",
		func(t *testing.T) {
			configPath := "gaussian_process_config.yaml"
			settings := simulator.NewLoadSettingsConfigFromYaml(configPath)
			gaussianProc := &GaussianProcessConditionalProbability{
				Kernel: &ConstantGaussianProcessCovarianceKernel{},
			}
			config := newSimpleLearningConfigForTests(settings, gaussianProc)
			learningObjective := learning.NewLearningObjective(config)
			_ = learningObjective.Evaluate(settings.OtherParams)
		},
	)
}
