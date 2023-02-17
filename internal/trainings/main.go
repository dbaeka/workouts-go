package main

import (
	"net/http"

	grpcClient "github.com/dbaeka/workouts-go/internal/common/client"
	_ "github.com/dbaeka/workouts-go/internal/common/logs"
	"github.com/dbaeka/workouts-go/internal/common/server"
	"github.com/dbaeka/workouts-go/internal/trainings/adapters"
	"github.com/dbaeka/workouts-go/internal/trainings/app"
	"github.com/dbaeka/workouts-go/internal/trainings/ports"
	"github.com/go-chi/chi/v5"
)

func main() {

	mySQLDB, err := adapters.NewMySQLConnection()
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}

	trainerClient, closeTrainerClient, err := grpcClient.NewTrainerClient()
	if err != nil {
		panic(err)
	}
	defer func() { _ = closeTrainerClient() }()

	usersClient, closeUsersClient, err := grpcClient.NewUsersClient()
	if err != nil {
		panic(err)
	}
	defer func() { _ = closeUsersClient() }()

	trainingsRepository := adapters.NewMySQLTrainingsRepository(mySQLDB)
	trainerGrpc := adapters.NewTrainerGrpc(trainerClient)
	usersGrpc := adapters.NewUsersGrpc(usersClient)

	trainingsService := app.NewTrainingsService(trainingsRepository, trainerGrpc, usersGrpc)

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(ports.NewHttpServer(trainingsService), router)
	})
}
