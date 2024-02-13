package learning

import (
	"github.com/umbralcalc/learnadex/pkg/params"
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
}

func (i *IterationWithObjective) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	i.logLikelihood.Configure(partitionIndex, settings)
	i.burnInSteps = settings.StateHistoryDepths[partitionIndex]
	i.cumulativeLogLikelihood = 0.0
}

func (i *IterationWithObjective) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	if timestepsHistory.CurrentStepNumber <= i.burnInSteps {
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

// OnlineLearningIteration facilitates online learning by iterating over
// successive parameter values and rerunning the optimiser over the specified
// streaming window every 'refitSteps' number of steps.
type OnlineLearningIteration struct {
	config                *LearningConfig
	learnerStreamSettings *simulator.Settings
	learnerStreamMappings *params.OptimiserParamsMappings
	learnerStreamIndices  []int
	refitSteps            int
}

func (o *OnlineLearningIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	o.learnerStreamIndices = make([]int, 0)
	learnerHistoryDepths := make([]int, 0)
	var windowLength int
	for i, index := range settings.OtherParams[partitionIndex].
		IntParams["streamer_partition_indices"] {
		if i == 0 {
			windowLength = settings.StateHistoryDepths[index]
		}
		if (windowLength != settings.StateHistoryDepths[index]) ||
			(settings.StateHistoryDepths[index] !=
				settings.StateHistoryDepths[partitionIndex]) {
			panic("all state history depths for " +
				"streamer_partition_indices must be the same" +
				"as for the online learning iteration - " +
				"use the 'learner_history_depths' parameter if " +
				"you want to vary each learner's window size")
		}
		o.learnerStreamIndices = append(o.learnerStreamIndices, int(index))
		learnerHistoryDepths = append(
			learnerHistoryDepths,
			int(settings.OtherParams[partitionIndex].
				IntParams["learner_history_depths"][i]),
		)
	}
	o.learnerStreamSettings = params.CopySettingsForPartitions(
		o.learnerStreamIndices,
		settings,
	)
	o.learnerStreamSettings.StateHistoryDepths = learnerHistoryDepths
	o.learnerStreamMappings = params.NewOptimiserParamsMappings(
		o.learnerStreamSettings.OtherParams,
	)
	o.refitSteps = int(
		settings.OtherParams[partitionIndex].IntParams["refit_steps"][0])
}

func (o *OnlineLearningIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	if timestepsHistory.CurrentStepNumber == 1 {
		// reset the initial condition to the ones put into
		// the learner stream parameters - easier for the user!
		stateHistories[partitionIndex].Values.SetRow(
			0,
			o.learnerStreamMappings.GetParamsForOptimiser(
				o.learnerStreamSettings.OtherParams,
			),
		)
	}
	if timestepsHistory.CurrentStepNumber%o.refitSteps != 0 {
		return stateHistories[partitionIndex].Values.RawRowView(0)
	}
	windowLength := stateHistories[partitionIndex].StateHistoryDepth
	learnerStreamImplementations := &simulator.Implementations{
		Iterations: make(
			[][]simulator.Iteration,
			len(o.learnerStreamIndices),
		),
		OutputCondition: &simulator.NilOutputCondition{},
		OutputFunction:  &simulator.NilOutputFunction{},
	}
	for i, index := range o.learnerStreamIndices {
		learnerStreamImplementations.Iterations[i] = []simulator.Iteration{
			&MemoryIteration{Data: stateHistories[index]},
		}
		o.learnerStreamSettings.InitStateValues[i] =
			stateHistories[index].Values.RawRowView(
				windowLength - 1,
			)
	}
	learnerStreamImplementations.TerminationCondition =
		&simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: windowLength - 1,
		}
	learnerStreamImplementations.TimestepFunction =
		&MemoryTimestepFunction{Data: timestepsHistory}
	o.learnerStreamSettings.InitTimeValue =
		timestepsHistory.Values.AtVec(windowLength - 1)
	newParamValues := o.config.Optimiser.Run(
		NewObjectiveEvaluator(
			learnerStreamImplementations,
			o.learnerStreamSettings,
			o.config,
		),
		o.learnerStreamMappings.UpdateParamsFromOptimiser(
			stateHistories[partitionIndex].Values.RawRowView(0),
			o.learnerStreamSettings.OtherParams,
		),
	)
	return o.learnerStreamMappings.GetParamsForOptimiser(newParamValues)
}

// NewOnlineLearningIteration creates a new OnlineLearningIteration.
func NewOnlineLearningIteration(
	config *LearningConfig,
) *OnlineLearningIteration {
	return &OnlineLearningIteration{config: config}
}
