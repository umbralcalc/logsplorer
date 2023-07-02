package learning

import (
	"testing"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"gonum.org/v1/gonum/mat"
)

func TestCsvFileDataStreamer(t *testing.T) {
	t.Run(
		"test that the file streamer works",
		func(t *testing.T) {
			config, _ := NewMemoryDataStreamingConfigFromCsv("test_file.csv", 0, []int{1, 2, 3}, true)
			_ = config.DataStreamer.NextValue(
				&simulator.CumulativeTimestepsHistory{
					NextIncrement:     1.0,
					Values:            mat.NewVecDense(2, []float64{1.0, 0.0}),
					StateHistoryDepth: 2,
				},
			)
		},
	)
}
