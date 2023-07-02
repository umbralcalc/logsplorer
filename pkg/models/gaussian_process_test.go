package models

import (
	"testing"

	"github.com/umbralcalc/learnadex/pkg/learning"
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/mat"
)

func TestGaussianProcess(t *testing.T) {
	t.Run(
		"test that the Gaussian process learning objective evaluates",
		func(t *testing.T) {
			configPath := "gaussian_process_config.yaml"
			settings := simulator.NewLoadSettingsConfigFromYaml(configPath)
			_, times := learning.NewMemoryDataStreamingConfigFromCsv(
				"test_file.csv",
				0,
				[]int{1, 2, 3},
				true,
			)
			settings.OtherParams[0].FloatParams["times"] = times
			for range settings.OtherParams[0].FloatParams["times"] {
				for i := 0; i < settings.StateWidths[0]; i++ {
					settings.OtherParams[0].FloatParams["flattened_means_in_time"] = append(
						settings.OtherParams[0].FloatParams["flattened_means_in_time"],
						0.0,
					)
				}
			}
			currentRow := 0
			for col := currentRow; col < settings.StateWidths[0]; col++ {
				for row := col; row < settings.StateWidths[0]; row++ {
					val := 0.0
					if row == col {
						val = 1.0
					}
					settings.OtherParams[0].FloatParams["upper_triangle_covariance_matrix"] = append(
						settings.OtherParams[0].FloatParams["upper_triangle_covariance_matrix"],
						val,
					)
				}
				currentRow += 1
			}
			config := newSimpleLearningConfigForTests(
				configPath,
				settings,
				NewGaussianProcessConditionalProbability(
					&ConstantGaussianProcessCovarianceKernel{
						covMatrix:  mat.NewSymDense(settings.StateWidths[0], nil),
						stateWidth: settings.StateWidths[0],
					},
					settings.OtherParams[0].FloatParams["times"],
					settings.StateWidths[0],
				),
			)
			learningObjective := learning.NewLearningObjective(config, settings)
			_ = learningObjective.Evaluate(settings.OtherParams)
			learningObjective.ResetIterators()
		},
	)
}
