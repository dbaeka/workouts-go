package query

import "context"

type AllTrainingsReadModel interface {
	AllTrainings(ctx context.Context) ([]Training, error)
}

type AllTrainingsHandler struct {
	readModel AllTrainingsReadModel
}

func NewAllTrainingsHandler(readModel AllTrainingsReadModel) AllTrainingsHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return AllTrainingsHandler{readModel: readModel}
}

func (h AllTrainingsHandler) Handle(ctx context.Context) (tr []Training, err error) {
	return h.readModel.AllTrainings(ctx)
}
