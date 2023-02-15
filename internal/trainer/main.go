package main

import (
	"context"
	"fmt"
	"github.com/dbaeka/workouts-go/internal/trainer/domain/hour"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/dbaeka/workouts-go/internal/common/genproto/trainer"
	_ "github.com/dbaeka/workouts-go/internal/common/logs"
	"github.com/dbaeka/workouts-go/internal/common/server"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}

	mySQLDB, err := NewMySQLConnection()
	if err != nil {
		panic(err)
	}

	hourFactory, err := hour.NewFactory(hour.FactoryConfig{
		MaxWeeksInTheFutureToSet: 6,
		MinUtcHour:               12,
		MaxUtcHour:               20,
	})
	if err != nil {
		panic(err)
	}

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		go loadFixtures(mySQLDB)

		server.RunHTTPServer(func(router chi.Router) http.Handler {
			return HandlerFromMux(HttpServer{firebaseDB, NewMySQLHourRepository(mySQLDB, hourFactory)}, router)
		})
	case "grpc":
		server.RunGRPCServer(func(server *grpc.Server) {
			svc := GrpcServer{hourRepository: NewMySQLHourRepository(mySQLDB, hourFactory)}
			trainer.RegisterTrainerServiceServer(server, svc)
		})
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
