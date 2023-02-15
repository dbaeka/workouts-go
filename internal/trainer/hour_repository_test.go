package main_test

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/dbaeka/workouts-go/internal/trainer"
	"github.com/dbaeka/workouts-go/internal/trainer/domain/hour"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// in general global state is not the best idea, but sometimes rules have some exceptions!
// in tests it's just simpler to re-use one instance of the factory
var testHourFactory = hour.MustNewFactory(hour.FactoryConfig{
	// 500 weeks gives us enough entropy to avoid duplicated dates
	// (even if duplicate dates should be not a problem)
	MaxWeeksInTheFutureToSet: 500,
	MinUtcHour:               0,
	MaxUtcHour:               24,
})

func newMySQLRepository(t *testing.T) *main.MySQLHourRepository {
	db, err := main.NewMySQLConnection()
	require.NoError(t, err)

	return main.NewMySQLHourRepository(db, testHourFactory)
}

func newValidAvailableHour(t *testing.T) *hour.Hour {
	hourTime := newValidHourTime()

	hourDto, err := testHourFactory.NewAvailableHour(hourTime)
	require.NoError(t, err)

	return hourDto
}

// usedHours is storing hours used during the test,
// to ensure that within one test run we are not using the same hour
// (it should be not a problem between test runs)
var usedHours = sync.Map{}

func newValidHourTime() time.Time {
	for {
		minTime := time.Now().AddDate(0, 0, 1)

		minTimestamp := minTime.Unix()
		maxTimestamp := minTime.AddDate(0, 0, testHourFactory.Config().MaxWeeksInTheFutureToSet*7).Unix()

		t := time.Unix(rand.Int63n(maxTimestamp-minTimestamp)+minTimestamp, 0).Truncate(time.Hour).Local()

		_, alreadyUsed := usedHours.LoadOrStore(t.Unix(), true)
		if !alreadyUsed {
			return t
		}
	}
}

func assertHourInRepository(ctx context.Context, t *testing.T, repo hour.Repository, hour *hour.Hour) {
	require.NotNil(t, hour)

	hourFromRepo, err := repo.GetOrCreateHour(ctx, hour.Time())
	require.NoError(t, err)

	assert.Equal(t, hour, hourFromRepo)
}

type Repository struct {
	Name       string
	Repository hour.Repository
}

func createRepositories(t *testing.T) []Repository {
	return []Repository{
		{
			Name:       "MySQL",
			Repository: newMySQLRepository(t),
		},
		{
			Name:       "memory",
			Repository: main.NewMemoryHourRepository(testHourFactory),
		},
	}
}

func testUpdateHour(t *testing.T, repository hour.Repository) {
	t.Helper()
	ctx := context.Background()

	testCases := []struct {
		Name       string
		CreateHour func(*testing.T) *hour.Hour
	}{
		{
			Name: "available_hour",
			CreateHour: func(t *testing.T) *hour.Hour {
				return newValidAvailableHour(t)
			},
		},
		{
			Name: "not_available_hour",
			CreateHour: func(t *testing.T) *hour.Hour {
				h := newValidAvailableHour(t)
				require.NoError(t, h.MakeNotAvailable())

				return h
			},
		},
		{
			Name: "hour_with_training",
			CreateHour: func(t *testing.T) *hour.Hour {
				h := newValidAvailableHour(t)
				require.NoError(t, h.ScheduleTraining())

				return h
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			newHour := tc.CreateHour(t)

			err := repository.UpdateHour(ctx, newHour.Time(), func(_ *hour.Hour) (*hour.Hour, error) {
				// UpdateHour provides us existing/new *hour.Hour,
				// but we are ignoring this hour and persisting result of `CreateHour`
				// we can assert this hour later in assertHourInRepository
				return newHour, nil
			})
			require.NoError(t, err)

			assertHourInRepository(ctx, t, repository, newHour)
		})
	}
}

