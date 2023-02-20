package app

import (
	"github.com/dbaeka/workouts-go/internal/trainings/app/command"
	"github.com/dbaeka/workouts-go/internal/trainings/app/query"
)

type Commands struct {
	ApproveTrainingReschedule command.ApproveTrainingRescheduleHandler
	CancelTraining            command.CancelTrainingHandler
	RejectTrainingReschedule  command.RejectTrainingRescheduleHandler
	RescheduleTraining        command.RescheduleTrainingHandler
	RequestTrainingReschedule command.RequestTrainingRescheduleHandler
	ScheduleTraining          command.ScheduleTrainingHandler
}

type Queries struct {
	AllTrainings     query.AllTrainingsHandler
	TrainingsForUser query.TrainingsForUserHandler
}

type Application struct {
	Commands Commands
	Queries  Queries
}
