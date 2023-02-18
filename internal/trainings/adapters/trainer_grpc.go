package adapters

import (
	"context"
	"time"

	"github.com/dbaeka/workouts-go/internal/common/genproto/trainer"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TrainerGrpc struct {
	client trainer.TrainerServiceClient
}

func NewTrainerGrpc(client trainer.TrainerServiceClient) TrainerGrpc {
	return TrainerGrpc{client: client}
}

func (s TrainerGrpc) ScheduleTraining(ctx context.Context, trainingTime time.Time) error {
	timestamp := timestamppb.New(trainingTime)
	if timestamp == nil {
		return errors.New("unable to convert time to proto timestamp")
	}

	_, err := s.client.ScheduleTraining(ctx, &trainer.UpdateHourRequest{
		Time: timestamp,
	})

	return err
}

func (s TrainerGrpc) CancelTraining(ctx context.Context, trainingTime time.Time) error {
	timestamp := timestamppb.New(trainingTime)
	if timestamp == nil {
		return errors.New("unable to convert time to proto timestamp")
	}

	_, err := s.client.CancelTraining(ctx, &trainer.UpdateHourRequest{
		Time: timestamp,
	})

	return err
}

func (s TrainerGrpc) MoveTraining(
	ctx context.Context,
	newTime time.Time,
	originalTrainingTime time.Time,
) error {
	err := s.ScheduleTraining(ctx, newTime)
	if err != nil {
		return errors.Wrap(err, "unable to schedule training")
	}

	err = s.CancelTraining(ctx, originalTrainingTime)
	if err != nil {
		return errors.Wrap(err, "unable to cancel training")
	}

	return nil
}
