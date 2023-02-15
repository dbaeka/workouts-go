package main

import (
	"context"
	"github.com/dbaeka/workouts-go/internal/trainings/adapters"
	"github.com/dbaeka/workouts-go/internal/trainings/app"
	"github.com/dbaeka/workouts-go/internal/trainings/ports"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	grpcClient "github.com/dbaeka/workouts-go/internal/common/client"
	"github.com/dbaeka/workouts-go/internal/common/logs"
	"github.com/dbaeka/workouts-go/internal/common/server"
	"github.com/go-chi/chi/v5"
)

func main() {
	logs.Init()

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
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

	trainingsRepository := adapters.NewTrainingsFirestoreRepository(client)
	trainerGrpc := adapters.NewTrainerGrpc(trainerClient)
	usersGrpc := adapters.NewUsersGrpc(usersClient)

	trainingsService := app.NewTrainingsService(trainingsRepository, trainerGrpc, usersGrpc)

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(ports.NewHttpServer(trainingsService), router)
	})
}
