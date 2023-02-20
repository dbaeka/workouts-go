package service

import (
	"context"

	"github.com/dbaeka/workouts-go/internal/trainer/adapters"
	"github.com/dbaeka/workouts-go/internal/trainer/app"
	"github.com/dbaeka/workouts-go/internal/trainer/app/command"
	"github.com/dbaeka/workouts-go/internal/trainer/app/query"
	"github.com/dbaeka/workouts-go/internal/trainer/domain/hour"
)

func NewApplication(_ context.Context) app.Application {
	mysqlDB, err := adapters.NewMySQLConnection()
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

	datesRepository := adapters.NewMySQLDatesRepository(mysqlDB, hourFactory)
	if datesRepository == nil {
		panic(err)
	}

	hourRepository := adapters.NewMySQLHourRepository(mysqlDB, hourFactory)

	return app.Application{
		Commands: app.Commands{
			CancelTraining:       command.NewCancelTrainingHandler(hourRepository),
			ScheduleTraining:     command.NewScheduleTrainingHandler(hourRepository),
			MakeHoursAvailable:   command.NewMakeHoursAvailableHandler(hourRepository),
			MakeHoursUnavailable: command.NewMakeHoursUnavailableHandler(hourRepository),
		},
		Queries: app.Queries{
			HourAvailability:      query.NewHourAvailabilityHandler(hourRepository),
			TrainerAvailableHours: query.NewAvailableHoursHandler(datesRepository),
		},
	}
}
