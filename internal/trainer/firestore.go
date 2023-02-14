package main

import (
	"context"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

type db struct {
	firestoreClient *firestore.Client
}

func (d db) TrainerHoursCollection() *firestore.CollectionRef {
	return d.firestoreClient.Collection("trainer-hours")
}

func (d db) DocumentRef(dateTimeToUpdate time.Time) *firestore.DocumentRef {
	return d.TrainerHoursCollection().Doc(dateTimeToUpdate.Format("2006-01-02"))
}

func (d db) GetDates(ctx context.Context, params *GetTrainerAvailableHoursParams) ([]Date, error) {
	dates, err := d.QueryDates(params, ctx)
	if err != nil {
		return nil, err
	}
	dates = addMissingDates(params, dates)

	for _, date := range dates {
		sort.Slice(date.Hours, func(i, j int) bool { return date.Hours[i].Hour.Before(date.Hours[j].Hour) })
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Date.Before(dates[j].Date.Time) })

	return dates, nil
}

func (d db) QueryDates(params *GetTrainerAvailableHoursParams, ctx context.Context) ([]Date, error) {
	iter := d.
		TrainerHoursCollection().
		Where("Date.Time", ">=", params.DateFrom).
		Where("Date.Time", "<=", params.DateTo).
		Documents(ctx)

	var dates []Date

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		date := Date{}
		if err := doc.DataTo(&date); err != nil {
			return nil, err
		}
		date = setDefaultAvailability(date)
		dates = append(dates, date)
	}

	return dates, nil
}

const (
	week = time.Hour * 24 * 7

	weeksAllowedToSet = 6
)

func setAvailability(hourUpdate time.Time, date Date, availabilityToSet bool) (Date, error) {
	found := false

	for i := range date.Hours {
		if date.Hours[i].Hour.Equal(hourUpdate) {
			if date.Hours[i].HasTrainingScheduled && (availabilityToSet == false) {
				return Date{}, errors.New("training already scheduled, cannot make  unavailable")
			}

			date.Hours[i].Available = availabilityToSet
			found = true
		}
	}
	if !found {
		newHour := Hour{
			Available: availabilityToSet,
			Hour:      hourUpdate,
		}
		date.Hours = append(date.Hours, newHour)
	}

	return date, nil
}
