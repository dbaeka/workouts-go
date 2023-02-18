package main

import (
	"context"
	"net/http"

	"github.com/dbaeka/workouts-go/internal/trainings/app/command"
	"github.com/dbaeka/workouts-go/internal/trainings/app/query"

	grpcClient "github.com/dbaeka/workouts-go/internal/common/client"
	_ "github.com/dbaeka/workouts-go/internal/common/logs"
	"github.com/dbaeka/workouts-go/internal/common/server"
	"github.com/dbaeka/workouts-go/internal/trainings/adapters"
	"github.com/dbaeka/workouts-go/internal/trainings/app"
	"github.com/dbaeka/workouts-go/internal/trainings/ports"
	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()
	newApp, cleanup := newApplication(ctx)
	defer cleanup()
	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(ports.NewHttpServer(newApp), router)
	})
}

func newApplication(_ context.Context) (app.Application, func()) {
	mySQLDB, err := adapters.NewMySQLConnection()
	if err != nil {
		panic(err)
	}

	trainerClient, closeTrainerClient, err := grpcClient.NewTrainerClient()
	if err != nil {
		panic(err)
	}
	usersClient, closeUsersClient, err := grpcClient.NewUsersClient()
	if err != nil {
		panic(err)
	}

	trainingsRepository := adapters.NewMySQLTrainingsRepository(mySQLDB)
	trainerGrpc := adapters.NewTrainerGrpc(trainerClient)
	usersGrpc := adapters.NewUsersGrpc(usersClient)

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

	cleanup := func() {
		_ = closeTrainerClient()
		_ = closeUsersClient
	}

	return newApp, cleanup
}
