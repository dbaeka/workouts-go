package main

import (
	"context"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/dbaeka/workouts-go/internal/common/auth"
	"github.com/dbaeka/workouts-go/internal/common/genproto/trainer"
	"github.com/dbaeka/workouts-go/internal/common/genproto/users"
	"github.com/dbaeka/workouts-go/internal/common/server/httperr"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type HttpServer struct {
	db            db
	trainerClient trainer.TrainerServiceClient
	usersClient   users.UsersServiceClient
}

func (h HttpServer) GetTrainings(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.Unauthorized("no-user-found", err, w, r)
		return
	}

	trainings, err := h.db.GetTrainings(r.Context(), user)
	if err != nil {
		httperr.InternalError("cannot-get-trainings", err, w, r)
		return
	}

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

	training := &Training{
		Notes:    postTraining.Notes,
		Time:     postTraining.Time,
		User:     user.DisplayName,
		UserUuid: uuid.MustParse(user.UUID),
		Uuid:     uuid.New(),
	}

	collection := h.db.TrainingsCollection()

	err = h.db.firestoreClient.RunTransaction(r.Context(), func(ctx context.Context, tx *firestore.Transaction) error {
		docs, err := tx.Documents(collection.Where("Time", "==", postTraining.Time)).GetAll()
		if err != nil {
			return errors.Wrap(err, "unable to get actual docs")
		}
		if len(docs) > 0 {
			return errors.Errorf("there is training already at %s", postTraining.Time)
		}

		_, err = h.usersClient.UpdateTrainingBalance(ctx, &users.UpdateTrainingBalanceRequest{
			UserId:       user.UUID,
			AmountChange: -1,
		})
		if err != nil {
			return errors.Wrap(err, "unable to change trainings balance")
		}

		timestamp := timestamppb.New(postTraining.Time)
		if timestamp == nil {
			return errors.Wrap(err, "unable to convert time to proto timestamp")
		}
		_, err = h.trainerClient.UpdateHour(ctx, &trainer.UpdateHourRequest{
			Time:                 timestamp,
			HasTrainingScheduled: true,
			Available:            false,
		})
		if err != nil {
			return errors.Wrap(err, "unable to update trainer hour")
		}

		return tx.Create(collection.Doc(training.Uuid.String()), training)
	})
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

	trainingsCollection := h.db.TrainingsCollection()

	err = h.db.firestoreClient.RunTransaction(r.Context(), func(ctx context.Context, tx *firestore.Transaction) error {
		trainingDocumentRef := trainingsCollection.Doc(trainingUUID.String())

		firestoreTraining, err := tx.Get(trainingDocumentRef)
		if err != nil {
			return errors.Wrap(err, "unable to get actual docs")
		}

		training := &Training{}
		err = firestoreTraining.DataTo(training)
		if err != nil {
			return errors.Wrap(err, "unable to load document")
		}

		if user.Role != "trainer" && training.UserUuid.String() != user.UUID {
			return errors.Errorf("user '%s' is trying to cancel training of user '%s'", user.UUID, training.UserUuid)
		}

		var trainingBalanceDelta int64
		if training.canBeCancelled() { // More than 24 hour cancel
			// just give training back
			trainingBalanceDelta = 1
		} else {
			if user.Role == "trainer" {
				// 1 for cancelled training +1 fine for cancelling by trainer less than 24h before training
				trainingBalanceDelta = 2
			} else {
				// fine for cancelling less than 24h before training
				trainingBalanceDelta = 0
			}
		}

		if trainingBalanceDelta != 0 {
			_, err := h.usersClient.UpdateTrainingBalance(ctx, &users.UpdateTrainingBalanceRequest{
				UserId:       training.UserUuid.String(),
				AmountChange: trainingBalanceDelta,
			})
			if err != nil {
				return errors.Wrap(err, "unable to change trainings balance")
			}
		}

		timestamp := timestamppb.New(training.Time)
		if timestamp == nil {
			return errors.Wrap(err, "unable to convert time to proto timestamp")
		}
		_, err = h.trainerClient.UpdateHour(ctx, &trainer.UpdateHourRequest{
			Time:                 timestamp,
			HasTrainingScheduled: false,
			Available:            true,
		})
		if err != nil {
			return errors.Wrap(err, "unable to update trainer hour")
		}

		return tx.Delete(trainingDocumentRef)
	})
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

	collection := h.db.TrainingsCollection()

	err = h.db.firestoreClient.RunTransaction(r.Context(), func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(h.db.TrainingsCollection().Doc(trainingUUID.String()))
		if err != nil {
			return errors.Wrap(err, "could not find training")
		}

		docs, err := tx.Documents(collection.Where("Time", "==", rescheduleTraining.Time)).GetAll()
		if err != nil {
			return errors.Wrap(err, "unable to get actual docs")
		}
		if len(docs) > 0 {
			return errors.Errorf("there is training already at %s", rescheduleTraining.Time)
		}

		var training Training
		err = doc.DataTo(&training)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal training")
		}

		if training.canBeCancelled() {
			err = h.rescheduleTraining(ctx, training.Time, rescheduleTraining.Time)
			if err != nil {
				return errors.Wrap(err, "unable to reschedule training")
			}

			training.Time = rescheduleTraining.Time
			training.Notes = rescheduleTraining.Notes
		} else {
			training.ProposedTime = &rescheduleTraining.Time
			training.MoveProposedBy = &user.Role
			training.Notes = rescheduleTraining.Notes
		}

		return tx.Set(collection.Doc(training.Uuid.String()), training)
	})
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

	err = h.db.firestoreClient.RunTransaction(r.Context(), func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(h.db.TrainingsCollection().Doc(trainingUUID.String()))
		if err != nil {
			return errors.Wrap(err, "could not find training")
		}

		var training Training
		err = doc.DataTo(&training)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal training")
		}

		if training.ProposedTime == nil {
			return errors.New("training has no proposed time")
		}
		if training.MoveProposedBy == nil {
			return errors.New("training has no MoveProposedBy")
		}
		if *training.MoveProposedBy == "trainer" && training.UserUuid.String() != user.UUID {
			return errors.Errorf("user '%s' cannot approve reschedule of user '%s'", user.UUID, training.UserUuid)
		}
		if *training.MoveProposedBy == user.Role {
			return errors.New("reschedule cannot be accepted by requesting person")
		}

		training.Time = *training.ProposedTime
		training.ProposedTime = nil

		return tx.Set(h.db.TrainingsCollection().Doc(training.Uuid.String()), training)
	})
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

	err = h.db.firestoreClient.RunTransaction(r.Context(), func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(h.db.TrainingsCollection().Doc(trainingUUID.String()))
		if err != nil {
			return errors.Wrap(err, "could not find training")
		}

		var training Training
		err = doc.DataTo(&training)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal training")
		}

		if training.MoveProposedBy == nil {
			return errors.New("training has no MoveProposedBy")
		}
		if *training.MoveProposedBy != "trainer" && training.UserUuid.String() != user.UUID {
			return errors.Errorf("user '%s' cannot approve reschedule of user '%s'", user.UUID, training.UserUuid)
		}

		training.ProposedTime = nil

		return tx.Set(h.db.TrainingsCollection().Doc(training.Uuid.String()), training)
	})
	if err != nil {
		httperr.InternalError("cannot-update-training", err, w, r)
		return
	}
}

func (h HttpServer) rescheduleTraining(ctx context.Context, oldTime, newTime time.Time) error {
	oldTimeProto := timestamppb.New(oldTime)
	if oldTimeProto == nil {
		return errors.New("unable to convert time to proto timestamp")
	}

	newTimeProto := timestamppb.New(newTime)
	if newTimeProto == nil {
		return errors.New("unable to convert time to proto timestamp")
	}

	_, err := h.trainerClient.UpdateHour(ctx, &trainer.UpdateHourRequest{
		Time:                 newTimeProto,
		HasTrainingScheduled: true,
		Available:            false,
	})
	if err != nil {
		return errors.Wrap(err, "unable to update trainer hour")
	}

	_, err = h.trainerClient.UpdateHour(ctx, &trainer.UpdateHourRequest{
		Time:                 oldTimeProto,
		HasTrainingScheduled: false,
		Available:            true,
	})
	if err != nil {
		return errors.Wrap(err, "unable to update trainer hour")
	}

	return nil
}
