package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dbaeka/workouts-go/internal/common/genproto/users"
	_ "github.com/dbaeka/workouts-go/internal/common/logs"
	"github.com/dbaeka/workouts-go/internal/common/server"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

func main() {
	mysqlDB, err := NewMySQLConnection()
	if err != nil {
		panic(err)
	}
	mysqlRepo := db{mysqlDB}

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		go loadFixtures()
		server.RunHTTPServer(func(router chi.Router) http.Handler {
			return HandlerFromMux(HttpServer{mysqlRepo}, router) // Function from OpenAPI generated code
		})
	case "grpc":
		server.RunGRPCServer(func(server *grpc.Server) {
			svc := GrpcServer{db: mysqlRepo}
			users.RegisterUsersServiceServer(server, svc)
		})
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
