package main

import (
	"net/http"

	"github.com/dbaeka/workouts-go/internal/common/auth"
	"github.com/dbaeka/workouts-go/internal/common/server/httperr"
	"github.com/go-chi/render"
)

// HttpServer defines interface to satisfy OpenAPI generated code functions for this API
type HttpServer struct {
	db db
}

func (h HttpServer) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	authUser, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.Unauthorized("no-user-found", err, w, r)
		return
	}

	user, err := h.db.GetUser(r.Context(), authUser.UUID)
	if err != nil {
		httperr.InternalError("cannot-get-user", err, w, r)
	}

	user.Role = authUser.Role
	user.DisplayName = authUser.DisplayName

	render.Respond(w, r, user)
}
