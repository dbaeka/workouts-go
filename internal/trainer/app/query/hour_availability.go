package query

import (
	"context"
	"github.com/dbaeka/workouts-go/internal/trainer/domain/hour"
	"time"
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
