package query

import (
	"context"
	"time"

	"github.com/dbaeka/workouts-go/internal/trainer/domain/hour"
)

type HourAvailabilityHandler struct {
	hourRepo hour.Repository
}

func NewHourAvailabilityHandler(hourRepo hour.Repository) HourAvailabilityHandler {
	if hourRepo == nil {
		panic("nil hourRepo")
	}

	return HourAvailabilityHandler{hourRepo: hourRepo}
}

func (h HourAvailabilityHandler) Handle(ctx context.Context, time time.Time) (bool, error) {
	hr, err := h.hourRepo.GetHour(ctx, time)
	if err != nil {
		return false, err
	}

	return hr.IsAvailable(), nil
}
