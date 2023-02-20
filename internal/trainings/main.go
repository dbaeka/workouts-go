package main

import (
	"context"
	"net/http"

	"github.com/dbaeka/workouts-go/internal/trainings/service"

	_ "github.com/dbaeka/workouts-go/internal/common/logs"
	"github.com/dbaeka/workouts-go/internal/common/server"
	"github.com/dbaeka/workouts-go/internal/trainings/ports"
	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()
	newApp, cleanup := service.NewApplication(ctx)
	defer cleanup()
	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(ports.NewHttpServer(newApp), router)
	})
}
