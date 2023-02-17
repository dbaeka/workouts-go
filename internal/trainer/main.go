package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dbaeka/workouts-go/internal/common/genproto/trainer"
	_ "github.com/dbaeka/workouts-go/internal/common/logs"
	"github.com/dbaeka/workouts-go/internal/common/server"
	"github.com/dbaeka/workouts-go/internal/trainer/adapters"
	"github.com/dbaeka/workouts-go/internal/trainer/app"
	"github.com/dbaeka/workouts-go/internal/trainer/domain/hour"
	"github.com/dbaeka/workouts-go/internal/trainer/ports"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

func main() {
	mySQLDB, err := adapters.NewMySQLConnection()
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

	datesRepository := adapters.NewMySQLDatesRepository(mySQLDB, hourFactory)
	if datesRepository == nil {
		panic(err)
	}

	hourRepository := adapters.NewMySQLHourRepository(mySQLDB, hourFactory)
	service := app.NewHourService(datesRepository, hourRepository)

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		go loadFixtures(datesRepository)

		server.RunHTTPServer(func(router chi.Router) http.Handler {
			return ports.HandlerFromMux(ports.NewHttpServer(service), router)
		})
	case "grpc":
		server.RunGRPCServer(func(server *grpc.Server) {
			svc := ports.NewGrpcServer(hourRepository)
			trainer.RegisterTrainerServiceServer(server, svc)
		})
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
