package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/dbaeka/workouts-go/internal/common/genproto/users"
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
	firebaseDB := db{firestoreClient}

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		go loadFixtures()
		server.RunHTTPServer(func(router chi.Router) http.Handler {
			return HandlerFromMux(HttpServer{firebaseDB}, router) // Function from OpenAPI generated code
		})
	case "grpc":
		server.RunGRPCServer(func(server *grpc.Server) {
			svc := GrpcServer{db: firebaseDB}
			users.RegisterUsersServiceServer(server, svc)
		})
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
