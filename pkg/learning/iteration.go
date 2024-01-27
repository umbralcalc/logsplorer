package learning

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// LogLikelihood defines the interface that must be implemented in order to
// create an objective function which is applied to the data in an iterative
// fashion by the DataIteration class.
type LogLikelihood interface {
	Configure(partitionIndex int, settings *simulator.Settings)
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
	settings *simulator.Settings,
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

func copySettingsForPartitions(
	partitionIndices []int,
	settings *simulator.Settings,
) *simulator.Settings {
	settingsCopy := &simulator.Settings{}
	settingsCopy.InitTimeValue = settings.InitTimeValue
	settingsCopy.TimestepsHistoryDepth = settings.TimestepsHistoryDepth
	for _, index := range partitionIndices {
		paramsCopy := *settings.OtherParams[index]
		settingsCopy.OtherParams = append(
			settingsCopy.OtherParams,
			&paramsCopy,
		)
		settingsCopy.InitStateValues = append(
			settingsCopy.InitStateValues,
			settings.InitStateValues[index],
		)
		settingsCopy.Seeds = append(
			settingsCopy.Seeds,
			settings.Seeds[index],
		)
		settingsCopy.StateWidths = append(
			settingsCopy.StateWidths,
			settings.StateWidths[index],
		)
		settingsCopy.StateHistoryDepths = append(
			settingsCopy.StateHistoryDepths,
			settings.StateHistoryDepths[index],
		)
	}
	return settingsCopy
}

// OnlineLearningIteration facilitates online learning by iterating over
// successive parameter values and rerunning the optimiser over the specified
// streaming window every 'refitSteps' number of steps.
type OnlineLearningIteration struct {
	config                  *LearningConfig
	streamerImplementations *simulator.Implementations
	streamerSettings        *simulator.Settings
	streamerMappings        *OptimiserParamsMappings
	streamerIndices         []int
	windowIterations        []*MemoryIteration
	windowTimesteps         *MemoryTimestepFunction
	windowEdgeTimes         []float64
	windowSize              int
	refitSteps              int
	stepsTaken              int
}

func (o *OnlineLearningIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	o.streamerIndices = make([]int, 0)
	o.windowEdgeTimes = make([]float64, 0)
	o.windowIterations = make([]*MemoryIteration, 0)
	for i, index := range settings.OtherParams[partitionIndex].
		IntParams["streamer_partition_indices"] {
		if i == 0 {
			o.windowSize = settings.StateHistoryDepths[index]
		} else {
			if o.windowSize != settings.StateHistoryDepths[index] {
				panic("state_history_depth for streamers must all " +
					"be the same when using online learning")
			}
		}
		o.streamerIndices = append(o.streamerIndices, int(index))
		o.windowEdgeTimes = append(o.windowEdgeTimes, 0.0)
		o.windowIterations = append(
			o.windowIterations,
			&MemoryIteration{Data: make(map[float64][]float64)},
		)
	}
	o.windowTimesteps = &MemoryTimestepFunction{
		NextIncrements: make(map[float64]float64),
	}
	o.streamerImplementations = &simulator.Implementations{
		Iterations: make(
			[]simulator.Iteration,
			len(o.streamerIndices),
		),
		OutputCondition: &simulator.NilOutputCondition{},
		OutputFunction:  &simulator.NilOutputFunction{},
	}
	o.streamerSettings = copySettingsForPartitions(
		o.streamerIndices,
		settings,
	)
	o.streamerMappings = NewOptimiserParamsMappings(
		o.streamerSettings.OtherParams,
	)
	o.refitSteps = int(
		settings.OtherParams[partitionIndex].IntParams["refit_steps"][0])
	o.stepsTaken = 0
}

func (o *OnlineLearningIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	if o.stepsTaken == 0 {
		// reset the initial condition to the ones put into
		// the streamer parameters - easier for the user!
		stateHistories[partitionIndex].Values.SetRow(
			0,
			o.streamerMappings.GetParamsForOptimiser(
				o.streamerSettings.OtherParams,
			),
		)
	}
	o.stepsTaken += 1
	nextTime := timestepsHistory.Values.AtVec(0) +
		timestepsHistory.NextIncrement
	o.windowTimesteps.NextIncrements[nextTime] =
		timestepsHistory.NextIncrement
	for i, oldEdgeTime := range o.windowEdgeTimes {
		if o.stepsTaken > o.windowSize {
			delete(o.windowIterations[i].Data, oldEdgeTime)
			if i == 0 {
				delete(o.windowTimesteps.NextIncrements, oldEdgeTime)
			}
		}
		o.windowEdgeTimes[i] = minKey(o.windowIterations[i].Data)
		o.windowIterations[i].Data[nextTime] =
			stateHistories[o.streamerIndices[i]].Values.RawRowView(0)
		o.streamerImplementations.Iterations[i] = o.windowIterations[i]
		o.streamerImplementations.TerminationCondition =
			&simulator.NumberOfStepsTerminationCondition{
				MaxNumberOfSteps: len(o.windowIterations[i].Data) - 1,
			}
		o.streamerImplementations.TimestepFunction = o.windowTimesteps
		o.streamerSettings.InitTimeValue = o.windowEdgeTimes[i]
		o.streamerSettings.InitStateValues[i] =
			o.windowIterations[i].Data[o.windowEdgeTimes[i]]
	}
	if o.stepsTaken%o.refitSteps != 0 {
		return stateHistories[partitionIndex].Values.RawRowView(0)
	}
	newParamValues := o.config.Optimiser.Run(
		NewObjectiveEvaluator(
			o.streamerImplementations,
			o.streamerSettings,
			o.config,
		),
		o.streamerMappings.UpdateParamsFromOptimiser(
			stateHistories[partitionIndex].Values.RawRowView(0),
			o.streamerSettings.OtherParams,
		),
	)
	return o.streamerMappings.GetParamsForOptimiser(newParamValues)
}

// NewOnlineLearningIteration creates a new OnlineLearningIteration.
func NewOnlineLearningIteration(
	config *LearningConfig,
) *OnlineLearningIteration {
	return &OnlineLearningIteration{config: config}
}
