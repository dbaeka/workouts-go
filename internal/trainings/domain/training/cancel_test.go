package training_test

import (
	"testing"
	"time"

	"github.com/dbaeka/workouts-go/internal/trainings/domain/training"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTraining_Cancel(t *testing.T) {
	tr := newExampleTraining(t)
	// it's always a good idea to ensure about pre-conditions in the test ;-)
	assert.False(t, tr.IsCanceled())

	err := tr.Cancel()
	require.NoError(t, err)
	assert.True(t, tr.IsCanceled())
}

func TestTraining_Cancel_already_canceled(t *testing.T) {
	tr := newCanceledTraining(t)

	assert.EqualError(t, tr.Cancel(), training.ErrTrainingAlreadyCanceled.Error())
}

func TestTraining_MoreThanDayUntilTraining(t *testing.T) {
	trainingNow := newExampleTrainingWithTime(t, time.Now())
	assert.False(t, trainingNow.CanBeCanceledForFree())

	trainingInTwoDays := newExampleTrainingWithTime(t, time.Now().AddDate(0, 0, 2))
	assert.True(t, trainingInTwoDays.CanBeCanceledForFree())
}
