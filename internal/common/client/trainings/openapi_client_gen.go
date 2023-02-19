// Package trainings provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package trainings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetTrainings request
	GetTrainings(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateTraining request with any body
	CreateTrainingWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateTraining(ctx context.Context, body CreateTrainingJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CancelTraining request
	CancelTraining(ctx context.Context, trainingUUID openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ApproveRescheduleTraining request
	ApproveRescheduleTraining(ctx context.Context, trainingUUID openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RejectRescheduleTraining request
	RejectRescheduleTraining(ctx context.Context, trainingUUID openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RequestRescheduleTraining request with any body
	RequestRescheduleTrainingWithBody(ctx context.Context, trainingUUID openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	RequestRescheduleTraining(ctx context.Context, trainingUUID openapi_types.UUID, body RequestRescheduleTrainingJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RescheduleTraining request with any body
	RescheduleTrainingWithBody(ctx context.Context, trainingUUID openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	RescheduleTraining(ctx context.Context, trainingUUID openapi_types.UUID, body RescheduleTrainingJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetTrainings(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetTrainingsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateTrainingWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateTrainingRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateTraining(ctx context.Context, body CreateTrainingJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateTrainingRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CancelTraining(ctx context.Context, trainingUUID openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCancelTrainingRequest(c.Server, trainingUUID)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ApproveRescheduleTraining(ctx context.Context, trainingUUID openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewApproveRescheduleTrainingRequest(c.Server, trainingUUID)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RejectRescheduleTraining(ctx context.Context, trainingUUID openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRejectRescheduleTrainingRequest(c.Server, trainingUUID)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RequestRescheduleTrainingWithBody(ctx context.Context, trainingUUID openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRequestRescheduleTrainingRequestWithBody(c.Server, trainingUUID, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RequestRescheduleTraining(ctx context.Context, trainingUUID openapi_types.UUID, body RequestRescheduleTrainingJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRequestRescheduleTrainingRequest(c.Server, trainingUUID, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RescheduleTrainingWithBody(ctx context.Context, trainingUUID openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRescheduleTrainingRequestWithBody(c.Server, trainingUUID, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RescheduleTraining(ctx context.Context, trainingUUID openapi_types.UUID, body RescheduleTrainingJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRescheduleTrainingRequest(c.Server, trainingUUID, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetTrainingsRequest generates requests for GetTrainings
func NewGetTrainingsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/trainings")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateTrainingRequest calls the generic CreateTraining builder with application/json body
func NewCreateTrainingRequest(server string, body CreateTrainingJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateTrainingRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateTrainingRequestWithBody generates requests for CreateTraining with any type of body
func NewCreateTrainingRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/trainings")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewCancelTrainingRequest generates requests for CancelTraining
func NewCancelTrainingRequest(server string, trainingUUID openapi_types.UUID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "trainingUUID", runtime.ParamLocationPath, trainingUUID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/trainings/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewApproveRescheduleTrainingRequest generates requests for ApproveRescheduleTraining
func NewApproveRescheduleTrainingRequest(server string, trainingUUID openapi_types.UUID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "trainingUUID", runtime.ParamLocationPath, trainingUUID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/trainings/%s/approve-reschedule", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewRejectRescheduleTrainingRequest generates requests for RejectRescheduleTraining
func NewRejectRescheduleTrainingRequest(server string, trainingUUID openapi_types.UUID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "trainingUUID", runtime.ParamLocationPath, trainingUUID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/trainings/%s/reject-reschedule", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewRequestRescheduleTrainingRequest calls the generic RequestRescheduleTraining builder with application/json body
func NewRequestRescheduleTrainingRequest(server string, trainingUUID openapi_types.UUID, body RequestRescheduleTrainingJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewRequestRescheduleTrainingRequestWithBody(server, trainingUUID, "application/json", bodyReader)
}

// NewRequestRescheduleTrainingRequestWithBody generates requests for RequestRescheduleTraining with any type of body
func NewRequestRescheduleTrainingRequestWithBody(server string, trainingUUID openapi_types.UUID, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "trainingUUID", runtime.ParamLocationPath, trainingUUID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/trainings/%s/request-reschedule", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewRescheduleTrainingRequest calls the generic RescheduleTraining builder with application/json body
func NewRescheduleTrainingRequest(server string, trainingUUID openapi_types.UUID, body RescheduleTrainingJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewRescheduleTrainingRequestWithBody(server, trainingUUID, "application/json", bodyReader)
}

// NewRescheduleTrainingRequestWithBody generates requests for RescheduleTraining with any type of body
func NewRescheduleTrainingRequestWithBody(server string, trainingUUID openapi_types.UUID, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "trainingUUID", runtime.ParamLocationPath, trainingUUID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/trainings/%s/reschedule", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetTrainings request
	GetTrainingsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetTrainingsResponse, error)

	// CreateTraining request with any body
	CreateTrainingWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTrainingResponse, error)

	CreateTrainingWithResponse(ctx context.Context, body CreateTrainingJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTrainingResponse, error)

	// CancelTraining request
	CancelTrainingWithResponse(ctx context.Context, trainingUUID openapi_types.UUID, reqEditors ...RequestEditorFn) (*CancelTrainingResponse, error)

	// ApproveRescheduleTraining request
	ApproveRescheduleTrainingWithResponse(ctx context.Context, trainingUUID openapi_types.UUID, reqEditors ...RequestEditorFn) (*ApproveRescheduleTrainingResponse, error)

	// RejectRescheduleTraining request
	RejectRescheduleTrainingWithResponse(ctx context.Context, trainingUUID openapi_types.UUID, reqEditors ...RequestEditorFn) (*RejectRescheduleTrainingResponse, error)

	// RequestRescheduleTraining request with any body
	RequestRescheduleTrainingWithBodyWithResponse(ctx context.Context, trainingUUID openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RequestRescheduleTrainingResponse, error)

	RequestRescheduleTrainingWithResponse(ctx context.Context, trainingUUID openapi_types.UUID, body RequestRescheduleTrainingJSONRequestBody, reqEditors ...RequestEditorFn) (*RequestRescheduleTrainingResponse, error)

	// RescheduleTraining request with any body
	RescheduleTrainingWithBodyWithResponse(ctx context.Context, trainingUUID openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RescheduleTrainingResponse, error)

	RescheduleTrainingWithResponse(ctx context.Context, trainingUUID openapi_types.UUID, body RescheduleTrainingJSONRequestBody, reqEditors ...RequestEditorFn) (*RescheduleTrainingResponse, error)
}

type GetTrainingsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Trainings
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetTrainingsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetTrainingsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateTrainingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r CreateTrainingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateTrainingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CancelTrainingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r CancelTrainingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CancelTrainingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ApproveRescheduleTrainingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r ApproveRescheduleTrainingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ApproveRescheduleTrainingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RejectRescheduleTrainingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r RejectRescheduleTrainingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RejectRescheduleTrainingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RequestRescheduleTrainingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r RequestRescheduleTrainingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RequestRescheduleTrainingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RescheduleTrainingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r RescheduleTrainingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RescheduleTrainingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetTrainingsWithResponse request returning *GetTrainingsResponse
func (c *ClientWithResponses) GetTrainingsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetTrainingsResponse, error) {
	rsp, err := c.GetTrainings(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTrainingsResponse(rsp)
}

// CreateTrainingWithBodyWithResponse request with arbitrary body returning *CreateTrainingResponse
func (c *ClientWithResponses) CreateTrainingWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTrainingResponse, error) {
	rsp, err := c.CreateTrainingWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTrainingResponse(rsp)
}

func (c *ClientWithResponses) CreateTrainingWithResponse(ctx context.Context, body CreateTrainingJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTrainingResponse, error) {
	rsp, err := c.CreateTraining(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTrainingResponse(rsp)
}

// CancelTrainingWithResponse request returning *CancelTrainingResponse
func (c *ClientWithResponses) CancelTrainingWithResponse(ctx context.Context, trainingUUID openapi_types.UUID, reqEditors ...RequestEditorFn) (*CancelTrainingResponse, error) {
	rsp, err := c.CancelTraining(ctx, trainingUUID, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCancelTrainingResponse(rsp)
}

// ApproveRescheduleTrainingWithResponse request returning *ApproveRescheduleTrainingResponse
func (c *ClientWithResponses) ApproveRescheduleTrainingWithResponse(ctx context.Context, trainingUUID openapi_types.UUID, reqEditors ...RequestEditorFn) (*ApproveRescheduleTrainingResponse, error) {
	rsp, err := c.ApproveRescheduleTraining(ctx, trainingUUID, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseApproveRescheduleTrainingResponse(rsp)
}

// RejectRescheduleTrainingWithResponse request returning *RejectRescheduleTrainingResponse
func (c *ClientWithResponses) RejectRescheduleTrainingWithResponse(ctx context.Context, trainingUUID openapi_types.UUID, reqEditors ...RequestEditorFn) (*RejectRescheduleTrainingResponse, error) {
	rsp, err := c.RejectRescheduleTraining(ctx, trainingUUID, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRejectRescheduleTrainingResponse(rsp)
}

// RequestRescheduleTrainingWithBodyWithResponse request with arbitrary body returning *RequestRescheduleTrainingResponse
func (c *ClientWithResponses) RequestRescheduleTrainingWithBodyWithResponse(ctx context.Context, trainingUUID openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RequestRescheduleTrainingResponse, error) {
	rsp, err := c.RequestRescheduleTrainingWithBody(ctx, trainingUUID, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRequestRescheduleTrainingResponse(rsp)
}

func (c *ClientWithResponses) RequestRescheduleTrainingWithResponse(ctx context.Context, trainingUUID openapi_types.UUID, body RequestRescheduleTrainingJSONRequestBody, reqEditors ...RequestEditorFn) (*RequestRescheduleTrainingResponse, error) {
	rsp, err := c.RequestRescheduleTraining(ctx, trainingUUID, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRequestRescheduleTrainingResponse(rsp)
}

// RescheduleTrainingWithBodyWithResponse request with arbitrary body returning *RescheduleTrainingResponse
func (c *ClientWithResponses) RescheduleTrainingWithBodyWithResponse(ctx context.Context, trainingUUID openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RescheduleTrainingResponse, error) {
	rsp, err := c.RescheduleTrainingWithBody(ctx, trainingUUID, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRescheduleTrainingResponse(rsp)
}

func (c *ClientWithResponses) RescheduleTrainingWithResponse(ctx context.Context, trainingUUID openapi_types.UUID, body RescheduleTrainingJSONRequestBody, reqEditors ...RequestEditorFn) (*RescheduleTrainingResponse, error) {
	rsp, err := c.RescheduleTraining(ctx, trainingUUID, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRescheduleTrainingResponse(rsp)
}

// ParseGetTrainingsResponse parses an HTTP response from a GetTrainingsWithResponse call
func ParseGetTrainingsResponse(rsp *http.Response) (*GetTrainingsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetTrainingsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Trainings
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseCreateTrainingResponse parses an HTTP response from a CreateTrainingWithResponse call
func ParseCreateTrainingResponse(rsp *http.Response) (*CreateTrainingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateTrainingResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseCancelTrainingResponse parses an HTTP response from a CancelTrainingWithResponse call
func ParseCancelTrainingResponse(rsp *http.Response) (*CancelTrainingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CancelTrainingResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseApproveRescheduleTrainingResponse parses an HTTP response from a ApproveRescheduleTrainingWithResponse call
func ParseApproveRescheduleTrainingResponse(rsp *http.Response) (*ApproveRescheduleTrainingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ApproveRescheduleTrainingResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseRejectRescheduleTrainingResponse parses an HTTP response from a RejectRescheduleTrainingWithResponse call
func ParseRejectRescheduleTrainingResponse(rsp *http.Response) (*RejectRescheduleTrainingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RejectRescheduleTrainingResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseRequestRescheduleTrainingResponse parses an HTTP response from a RequestRescheduleTrainingWithResponse call
func ParseRequestRescheduleTrainingResponse(rsp *http.Response) (*RequestRescheduleTrainingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RequestRescheduleTrainingResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseRescheduleTrainingResponse parses an HTTP response from a RescheduleTrainingWithResponse call
func ParseRescheduleTrainingResponse(rsp *http.Response) (*RescheduleTrainingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RescheduleTrainingResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}
