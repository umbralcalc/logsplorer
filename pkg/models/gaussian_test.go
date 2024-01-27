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
			settings := simulator.LoadSettingsFromYaml(configPath)
			gaussianProc := &GaussianConditionalProbability{
				Kernel: &ConstantGaussianCovarianceKernel{},
			}
			implementations, config :=
				newImplementationsAndSimpleLearningConfigForTests(
					settings,
					gaussianProc,
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
