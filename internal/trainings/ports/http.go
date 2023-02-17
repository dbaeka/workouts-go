package ports

import (
	"net/http"

	"github.com/dbaeka/workouts-go/internal/common/auth"
	"github.com/dbaeka/workouts-go/internal/common/server/httperr"
	"github.com/dbaeka/workouts-go/internal/trainings/app"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type HttpServer struct {
	service app.TrainingService
}

func NewHttpServer(service app.TrainingService) HttpServer {
	return HttpServer{
		service: service,
	}
}

func appTrainingsToResponse(appTrainings []app.Training) []Training {
	var trainings []Training
	for _, tm := range appTrainings {
		t := Training{
			CanBeCancelled:     tm.CanBeCancelled(),
			MoveProposedBy:     tm.MoveProposedBy,
			MoveRequiresAccept: tm.MoveRequiresAccept(),
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
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	var appTrainings []app.Training
	if user.Role == "trainer" {
		appTrainings, err = h.service.GetAllTrainings(r.Context())
	} else {
		appTrainings, err = h.service.GetTrainingsForUser(r.Context(), user)
	}

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	trainings := appTrainingsToResponse(appTrainings)
	trainingsResp := Trainings{trainings}

	render.Respond(w, r, trainingsResp)
}

func (h HttpServer) CreateTraining(w http.ResponseWriter, r *http.Request) {
	postTraining := PostTraining{}
	if err := render.Decode(r, &postTraining); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	if user.Role != "attendee" {
		httperr.Unauthorized("invalid-role", nil, w, r)
		return
	}

	err = h.service.CreateTraining(r.Context(), user, postTraining.Time, postTraining.Notes)

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func (h HttpServer) CancelTraining(w http.ResponseWriter, r *http.Request, trainingUUID openapi_types.UUID) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.service.CancelTraining(r.Context(), user, trainingUUID.String())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func (h HttpServer) RescheduleTraining(w http.ResponseWriter, r *http.Request, trainingUUID openapi_types.UUID) {
	rescheduleTraining := PostTraining{}
	if err := render.Decode(r, &rescheduleTraining); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.service.RescheduleTraining(r.Context(), user, trainingUUID.String(), rescheduleTraining.Time, rescheduleTraining.Notes)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func (h HttpServer) ApproveRescheduleTraining(w http.ResponseWriter, r *http.Request, trainingUUID openapi_types.UUID) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	err = h.service.ApproveTrainingReschedule(r.Context(), user, trainingUUID.String())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func (h HttpServer) RejectRescheduleTraining(w http.ResponseWriter, r *http.Request, trainingUUID openapi_types.UUID) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.service.RejectTrainingReschedule(r.Context(), user, trainingUUID.String())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}
