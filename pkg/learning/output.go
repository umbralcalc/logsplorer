package learning

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// ObjectiveOutputFunction is the interface that must be implemented for outputting
// data from the LearningObjective.
type ObjectiveOutputFunction interface {
	Output(partitionIndex int, objective float64, params *simulator.OtherParams)
}

// NilObjectiveOutputFunction outputs nothing from the LearningObjective.
type NilObjectiveOutputFunction struct{}

func (n *NilObjectiveOutputFunction) Output(
	partitionIndex int,
	objective float64,
	params *simulator.OtherParams,
) {

}

// StdoutObjectiveOutputFunction outputs data to the console from the LearningObjective.
type StdoutObjectiveOutputFunction struct{}

func (s *StdoutObjectiveOutputFunction) Output(
	partitionIndex int,
	objective float64,
	params *simulator.OtherParams,
) {
	fmt.Println(partitionIndex, objective, params)
}

// JsonLogObjectiveOutputFunction outputs data to log of json packets from
// the LearningObjective.
type JsonLogObjectiveOutputFunction struct {
	logger *logrus.Logger
}

func (j *JsonLogObjectiveOutputFunction) Output(
	partitionIndex int,
	objective float64,
	params *simulator.OtherParams,
) {
	outputPacket := struct {
		PartitionIndex int
		Objective      float64
		FloatParams    map[string][]float64
		IntParams      map[string][]int64
	}{
		PartitionIndex: partitionIndex,
		Objective:      objective,
		FloatParams:    params.FloatParams,
		IntParams:      params.IntParams,
	}
	jsonData, err := json.Marshal(outputPacket)
	if err != nil {
		log.Printf("Error encoding JSON: %s\n", err)
		panic(err)
	}
	j.logger.Info(string(jsonData))
}

// NewJsonLogObjectiveOutputFunction creates a new JsonLogObjectiveOutputFunction.
func NewJsonLogObjectiveOutputFunction(
	filePath string,
) *JsonLogObjectiveOutputFunction {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Error creating log file:", err)
		panic(err)
	}
	logger := logrus.New()
	logger.Out = file
	return &JsonLogObjectiveOutputFunction{logger: logger}
}
