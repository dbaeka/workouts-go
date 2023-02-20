package service

import (
	"context"

	grpcClient "github.com/dbaeka/workouts-go/internal/common/client"
	"github.com/dbaeka/workouts-go/internal/trainings/adapters"
	"github.com/dbaeka/workouts-go/internal/trainings/app"
	"github.com/dbaeka/workouts-go/internal/trainings/app/command"
	"github.com/dbaeka/workouts-go/internal/trainings/app/query"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	trainerClient, closeTrainerClient, err := grpcClient.NewTrainerClient()
	if err != nil {
		panic(err)
	}

	usersClient, closeUsersClient, err := grpcClient.NewUsersClient()
	if err != nil {
		panic(err)
	}
	trainerGrpc := adapters.NewTrainerGrpc(trainerClient)
	usersGrpc := adapters.NewUsersGrpc(usersClient)

	cleanup := func() {
		_ = closeTrainerClient()
		_ = closeUsersClient()
	}

	return newApplication(ctx, trainerGrpc, usersGrpc), cleanup
}

func NewComponentTestApplication(ctx context.Context) app.Application {
	return newApplication(ctx, TrainerServiceMock{}, UserServiceMock{})
}

func newApplication(_ context.Context, trainerGrpc command.TrainerService, usersGrpc command.UserService) app.Application {
	mysqlDB, err := adapters.NewMySQLConnection()
	if err != nil {
		panic(err)
	}

	trainingsRepository := adapters.NewMySQLTrainingsRepository(mysqlDB)

	newApp := app.Application{
		Commands: app.Commands{
			ApproveTrainingReschedule: command.NewApproveTrainingRescheduleHandler(trainingsRepository, usersGrpc, trainerGrpc),
			CancelTraining:            command.NewCancelTrainingHandler(trainingsRepository, usersGrpc, trainerGrpc),
			RejectTrainingReschedule:  command.NewRejectTrainingRescheduleHandler(trainingsRepository),
			RescheduleTraining:        command.NewRescheduleTrainingHandler(trainingsRepository, usersGrpc, trainerGrpc),
			RequestTrainingReschedule: command.NewRequestTrainingRescheduleHandler(trainingsRepository),
			ScheduleTraining:          command.NewScheduleTrainingHandler(trainingsRepository, usersGrpc, trainerGrpc),
		},
		Queries: app.Queries{
			AllTrainings:     query.NewAllTrainingsHandler(trainingsRepository),
			TrainingsForUser: query.NewTrainingsForUserHandler(trainingsRepository),
		},
	}

	return newApp
}
