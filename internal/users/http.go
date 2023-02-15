package main

import (
	"net"
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
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		err = h.db.UpdateLastIP(r.Context(), authUser.UUID, host)
		if err != nil {
			httperr.InternalError("internal-server-error", err, w, r)
			return
		}
	}

	user, err := h.db.GetUser(r.Context(), authUser.UUID)
	if err != nil {
		httperr.InternalError("cannot-get-user", err, w, r)
	}

	userResponse := User{
		DisplayName: authUser.DisplayName,
		Balance:     user.Balance,
		Role:        authUser.Role,
	}

	render.Respond(w, r, userResponse)
}
