// Package ports provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package ports

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /trainings)
	GetTrainings(w http.ResponseWriter, r *http.Request)

	// (POST /trainings)
	CreateTraining(w http.ResponseWriter, r *http.Request)

	// (DELETE /trainings/{trainingUUID})
	CancelTraining(w http.ResponseWriter, r *http.Request, trainingUUID openapi_types.UUID)

	// (PUT /trainings/{trainingUUID}/approve-reschedule)
	ApproveRescheduleTraining(w http.ResponseWriter, r *http.Request, trainingUUID openapi_types.UUID)

	// (PUT /trainings/{trainingUUID}/reject-reschedule)
	RejectRescheduleTraining(w http.ResponseWriter, r *http.Request, trainingUUID openapi_types.UUID)

	// (PUT /trainings/{trainingUUID}/request-reschedule)
	RequestRescheduleTraining(w http.ResponseWriter, r *http.Request, trainingUUID openapi_types.UUID)

	// (PUT /trainings/{trainingUUID}/reschedule)
	RescheduleTraining(w http.ResponseWriter, r *http.Request, trainingUUID openapi_types.UUID)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetTrainings operation middleware
func (siw *ServerInterfaceWrapper) GetTrainings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTrainings(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateTraining operation middleware
func (siw *ServerInterfaceWrapper) CreateTraining(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateTraining(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CancelTraining operation middleware
func (siw *ServerInterfaceWrapper) CancelTraining(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "trainingUUID" -------------
	var trainingUUID openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "trainingUUID", runtime.ParamLocationPath, chi.URLParam(r, "trainingUUID"), &trainingUUID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "trainingUUID", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CancelTraining(w, r, trainingUUID)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// ApproveRescheduleTraining operation middleware
func (siw *ServerInterfaceWrapper) ApproveRescheduleTraining(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "trainingUUID" -------------
	var trainingUUID openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "trainingUUID", runtime.ParamLocationPath, chi.URLParam(r, "trainingUUID"), &trainingUUID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "trainingUUID", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ApproveRescheduleTraining(w, r, trainingUUID)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// RejectRescheduleTraining operation middleware
func (siw *ServerInterfaceWrapper) RejectRescheduleTraining(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "trainingUUID" -------------
	var trainingUUID openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "trainingUUID", runtime.ParamLocationPath, chi.URLParam(r, "trainingUUID"), &trainingUUID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "trainingUUID", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RejectRescheduleTraining(w, r, trainingUUID)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// RequestRescheduleTraining operation middleware
func (siw *ServerInterfaceWrapper) RequestRescheduleTraining(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "trainingUUID" -------------
	var trainingUUID openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "trainingUUID", runtime.ParamLocationPath, chi.URLParam(r, "trainingUUID"), &trainingUUID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "trainingUUID", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RequestRescheduleTraining(w, r, trainingUUID)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// RescheduleTraining operation middleware
func (siw *ServerInterfaceWrapper) RescheduleTraining(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "trainingUUID" -------------
	var trainingUUID openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "trainingUUID", runtime.ParamLocationPath, chi.URLParam(r, "trainingUUID"), &trainingUUID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "trainingUUID", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RescheduleTraining(w, r, trainingUUID)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshallingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshallingParamError) Error() string {
	return fmt.Sprintf("Error unmarshalling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshallingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/trainings", wrapper.GetTrainings)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/trainings", wrapper.CreateTraining)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/trainings/{trainingUUID}", wrapper.CancelTraining)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/trainings/{trainingUUID}/approve-reschedule", wrapper.ApproveRescheduleTraining)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/trainings/{trainingUUID}/reject-reschedule", wrapper.RejectRescheduleTraining)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/trainings/{trainingUUID}/request-reschedule", wrapper.RequestRescheduleTraining)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/trainings/{trainingUUID}/reschedule", wrapper.RescheduleTraining)
	})

	return r
}
