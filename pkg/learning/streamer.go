package learning

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// DataStreamer defines the interface that must be implemented to
// support streaming data from any source to a LearningObjective.
type DataStreamer interface {
	NextValue(
		timestepsHistory *simulator.CumulativeTimestepsHistory,
	) []float64
}

// CsvFileDataStreamer provides a stream of data that has been read
// in from a csv file.
type CsvFileDataStreamer struct {
	data [][]float64
}

func (c *CsvFileDataStreamer) NextValue(
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	return c.data[int(timestepsHistory.Values.AtVec(0))]
}

// NewCsvFileDataStreamingConfig creates a new DataStreamingConfig based on
// read in csv file data and some specified columns for time and state.
func NewCsvFileDataStreamingConfig(
	filePath string,
	timeColumn int,
	stateColumns []int,
	skipHeaderRow bool,
) *DataStreamingConfig {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	// create this as a faster lookup
	stateColumnsMap := make(map[int]bool)
	for _, column := range stateColumns {
		stateColumnsMap[column] = true
	}

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	data := make([][]float64, 0)
	timeData := make([]float64, 0)
	for _, row := range records {
		if skipHeaderRow {
			skipHeaderRow = false
			continue
		}
		floatRow := make([]float64, 0)
		for i, r := range row {
			if i == timeColumn {
				dataPoint, err := strconv.ParseFloat(r, 64)
				if err != nil {
					fmt.Printf("Error converting string: %v", err)
				}
				timeData = append(timeData, dataPoint)
				continue
			}
			_, ok := stateColumnsMap[i]
			if !ok {
				continue
			}
			dataPoint, err := strconv.ParseFloat(r, 64)
			if err != nil {
				fmt.Printf("Error converting string: %v", err)
			}
			floatRow = append(floatRow, dataPoint)
		}
		data = append(data, floatRow)
	}
	return &DataStreamingConfig{
		DataStreamer:     &CsvFileDataStreamer{data: data},
		TimestepFunction: &simulator.ConstantTimestepFunction{Stepsize: 1},
		TerminationCondition: &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: len(timeData),
		},
	}
}
