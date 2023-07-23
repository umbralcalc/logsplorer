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
			_, times := learning.NewMemoryDataStreamingConfigFromCsv(
				"test_file.csv",
				0,
				[]int{1, 2, 3},
				true,
			)
			settings.OtherParams[0].FloatParams["times"] = times
			flattenedMeans := make([]float64, 0)
			for range settings.OtherParams[0].FloatParams["times"] {
				for i := 0; i < settings.StateWidths[0]; i++ {
					flattenedMeans = append(flattenedMeans, 0.0)
				}
			}
			settings.OtherParams[0].FloatParams["flattened_means_in_time"] = flattenedMeans
			currentRow := 0
			triangleCov := make([]float64, 0)
			for col := currentRow; col < settings.StateWidths[0]; col++ {
				for row := col; row < settings.StateWidths[0]; row++ {
					val := 0.0
					if row == col {
						val = 10.0
					}
					triangleCov = append(triangleCov, val)
				}
				currentRow += 1
			}
			settings.OtherParams[0].FloatParams["upper_triangle_covariance_matrix"] = triangleCov
			gaussianProc := &GaussianProcessConditionalProbability{
				Kernel: &ConstantGaussianProcessCovarianceKernel{},
			}
			gaussianProc.Configure(0, settings)
			config := newSimpleLearningConfigForTests(configPath, settings, gaussianProc)
			learningObjective := learning.NewLearningObjective(config, settings)
			_ = learningObjective.Evaluate(settings.OtherParams)
		},
	)
}
