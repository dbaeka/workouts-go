package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dbaeka/workouts-go/internal/trainer/service"

	"github.com/dbaeka/workouts-go/internal/common/genproto/trainer"
	_ "github.com/dbaeka/workouts-go/internal/common/logs"
	"github.com/dbaeka/workouts-go/internal/common/server"
	"github.com/dbaeka/workouts-go/internal/trainer/ports"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	newApp := service.NewApplication(ctx)

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		go loadFixtures(newApp)

		server.RunHTTPServer(func(router chi.Router) http.Handler {
			return ports.HandlerFromMux(ports.NewHttpServer(newApp), router)
		})
	case "grpc":
		server.RunGRPCServer(func(server *grpc.Server) {
			svc := ports.NewGrpcServer(newApp)
			trainer.RegisterTrainerServiceServer(server, svc)
		})
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
