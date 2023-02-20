package service

import (
	"context"
	"time"
)

type TrainerServiceMock struct {
}

func (t TrainerServiceMock) ScheduleTraining(_ context.Context, _ time.Time) error {
	return nil
}

func (t TrainerServiceMock) CancelTraining(_ context.Context, _ time.Time) error {
	return nil
}

func (t TrainerServiceMock) MoveTraining(_ context.Context, _ time.Time, _ time.Time) error {
	return nil
}

type UserServiceMock struct {
}

func (u UserServiceMock) UpdateTrainingBalance(_ context.Context, _ string, _ int) error {
	return nil
}
