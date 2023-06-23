package learning

import (
	"sync"

	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// Learner
type Learner struct {
	config          *LearnerConfig
	stochadexConfig *simulator.StochadexConfig
	dataIterations  []*DataIterator
}

func (l *Learner) ComputeObjective(
	newParams []*simulator.OtherParams,
) float64 {
	var wg sync.WaitGroup

	// set the incoming new params for each state partition
	for index := range l.stochadexConfig.Partitions {
		l.stochadexConfig.Partitions[index].Params.Other = newParams[index]
	}

	// instantiate a new batch of data iterators via the stochadex
	coordinator := simulator.NewPartitionCoordinator(l.stochadexConfig)

	// run the iterations over the data and terminate the for loop
	// when the end of data condition has been met
	for !coordinator.ReadyToTerminate() {
		coordinator.Step(&wg)
	}

	// store the objective value found in this step
	objective := 0.0
	for _, iteration := range l.dataIterations {
		objective += iteration.GetObjective()
	}

	// reset the stateful data iterators to be used again
	l.ResetIterators()

	return objective
}

func (l *Learner) ReceiveAndSendObjectives(
	inputChannel <-chan *LearnerInputMessage,
	outputChannel chan<- *LearnerOutputMessage,
) {
	inputMessage := <-inputChannel
	outputChannel <- &LearnerOutputMessage{
		Objective: l.ComputeObjective(inputMessage.NewParams),
	}
}

func (l *Learner) ResetIterators() {
	for i, objective := range l.config.Objectives {
		dataIterator := NewDataIterator(
			objective,
			l.config.Streaming[i].DataStreamer,
		)
		l.stochadexConfig.Partitions[i].Iteration = dataIterator
		l.dataIterations[i] = dataIterator
	}
}

// NewLearner creates a new Learner struct given a config and loaded settings.
func NewLearner(
	config *LearnerConfig,
	settings *simulator.LoadSettingsConfig,
) *Learner {
	// handle some initial typing nonsense
	iterations := make([]simulator.Iteration, 0)
	dataIterations := make([]*DataIterator, 0)
	for i, objective := range config.Objectives {
		dataIterator := NewDataIterator(objective, config.Streaming[i].DataStreamer)
		iterations = append(iterations, dataIterator)
		dataIterations = append(dataIterations, dataIterator)
	}
	implementations := &simulator.LoadImplementationsConfig{
		Iterations:           iterations,
		OutputCondition:      &simulator.NilOutputCondition{},
		OutputFunction:       &simulator.NilOutputFunction{},
		TerminationCondition: config.Streaming[0].TerminationCondition,
		TimestepFunction:     config.Streaming[0].TimestepFunction,
	}
	stochadexConfig := simulator.NewStochadexConfig(
		settings,
		implementations,
	)
	return &Learner{
		config:          config,
		stochadexConfig: stochadexConfig,
		dataIterations:  dataIterations,
	}
}
