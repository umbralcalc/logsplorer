package learning

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// LogLikelihood defines the interface that must be implemented in order to
// create an objective function which is applied to the data in an iterative
// fashion by the DataIteration class.
type LogLikelihood interface {
	Configure(partitionIndex int, settings *simulator.LoadSettingsConfig)
	Evaluate(
		params *simulator.OtherParams,
		partitionIndex int,
		stateHistories []*simulator.StateHistory,
		timestepsHistory *simulator.CumulativeTimestepsHistory,
	) float64
}

// IterationWithObjective allows for iteration through a given stochadex simulation
// while evaluating a specified objective function in the form of a cumulative
// log-likelihood as it goes.
type IterationWithObjective struct {
	logLikelihood           LogLikelihood
	iteration               simulator.Iteration
	cumulativeLogLikelihood float64
	burnInSteps             int
	stepsTaken              int
}

func (i *IterationWithObjective) Configure(
	partitionIndex int,
	settings *simulator.LoadSettingsConfig,
) {
	i.logLikelihood.Configure(partitionIndex, settings)
	i.burnInSteps = int(
		settings.OtherParams[partitionIndex].IntParams["burn_in_steps"][0],
	)
	i.cumulativeLogLikelihood = 0.0
	i.stepsTaken = 0
}

func (i *IterationWithObjective) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	i.stepsTaken += 1
	if i.stepsTaken <= i.burnInSteps {
		return i.iteration.Iterate(
			params,
			partitionIndex,
			stateHistories,
			timestepsHistory,
		)
	}
	i.cumulativeLogLikelihood += i.logLikelihood.Evaluate(
		params,
		partitionIndex,
		stateHistories,
		timestepsHistory,
	)
	return i.iteration.Iterate(
		params,
		partitionIndex,
		stateHistories,
		timestepsHistory,
	)
}

// Get the cumulative log-likelihood that has been calculated by iterating
// through the data and applying the LogLikehood.Evaluate method.
func (i *IterationWithObjective) GetObjective() float64 {
	return i.cumulativeLogLikelihood
}

func minKey(m map[float64][]float64) float64 {
	var min float64
	for k := range m {
		min = k
		break
	}
	for k := range m {
		if k < min {
			min = k
		}
	}
	return min
}

// OnlineLearningIteration facilitates online learning by iterating over
// successive parameter values and rerunning the optimiser over the specified
// streaming window every 'refitSteps' number of steps.
type OnlineLearningIteration struct {
	Config           *LearnadexConfig
	mappings         *OptimiserParamsMappings
	windowIterations []*MemoryIteration
	windowEdgeTimes  []float64
	streamerIndices  []int
	refitSteps       int
	burnInSteps      int
	stepsTaken       int
}

func (o *OnlineLearningIteration) Configure(
	partitionIndex int,
	settings *simulator.LoadSettingsConfig,
) {
	o.mappings = NewOptimiserParamsMappings(settings.OtherParams)
	o.streamerIndices = make([]int, 0)
	o.windowEdgeTimes = make([]float64, 0)
	o.windowIterations = make([]*MemoryIteration, 0)
	for _, index := range settings.OtherParams[partitionIndex].
		IntParams["streamer_partition_indices"] {
		o.streamerIndices = append(o.streamerIndices, int(index))
		o.windowEdgeTimes = append(o.windowEdgeTimes, 0.0)
		o.windowIterations = append(
			o.windowIterations,
			&MemoryIteration{Data: make(map[float64][]float64)},
		)
	}
	o.burnInSteps = int(
		settings.OtherParams[partitionIndex].IntParams["burn_in_steps"][0],
	)
	o.refitSteps = int(
		settings.OtherParams[partitionIndex].IntParams["refit_steps"][0],
	)
	o.stepsTaken = 0
}

func (o *OnlineLearningIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	o.stepsTaken += 1
	nextTime := timestepsHistory.Values.AtVec(0) +
		timestepsHistory.NextIncrement
	for i, oldEdgeTime := range o.windowEdgeTimes {
		if o.stepsTaken > o.burnInSteps {
			delete(o.windowIterations[i].Data, oldEdgeTime)
		}
		o.windowEdgeTimes[i] = minKey(o.windowIterations[i].Data)
		o.windowIterations[i].Data[nextTime] =
			stateHistories[o.streamerIndices[i]].Values.RawRowView(0)
		o.Config.Learning.Streaming.Iterations[i] = o.windowIterations[i]
		o.Config.Learning.StreamingSettings.InitTimeValue =
			o.windowEdgeTimes[i]
		o.Config.Learning.StreamingSettings.InitStateValues[i] =
			o.windowIterations[i].Data[o.windowEdgeTimes[i]]
	}
	if o.stepsTaken%o.refitSteps == 0 {
		return stateHistories[partitionIndex].Values.RawRowView(0)
	}
	newParamValues := o.Config.Optimiser.Run(
		NewObjectiveEvaluator(o.Config.Learning),
		o.mappings.UpdateParamsFromOptimiser(
			stateHistories[partitionIndex].Values.RawRowView(0),
			o.Config.Learning.StreamingSettings.OtherParams,
		),
	)
	return o.mappings.GetParamsForOptimiser(newParamValues)
}

// NewOnlineLearningIteration creates a new OnlineLearningIteration.
func NewOnlineLearningIteration(
	config *LearnadexConfig,
) *OnlineLearningIteration {
	return &OnlineLearningIteration{Config: config}
}
