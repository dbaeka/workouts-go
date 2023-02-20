package command

import (
	"context"
	"time"

	"github.com/dbaeka/workouts-go/internal/common/errors"
	"github.com/dbaeka/workouts-go/internal/trainer/domain/hour"
)

type MakeHoursAvailableHandler struct {
	hourRepo hour.Repository
}

func NewMakeHoursAvailableHandler(hourRepo hour.Repository) MakeHoursAvailableHandler {
	if hourRepo == nil {
		panic("hourRepo is nil")
	}

	return MakeHoursAvailableHandler{hourRepo: hourRepo}
}

func (c MakeHoursAvailableHandler) Handle(ctx context.Context, hours []time.Time) error {
	for _, hourToUpdate := range hours {
		if err := c.hourRepo.UpdateHour(ctx, hourToUpdate, func(hr *hour.Hour) (*hour.Hour, error) {
			if err := hr.MakeAvailable(); err != nil {
				return nil, err
			}
			return hr, nil
		}); err != nil {
			return errors.NewSlugError(err.Error(), "unable-to-update-availability")
		}
	}

	return nil
}
