package app

import (
	"time"

	"github.com/dbaeka/workouts-go/internal/trainer/domain/hour"
)

type Hour struct {
	Available            bool
	HasTrainingScheduled bool
	Hour                 time.Time
}

type Date struct {
	Date         time.Time
	HasFreeHours bool
	Hours        []Hour
}

func (d Date) FindHourInDate(timeToCheck time.Time) (*Hour, bool) {
	for i, hourDomain := range d.Hours {
		if hourDomain.Hour == timeToCheck {
			return &d.Hours[i], true
		}
	}

	return nil, false
}

type AvailableHoursRequest struct {
	DateFrom time.Time
	DateTo   time.Time
}

const (
	minHour = 12
	maxHour = 20
)

// setDefaultAvailability adds missing hours to Date model if they were not set
func setDefaultAvailability(date Date) Date {
HoursLoop:
	for hourDomain := minHour; hourDomain <= maxHour; hourDomain++ {
		hourDomain := time.Date(date.Date.Year(), date.Date.Month(), date.Date.Day(), hourDomain, 0, 0, 0, time.UTC)

		for i := range date.Hours {
			if date.Hours[i].Hour.Equal(hourDomain) {
				continue HoursLoop
			}
		}
		newHour := Hour{
			Available: false,
			Hour:      hourDomain,
		}

		date.Hours = append(date.Hours, newHour)
	}

	return date
}

func addMissingDates(dates []Date, from time.Time, to time.Time) []Date {
	for day := from.UTC(); day.Before(to) || day.Equal(to); day = day.Add(time.Hour * 24) {
		found := false
		for _, date := range dates {
			if date.Date.Equal(day) {
				found = true
				break
			}
		}

		if !found {
			date := Date{
				Date: day,
			}
			date = setDefaultAvailability(date)
			dates = append(dates, date)
		}
	}

	return dates
}

// DateFromHourDomain converts Date from Hour Domain.
func DateFromHourDomain(hour hour.Hour) (*Date, error) {
	hourDomain := Hour{
		Available:            hour.IsAvailable(),
		HasTrainingScheduled: hour.HasTrainingScheduled(),
		Hour:                 hour.Time(),
	}

	date := Date{
		Date:  hourDomain.Hour.Truncate(time.Hour * 24),
		Hours: []Hour{hourDomain},
	}

	return &date, nil
}
