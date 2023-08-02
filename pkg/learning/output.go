package learning

import (
	"encoding/json"
	"fmt"
	"net/http"

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

// HttpObjectiveOutputFunction outputs data to an HTTP server from the LearningObjective.
type HttpObjectiveOutputFunction struct {
	writer http.ResponseWriter
}

func (p *HttpObjectiveOutputFunction) Output(
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
	p.writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(p.writer).Encode(outputPacket)
}

// NewHttpObjectiveOutputFunction creates a new HttpObjectiveOutputFunction.
func NewHttpObjectiveOutputFunction(
	writer http.ResponseWriter,
) *HttpObjectiveOutputFunction {
	return &HttpObjectiveOutputFunction{writer: writer}
}