func testUpdateHour_parallel(t *testing.T, repository hour.Repository) {
	t.Helper()
	ctx := context.Background()

	hourTime := newValidHourTime()

	// we are adding available hour
	err := repository.UpdateHour(ctx, hourTime, func(h *hour.Hour) (*hour.Hour, error) {
		if err := h.MakeAvailable(); err != nil {
			return nil, err
		}
		return h, nil
	})
	require.NoError(t, err)

	workersCount := 10
	workersDone := sync.WaitGroup{}
	workersDone.Add(workersCount)

	startWorkers := make(chan struct{})
	trainingsScheduled := make(chan int, workersCount)

	// we are trying to do race condition, in practice only one worker should be able to finish transaction
	for worker := 0; worker < workersCount; worker++ {
		workerNum := worker

		go func() {
			defer workersDone.Done()
			<-startWorkers

			schedulingTraining := false

			err := repository.UpdateHour(ctx, hourTime, func(h *hour.Hour) (*hour.Hour, error) {
				if h.HasTrainingScheduled() {
					return h, nil
				}
				if err := h.ScheduleTraining(); err != nil {
					return nil, err
				}

				schedulingTraining = true

				return h, nil
			})

			if schedulingTraining && err == nil {
				trainingsScheduled <- workerNum
			}
		}()
	}

	close(startWorkers)
	workersDone.Wait()
	close(trainingsScheduled)

	var workersScheduledTraining []int

	for workerNum := range trainingsScheduled {
		workersScheduledTraining = append(workersScheduledTraining, workerNum)
	}

	assert.Len(t, workersScheduledTraining, 1, "only one worker should schedule training")
}

func testUpdateHour_rollback(t *testing.T, repository hour.Repository) {
	t.Helper()
	ctx := context.Background()

	hourTime := newValidHourTime()

	err := repository.UpdateHour(ctx, hourTime, func(h *hour.Hour) (*hour.Hour, error) {
		require.NoError(t, h.MakeAvailable())
		return h, nil
	})

	err = repository.UpdateHour(ctx, hourTime, func(h *hour.Hour) (*hour.Hour, error) {
		assert.True(t, h.IsAvailable())
		require.NoError(t, h.MakeNotAvailable())

		return h, errors.New("something went wrong")
	})
	require.Error(t, err)

	persistedHour, err := repository.GetOrCreateHour(ctx, hourTime)
	require.NoError(t, err)

	assert.True(t, persistedHour.IsAvailable(), "availability change was persisted, not rolled back")
}

// testHourRepository_update_existing is testing path of creating a new hour and updating this hour.
func testHourRepository_update_existing(t *testing.T, repository hour.Repository) {
	t.Helper()
	ctx := context.Background()

	testHour := newValidAvailableHour(t)

	err := repository.UpdateHour(ctx, testHour.Time(), func(_ *hour.Hour) (*hour.Hour, error) {
		return testHour, nil
	})
	require.NoError(t, err)
	assertHourInRepository(ctx, t, repository, testHour)

	var expectedHour *hour.Hour
	err = repository.UpdateHour(ctx, testHour.Time(), func(h *hour.Hour) (*hour.Hour, error) {
		if err := h.ScheduleTraining(); err != nil {
			return nil, err
		}
		expectedHour = h
		return h, nil
	})
	require.NoError(t, err)

	assertHourInRepository(ctx, t, repository, expectedHour)
}

func TestRepository(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

	repositories := createRepositories(t)

	for i := range repositories {
		// When you are looping over slice and later using iterated value in goroutine (here because of t.Parallel()),
		// you need to always create variable scoped in loop body!
		// More info here: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		r := repositories[i]

		t.Run(r.Name, func(t *testing.T) {
			// It's always a good idea to build all non-unit tests to be able to work in parallel.
			// Thanks to that, your tests will be always fast, and you will not be afraid to add more tests because of slowdown.
			t.Parallel()

			t.Run("testUpdateHour", func(t *testing.T) {
				t.Parallel()
				testUpdateHour(t, r.Repository)
			})
			t.Run("testUpdateHour_parallel", func(t *testing.T) {
				t.Parallel()
				testUpdateHour_parallel(t, r.Repository)
			})
			t.Run("testHourRepository_update_existing", func(t *testing.T) {
				t.Parallel()
				testHourRepository_update_existing(t, r.Repository)
			})
			t.Run("testUpdateHour_rollback", func(t *testing.T) {
				t.Parallel()
				testUpdateHour_rollback(t, r.Repository)
			})
		})
	}
}
