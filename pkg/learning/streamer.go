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
// support streaming data from any source to a Learner.
type DataStreamer interface {
	Reset()
	NextValue() []float64
}

// CsvFileDataStreamer provides a stream of data that has been read
// in from a csv file.
type CsvFileDataStreamer struct {
	data         [][]float64
	currentIndex int
}

func (c *CsvFileDataStreamer) Reset() {
	c.currentIndex = 0
}

func (c *CsvFileDataStreamer) NextValue() []float64 {
	nextValue := c.data[c.currentIndex]
	c.currentIndex += 1
	return nextValue
}

// NewCsvFileDataStreamer creates a new CsvFileDataStreamer given
// the path to the file and a list of column indices to read in.
func NewCsvFileDataStreamer(
	filePath string,
	columns []int,
) *CsvFileDataStreamer {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	// create this as a faster lookup
	columnsMap := make(map[int]bool)
	for _, column := range columns {
		columnsMap[column] = true
	}

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	data := make([][]float64, 0)
	for _, row := range records {
		floatRow := make([]float64, 0)
		for i, r := range row {
			_, ok := columnsMap[i]
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
	return &CsvFileDataStreamer{
		data:         data,
		currentIndex: 0,
	}
}

// CsvFileTimestepFunction provides a stream of timestep values that has been read
// in from a csv file.
type CsvFileTimestepFunction struct {
	data         []float64
	currentIndex int
}

func (c *CsvFileTimestepFunction) SetNextIncrement(
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) *simulator.CumulativeTimestepsHistory {
	timestepsHistory.NextIncrement = c.data[c.currentIndex]
	c.currentIndex += 1
	return timestepsHistory
}

// CsvFileTerminationCondition determines when the iterations over the data should
// stop based on having reached the end of the file.
type CsvFileTerminationCondition struct {
	dataLength int
}

func (c *CsvFileTerminationCondition) Terminate(
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
	overallTimesteps int,
) bool {
	if overallTimesteps == c.dataLength {
		return true
	}
	return false
}

// NewCsvFileDataStreamingConfig creates a new DataStreamingConfig based on
// read in csv file data and some specified columns for time and state.
func NewCsvFileDataStreamingConfig(
	filePath string,
	timeColumn int,
	stateColumns []int,
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
		DataStreamer: &CsvFileDataStreamer{
			data:         data,
			currentIndex: 0,
		},
		TimestepFunction: &CsvFileTimestepFunction{
			data:         timeData,
			currentIndex: 0,
		},
		TerminationCondition: &CsvFileTerminationCondition{
			dataLength: len(timeData),
		},
	}
}
