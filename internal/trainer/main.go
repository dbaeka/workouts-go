package main

import (
	"context"
	"fmt"
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
	firebaseClient, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}

	firebaseDB := db{firebaseClient}

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		go loadFixtures(firebaseDB)

		server.RunHTTPServer(func(router chi.Router) http.Handler {
			return HandlerFromMux(HttpServer{firebaseDB}, router)
		})
	case "grpc":
		server.RunGRPCServer(func(server *grpc.Server) {
			svc := GrpcServer{db: firebaseDB}
			trainer.RegisterTrainerServiceServer(server, svc)
		})
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
