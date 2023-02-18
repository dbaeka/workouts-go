package command

import (
	"context"
	"github.com/dbaeka/workouts-go/internal/common/errors"
	"github.com/dbaeka/workouts-go/internal/trainer/domain/hour"
	"time"
)

type ScheduleTrainingHandler struct {
	hourRepo hour.Repository
}

func NewScheduleTrainingHandler(hourRepo hour.Repository) ScheduleTrainingHandler {
	if hourRepo == nil {
		panic("nil hourRepo")
	}

	return ScheduleTrainingHandler{hourRepo: hourRepo}
}

func (h ScheduleTrainingHandler) Handle(ctx context.Context, hourToSchedule time.Time) error {
	if err := h.hourRepo.UpdateHour(ctx, hourToSchedule, func(hr *hour.Hour) (*hour.Hour, error) {
		if err := hr.ScheduleTraining(); err != nil {
			return nil, err
		}
		return hr, nil
	}); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-update-availability")
	}

	return nil
}
