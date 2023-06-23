package learning

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// OptimisationAlgorithm
type OptimisationAlgorithm interface {
	GetNewParams(
		objectiveHistory []float64,
		paramHistory [][]*simulator.OtherParams,
	) [][]*simulator.OtherParams
	Terminate(
		objectiveHistories [][]float64,
		paramHistories [][][]*simulator.OtherParams,
	) bool
}

// LearningOptimiser
type LearningOptimiser struct {
	Learners             []Learner
	config               *OptimiserConfig
	objectiveHistories   [][]float64
	paramHistories       [][][]*simulator.OtherParams
	newWorkChannels      [](chan *LearnerInputMessage)
	workReturnedChannels [](chan *LearnerOutputMessage)
	numberOfLearners     int
}

func (l *LearningOptimiser) UpdateObjectiveHistories() {
	for index := 0; index < l.numberOfLearners; index++ {
		outputMessage := <-l.workReturnedChannels[index]
		l.objectiveHistories[index] = append(
			[]float64{outputMessage.Objective},
			l.objectiveHistories[index]...,
		)
		if len(l.objectiveHistories[index]) > l.config.HistoryDepth {
			l.objectiveHistories[index] =
				l.objectiveHistories[index][:l.config.HistoryDepth]
		}
	}
}

func (l *LearningOptimiser) RunNewBatch() {
	for index := 0; index < l.numberOfLearners; index++ {
		i := index
		go func() {
			l.Learners[i].ReceiveAndSendObjectives(
				l.newWorkChannels[i],
				l.workReturnedChannels[i],
			)
		}()
	}
	// send messages on the new work channels to ask for the next objective
	// values in the case of each learner
	for index := 0; index < l.numberOfLearners; index++ {
		l.newWorkChannels[index] <- &LearnerInputMessage{
			NewParams: l.paramHistories[index][0],
		}
	}
}

func (l *LearningOptimiser) AddNewParamsToHistory() {
	for index := 0; index < l.numberOfLearners; index++ {
		l.paramHistories[index] = append(
			l.config.Algorithm.GetNewParams(
				l.objectiveHistories[index],
				l.paramHistories[index],
			),
			l.paramHistories[index]...,
		)
		if len(l.paramHistories[index]) > l.config.HistoryDepth {
			l.paramHistories[index] =
				l.paramHistories[index][:l.config.HistoryDepth]
		}
	}
}

func (l *LearningOptimiser) Run() {

	for !l.config.Algorithm.Terminate(
		l.objectiveHistories,
		l.paramHistories,
	) {
		l.AddNewParamsToHistory()
		l.RunNewBatch()
		l.UpdateObjectiveHistories()
	}

}
