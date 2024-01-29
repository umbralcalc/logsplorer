package learning

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// ObjectiveOutputFunction is the interface that must be implemented for outputting
// data from the LearningObjective.
type ObjectiveOutputFunction interface {
	Output(
		partitionIndex int,
		time float64,
		objective float64,
		params *simulator.OtherParams,
	)
}

// NilObjectiveOutputFunction outputs nothing from the LearningObjective.
type NilObjectiveOutputFunction struct{}

func (n *NilObjectiveOutputFunction) Output(
	partitionIndex int,
	time float64,
	objective float64,
	params *simulator.OtherParams,
) {

}

// StdoutObjectiveOutputFunction outputs data to the console from the LearningObjective.
type StdoutObjectiveOutputFunction struct{}

func (s *StdoutObjectiveOutputFunction) Output(
	partitionIndex int,
	time float64,
	objective float64,
	params *simulator.OtherParams,
) {
	fmt.Println(partitionIndex, time, objective, params)
}

// JsonLogEntry is the format in which the logs are serialised when using the
// JsonLogObjectiveOutputFunction.
type JsonLogEntry struct {
	PartitionIndex int                  `json:"partition_index"`
	Time           float64              `json:"time"`
	Objective      float64              `json:"objective"`
	FloatParams    map[string][]float64 `json:"float_params"`
	IntParams      map[string][]int64   `json:"int_params"`
}

// JsonLogObjectiveOutputFunction outputs data to log of json packets from
// the LearningObjective.
type JsonLogObjectiveOutputFunction struct {
	file *os.File
}

func (j *JsonLogObjectiveOutputFunction) Output(
	partitionIndex int,
	time float64,
	objective float64,
	params *simulator.OtherParams,
) {
	logEntry := JsonLogEntry{
		PartitionIndex: partitionIndex,
		Time:           time,
		Objective:      objective,
		FloatParams:    params.FloatParams,
		IntParams:      params.IntParams,
	}
	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Error encoding JSON: %s\n", err)
		panic(err)
	}
	jsonData = append(jsonData, []byte("\n")...)
	_, err = j.file.Write(jsonData)
	if err != nil {
		panic(err)
	}
}

// NewJsonLogObjectiveOutputFunction creates a new JsonLogObjectiveOutputFunction.
func NewJsonLogObjectiveOutputFunction(
	filePath string,
) *JsonLogObjectiveOutputFunction {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Error creating log file:", err)
		panic(err)
	}
	return &JsonLogObjectiveOutputFunction{file: file}
}
