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

// MemoryDataStreamer provides a stream of data which is already know from a separate
// data source and is held in memory.
type MemoryDataStreamer struct {
	Data map[float64][]float64
}

func (m *MemoryDataStreamer) NextValue(
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	return m.Data[timestepsHistory.Values.AtVec(0)]
}

// MemoryTimestepFunction provides a stream of timesteps which already known from
// a separate data source and is held in memory.
type MemoryTimestepFunction struct {
	NextIncrements map[float64]float64
}

func (m *MemoryTimestepFunction) SetNextIncrement(
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) *simulator.CumulativeTimestepsHistory {
	timestepsHistory.NextIncrement = m.NextIncrements[timestepsHistory.Values.AtVec(0)]
	return timestepsHistory
}

// NewMemoryDataStreamingConfigFromCsv creates a new DataStreamingConfig for a
// MemoryDataStreamer based on data that is read in from the provided csv file
// and some specified columns for time and state.
func NewMemoryDataStreamingConfigFromCsv(
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
	var time, increment float64
	data := make(map[float64][]float64)
	timeIncrements := make(map[float64]float64)
	for j, row := range records {
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
				time = dataPoint
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
		data[time] = floatRow
		if j < len(records)-1 {
			dataPoint, err := strconv.ParseFloat(records[j+1][timeColumn], 64)
			if err != nil {
				fmt.Printf("Error converting string: %v", err)
			}
			increment = dataPoint - time
		}
		timeIncrements[time] = increment
	}
	return &DataStreamingConfig{
		DataStreamer:     &MemoryDataStreamer{Data: data},
		TimestepFunction: &MemoryTimestepFunction{NextIncrements: timeIncrements},
		TerminationCondition: &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: len(timeIncrements),
		},
	}
}
