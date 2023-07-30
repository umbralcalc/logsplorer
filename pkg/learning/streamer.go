package learning

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// MemoryIteration provides a stream of data which is already know from a
// separate data source and is held in memory.
type MemoryIteration struct {
	Data map[float64][]float64
}

func (m *MemoryIteration) Configure(
	partitionIndex int,
	settings *simulator.LoadSettingsConfig,
) {
}

func (m *MemoryIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	return m.Data[timestepsHistory.Values.AtVec(0)]
}

// NewMemoryIterationFromCsv creates a new MemoryIteration based on data
// that is read in from the provided csv file and some specified columns
// for time and state.
func NewMemoryIterationFromCsv(
	filePath string,
	timeColumn int,
	stateColumns []int,
	skipHeaderRow bool,
) *MemoryIteration {
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
	var time float64
	data := make(map[float64][]float64)
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
	}
	return &MemoryIteration{Data: data}
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

// NewMemoryTimestepFunctionFromCsv creates a new MemoryTimestepFunction
// based on data that is read in from the provided csv file and some specified
// columns for time and state.
func NewMemoryTimestepFunctionFromCsv(
	filePath string,
	timeColumn int,
	skipHeaderRow bool,
) *MemoryTimestepFunction {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	var increment float64
	timeIncrements := make(map[float64]float64)
	for j, row := range records {
		if skipHeaderRow {
			skipHeaderRow = false
			continue
		}
		time, err := strconv.ParseFloat(row[timeColumn], 64)
		if err != nil {
			fmt.Printf("Error converting string: %v", err)
		}
		if j < len(records)-1 {
			dataPoint, err := strconv.ParseFloat(records[j+1][timeColumn], 64)
			if err != nil {
				fmt.Printf("Error converting string: %v", err)
			}
			increment = dataPoint - time
		}
		timeIncrements[time] = increment
	}
	return &MemoryTimestepFunction{NextIncrements: timeIncrements}
}
