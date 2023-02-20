package training_test

import (
	"testing"
	"time"

	"github.com/dbaeka/workouts-go/internal/trainings/domain/training"
	"github.com/stretchr/testify/assert"
)

func TestTraining_RescheduleTraining(t *testing.T) {
	tr := newExampleTraining(t)

	oldTime := tr.Time()
	newTime := time.Now().AddDate(0, 0, 14).Round(time.Hour)

	// it's always a good idea to ensure about pre-conditions in the test ;-)
	assert.False(t, oldTime.Equal(newTime))

	err := tr.RescheduleTraining(newTime)
	assert.NoError(t, err)
	assert.True(t, tr.Time().Equal(newTime))
}

func TestTraining_RescheduleTraining_less_than_24h_before(t *testing.T) {
	originalTime := time.Now().Round(time.Hour)
	rescheduleRequestTime := originalTime.AddDate(0, 0, 5)

	tr := newExampleTrainingWithTime(t, originalTime)

	err := tr.RescheduleTraining(rescheduleRequestTime)

	assert.EqualError(t, err, training.CantRescheduleBeforeTimeError{
		TrainingTime: tr.Time(),
	}.Error())
}
