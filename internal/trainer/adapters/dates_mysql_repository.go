package adapters

import (
	"context"
	"database/sql"
	"time"

	"github.com/dbaeka/workouts-go/internal/trainer/domain/hour"
	"github.com/pkg/errors"

	"github.com/dbaeka/workouts-go/internal/trainer/app"
	"github.com/jmoiron/sqlx"
)

type MySQLDatesRepository struct {
	db          *sqlx.DB
	hourFactory hour.Factory
}

func NewMySQLDatesRepository(db *sqlx.DB, hourFactory hour.Factory) *MySQLDatesRepository {
	if db == nil {
		panic("missing db")
	}
	if hourFactory.IsZero() {
		panic("missing hourFactory")
	}
	return &MySQLDatesRepository{db: db}
}

func (m MySQLDatesRepository) GetDates(ctx context.Context, from time.Time, to time.Time) ([]app.Date, error) {
	return m.getDates(ctx, m.db, from, to)
}

func (m MySQLDatesRepository) getDates(
	ctx context.Context,
	db sqlContextGetter,
	from time.Time,
	to time.Time,
) ([]app.Date, error) {
	var dbHours []mysqlHour

	query := "SELECT * FROM `hours` WHERE `hour` BETWEEN ? AND ?"

	err := db.GetContext(ctx, &dbHours, query, from, to)
	if errors.Is(err, sql.ErrNoRows) {
		// in reality this date exists, even if it's not persisted, addMissingDates will handle this
		return []app.Date{}, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "unable to get hour from db")
	}

	var dates []app.Date

	for _, dbHour := range dbHours {
		availability, err := hour.NewAvailabilityFromString(dbHour.Availability)
		if err != nil {
			return nil, err
		}

		domainHour, err := m.hourFactory.UnmarshalHourFromDatabase(dbHour.Hour.Local(), availability)
		domainDate, err := app.DateFromHourDomain(*domainHour)
		if err != nil {
			return nil, err
		}
		dates = append(dates, *domainDate)
	}

	return dates, nil
}

func (m MySQLDatesRepository) CanLoadFixtures(ctx context.Context, daysToSet int) (bool, error) {
	var dbHours []mysqlHour

	query := "SELECT * FROM `hours` LIMIT ?"

	err := m.db.GetContext(ctx, &dbHours, query, daysToSet)
	if errors.Is(err, sql.ErrNoRows) {
		return true, nil
	} else if err != nil {
		return false, errors.Wrap(err, "unable to get hour from db")
	}

	return len(dbHours) < daysToSet, nil
}
