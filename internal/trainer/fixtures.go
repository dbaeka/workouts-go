package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/dbaeka/workouts-go/internal/common/client"
	"github.com/dbaeka/workouts-go/internal/common/genproto/trainer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func loadFixtures(db db) {
	start := time.Now()
	ctx := context.Background()

	logrus.Debug("Waiting for trainer service")
	working := client.WaitForTrainerService(time.Second * 30)
	if !working {
		logrus.Error("Trainer gRPC service is not up")
		return
	}

	logrus.WithField("after", time.Now().Sub(start)).Debug("Trainer service is available")

	var canLoad bool
	var err error

	for {
		canLoad, err = canLoadFixtures(ctx, db)
		if err == nil {
			break
		}
		logrus.WithError(err).Error("Cannot check if fixtures can be loaded")
		time.Sleep(10 * time.Second)
	}

	if !canLoad {
		logrus.Debug("Trainer fixtures are already loaded")
		return
	}

	for {
		err = loadTrainerFixtures(ctx)
		if err == nil {
			break
		}

		logrus.WithError(err).Error("Cannot load trainer fixtures")
		time.Sleep(10 * time.Second)
	}

	logrus.WithField("after", time.Now().Sub(start)).Debug("Trainer fixtures loaded")
}

const daysToSet = 30

func loadTrainerFixtures(ctx context.Context) error {
	trainerClient, closeTrainerClient, err := client.NewTrainerClient()
	if err != nil {
		return err
	}
	defer func() { _ = closeTrainerClient() }()

	maxDate := time.Now().Add(time.Hour * 24 * daysToSet)
	localRand := rand.New(rand.NewSource(3))

	for date := time.Now(); date.Before(maxDate); date = date.Add(time.Hour * 24) {
		for hour := 12; hour <= 20; hour++ {
			trainingTime := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.UTC)

			ts := timestamppb.New(trainingTime)
			if ts == nil {
				return errors.Wrapf(err, "unable to marshal time %s", trainingTime)
			}

			if localRand.NormFloat64() > 0 {
				_, err = trainerClient.UpdateHour(ctx, &trainer.UpdateHourRequest{
					Time:                 ts,
					HasTrainingScheduled: false,
					Available:            true,
				})
				if err != nil {
					return errors.Wrap(err, "unable to update hour")
				}
			}
		}
	}

	return nil
}

func canLoadFixtures(ctx context.Context, db db) (bool, error) {
	documents, err := db.TrainerHoursCollection().Limit(daysToSet).Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}

	return len(documents) < daysToSet, nil
}
