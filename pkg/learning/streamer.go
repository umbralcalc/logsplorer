package learning

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

// DataStreamer defines the interface that must be implemented to
// support streaming data from any source to a Learner. The NextValue()
// method should iterate through the data each time it is called and
// must ouput nil when there is no more data.
type DataStreamer interface {
	NextValue() []float64
}

// CsvFileDataStreamer
type CsvFileDataStreamer struct {
	data         [][]float64
	dataLength   int
	currentIndex int
}

func (c *CsvFileDataStreamer) NextValue() []float64 {
	if c.currentIndex == c.dataLength {
		return nil
	}
	nextValue := c.data[c.currentIndex]
	c.currentIndex += 1
	return nextValue
}

func NewCsvFileDataStreamer(filePath string) *CsvFileDataStreamer {
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
	data := make([][]float64, 0)
	for _, row := range records {
		floatRow := make([]float64, 0)
		for _, r := range row {
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
		dataLength:   len(data),
		currentIndex: 0,
	}
}
