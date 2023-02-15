package httperr

import (
	"github.com/dbaeka/workouts-go/internal/common/errors"
	"net/http"

	"github.com/dbaeka/workouts-go/internal/common/logs"
	"github.com/go-chi/render"
)

func InternalError(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Internal server error", http.StatusInternalServerError)
}

func Unauthorized(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Unauthorized", http.StatusUnauthorized)
}

func BadRequest(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Bad request", http.StatusBadRequest)
}

func httpRespondWithError(err error, slug string, w http.ResponseWriter, r *http.Request, logMsg string, status int) {
	logs.GetLogEntry(r).WithError(err).WithField("error-slug", slug).Warn(logMsg)
	resp := ErrorResponse{slug, status}

	if err := render.Render(w, r, resp); err != nil {
		panic(err)
	}
}

func RespondWithSlugError(err error, w http.ResponseWriter, r *http.Request) {
	slugError, ok := err.(errors.SlugError)
	if !ok {
		InternalError("internal-server-error", err, w, r)
		return
	}

	switch slugError.ErrorType() {
	case errors.ErrorTypeAuthorization:
		Unauthorized(slugError.Slug(), slugError, w, r)
	case errors.ErrorTypeIncorrectInput:
		BadRequest(slugError.Slug(), slugError, w, r)
	default:
		InternalError(slugError.Slug(), slugError, w, r)
	}
}

type ErrorResponse struct {
	Slug       string `json:"slug"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
