package query

import (
	"context"
	"github.com/dbaeka/workouts-go/internal/common/errors"
	"github.com/sirupsen/logrus"
	"time"
)

type AvailableHoursReadModel interface {
	AvailableHours(ctx context.Context, from time.Time, to time.Time) ([]Date, error)
}

type AvailableHoursHandler struct {
	readModel AvailableHoursReadModel
}

func NewAvailableHoursHandler(readModel AvailableHoursReadModel) AvailableHoursHandler {
	return AvailableHoursHandler{readModel: readModel}
}

type AvailableHours struct {
	From time.Time
	To   time.Time
}

func (h AvailableHoursHandler) Handle(ctx context.Context, query AvailableHours) (d []Date, err error) {
	start := time.Now()
	defer func() {
		logrus.
			WithError(err).
			WithField("duration", time.Since(start)).
			Debug("AvailableHoursHandler executed")
	}()

	if query.From.After(query.To) {
		return nil, errors.NewIncorrectInputError("Date from after date to", "date-from-after-date-to")
	}

	return h.readModel.AvailableHours(ctx, query.From, query.To)
}
