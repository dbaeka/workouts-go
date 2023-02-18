package adapters

import (
	"context"
	"database/sql"
	"github.com/dbaeka/workouts-go/internal/trainer/app/query"
	"sort"
	"time"

	"github.com/dbaeka/workouts-go/internal/trainer/domain/hour"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type mysqlDate struct {
	Date         time.Time   `db:"date"`
	HasFreeHours bool        `db:"has_free_hours"`
	Hours        []mysqlHour `db:"hours"`
}

type MySQLDatesRepository struct {
	db      *sqlx.DB
	factory hour.Factory
}

func NewMySQLDatesRepository(db *sqlx.DB, factory hour.Factory) *MySQLDatesRepository {
	if db == nil {
		panic("missing db")
	}
	return &MySQLDatesRepository{db: db, factory: factory}
}

func (m MySQLDatesRepository) dateModelToApp(dm mysqlDate) query.Date {
	var hours []query.Hour
	for _, h := range dm.Hours {
		availability, err := hour.NewAvailabilityFromString(h.Availability)
		if err != nil {
			// skipping rows with invalid enum but unlikely to happen based on db design
			continue
		}
		domainHour, err := m.factory.UnmarshalHourFromDatabase(h.Hour.Local(), availability)
		if err != nil {
			// skipping rows with invalid enum but unlikely to happen based on db design
			continue
		}
		hours = append(hours, query.Hour{
			Available:            domainHour.IsAvailable(),
			HasTrainingScheduled: domainHour.HasTrainingScheduled(),
			Hour:                 h.Hour,
		})
	}

	return query.Date{
		Date:         dm.Date,
		HasFreeHours: dm.HasFreeHours,
		Hours:        hours,
	}
}

func addMissingDates(dates []query.Date, from time.Time, to time.Time) []query.Date {
	for day := from.UTC(); day.Before(to) || day.Equal(to); day = day.AddDate(0, 0, 1) {
		found := false
		for _, date := range dates {
			if date.Date.Equal(day) {
				found = true
				break
			}
		}

		if !found {
			date := query.Date{
				Date: day,
			}
			dates = append(dates, date)
		}
	}

	return dates
}

// setDefaultAvailability adds missing hours to Date model if they were not set
func (m MySQLDatesRepository) setDefaultAvailability(date query.Date) query.Date {
HoursLoop:
	for h := m.factory.Config().MinUtcHour; h <= m.factory.Config().MaxUtcHour; h++ {
		hr := time.Date(date.Date.Year(), date.Date.Month(), date.Date.Day(), h, 0, 0, 0, time.UTC)

		for i := range date.Hours {
			if date.Hours[i].Hour.Equal(hr) {
				continue HoursLoop
			}
		}
		newHour := query.Hour{
			Available: false,
			Hour:      hr,
		}

		date.Hours = append(date.Hours, newHour)
	}

	return date
}

func (m MySQLDatesRepository) AvailableHours(ctx context.Context, from time.Time, to time.Time) ([]query.Date, error) {
	return m.getDates(ctx, m.db, from, to)
}

func (m MySQLDatesRepository) getDates(
	ctx context.Context,
	db sqlContextGetter,
	from time.Time,
	to time.Time,
) ([]query.Date, error) {
	var dbDates []mysqlDate

	sqlQuery := "SELECT dates.*, hours.* FROM `dates` JOIN hours ON dates.date = DATE(hours.hour) WHERE dates.date BETWEEN ? AND ?"

	err := db.GetContext(ctx, &dbDates, sqlQuery, from, to)
	if errors.Is(err, sql.ErrNoRows) {
		// in reality this date exists, even if it's not persisted, addMissingDates will handle this
		return []query.Date{}, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "unable to get hour from db")
	}

	var dates []query.Date

	for _, dbDate := range dbDates {
		dates = append(dates, m.dateModelToApp(dbDate))
	}

	dates = addMissingDates(dates, from, to)
	for i, date := range dates {
		date = m.setDefaultAvailability(date)
		sort.Slice(date.Hours, func(i, j int) bool { return date.Hours[i].Hour.Before(date.Hours[j].Hour) })
		dates[i] = date
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Date.Before(dates[j].Date) })

	return dates, nil
}
