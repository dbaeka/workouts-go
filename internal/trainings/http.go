package main

import (
	"github.com/dbaeka/workouts-go/internal/common/auth"
	"github.com/dbaeka/workouts-go/internal/common/server/httperr"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"net/http"
)

type HttpServer struct {
	db db
}

func trainingModelsToResponse(models []TrainingModel) []Training {
	var trainings []Training
	for _, tm := range models {
		t := Training{
			CanBeCancelled:     tm.canBeCancelled(),
			MoveProposedBy:     tm.MoveProposedBy,
			MoveRequiresAccept: !tm.canBeCancelled(),
			Notes:              tm.Notes,
			ProposedTime:       tm.ProposedTime,
			Time:               tm.Time,
			User:               tm.User,
			UserUuid:           uuid.MustParse(tm.UserUUID),
			Uuid:               uuid.MustParse(tm.UUID),
		}

		trainings = append(trainings, t)
	}

	return trainings
}

func (h HttpServer) GetTrainings(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.Unauthorized("no-user-found", err, w, r)
		return
	}

	trainingModels, err := h.db.GetTrainings(r.Context(), user)
	if err != nil {
		httperr.InternalError("cannot-get-trainings", err, w, r)
		return
	}

	trainings := trainingModelsToResponse(trainingModels)
	trainingsResp := Trainings{trainings}

	render.Respond(w, r, trainingsResp)
}

func (h HttpServer) CreateTraining(w http.ResponseWriter, r *http.Request) {
	postTraining := PostTraining{}
	if err := render.Decode(r, &postTraining); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	// sanity check
	if len(postTraining.Notes) > 1000 {
		httperr.BadRequest("note-too-big", nil, w, r)
		return
	}

	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.Unauthorized("no-user-found", err, w, r)
		return
	}
	if user.Role != "attendee" {
		httperr.Unauthorized("invalid-role", nil, w, r)
		return
	}

	training := TrainingModel{
		Notes:    postTraining.Notes,
		Time:     postTraining.Time,
		User:     user.DisplayName,
		UserUUID: user.UUID,
		UUID:     uuid.New().String(),
	}

	err = h.db.CreateTraining(r.Context(), user, training)
	if err != nil {
		httperr.InternalError("cannot-create-training", err, w, r)
		return
	}
}

func (h HttpServer) CancelTraining(w http.ResponseWriter, r *http.Request, trainingUUID openapi_types.UUID) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.Unauthorized("no-user-found", err, w, r)
		return
	}

	err = h.db.CancelTraining(r.Context(), user, trainingUUID.String())
	if err != nil {
		httperr.InternalError("cannot-update-training", err, w, r)
		return
	}
}

func (h HttpServer) RescheduleTraining(w http.ResponseWriter, r *http.Request, trainingUUID openapi_types.UUID) {
	rescheduleTraining := PostTraining{}
	if err := render.Decode(r, &rescheduleTraining); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	// sanity check
	if len(rescheduleTraining.Notes) > 1000 {
		httperr.BadRequest("note-too-big", nil, w, r)
		return
	}

	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.Unauthorized("no-user-found", err, w, r)
		return
	}

	err = h.db.RescheduleTraining(r.Context(), user, trainingUUID.String(), rescheduleTraining.Time, rescheduleTraining.Notes)
	if err != nil {
		httperr.InternalError("cannot-update-training", err, w, r)
		return
	}
}

func (h HttpServer) ApproveRescheduleTraining(w http.ResponseWriter, r *http.Request, trainingUUID openapi_types.UUID) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.Unauthorized("no-user-found", err, w, r)
		return
	}
	err = h.db.ApproveTrainingReschedule(r.Context(), user, trainingUUID.String())
	if err != nil {
		httperr.InternalError("cannot-update-training", err, w, r)
		return
	}
}

func (h HttpServer) RejectRescheduleTraining(w http.ResponseWriter, r *http.Request, trainingUUID openapi_types.UUID) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.Unauthorized("no-user-found", err, w, r)
		return
	}

	err = h.db.RejectTrainingReschedule(r.Context(), user, trainingUUID.String())
	if err != nil {
		httperr.InternalError("cannot-update-training", err, w, r)
		return
	}
}
