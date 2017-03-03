package user

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func MakeHTTPHandler(ctx context.Context, s Service, logger log.Logger) http.Handler {
	e := MakeEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}
	registerHandler := httptransport.NewServer(
		ctx,
		e.RegisterEndpoint,
		decodeRegisterRequest,
		encodeResponse,
		options...,
	)
	loginHandler := httptransport.NewServer(
		ctx,
		e.LoginEndpoint,
		decodeLoginRequest,
		encodeResponse,
		options...,
	)
	resetPasswordHandler := httptransport.NewServer(
		ctx,
		e.ResetPasswordEndpoint,
		decodeResetPasswordRequest,
		encodeResponse,
		options...,
	)
	changePasswordHandler := httptransport.NewServer(
		ctx,
		e.ChangePasswordEndpoint,
		decodeChangePasswordRequest,
		encodeResponse,
		options...,
	)
	listHandler := httptransport.NewServer(
		ctx,
		e.ListEndpoint,
		decodeListRequest,
		encodeResponse,
		options...,
	)

	r := mux.NewRouter()

	r.Handle("/users/v1/register", registerHandler).Methods("POST")
	r.Handle("/users/v1/login", loginHandler).Methods("POST")
	r.Handle("/users/v1/reset-password", resetPasswordHandler).Methods("POST")
	r.Handle("/users/v1/change-password", changePasswordHandler).Methods("POST")
	r.Handle("/users/v1/list", listHandler).Methods("GET")

	return r
}
func decodeRegisterRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var r registerRequest
	err := json.NewDecoder(req.Body).Decode(&r)
	return r, err
}

func decodeLoginRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var r loginRequest
	err := json.NewDecoder(req.Body).Decode(&r)
	return r, err
}

func decodeResetPasswordRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var r resetPasswordRequest
	err := json.NewDecoder(req.Body).Decode(&r)
	return r, err
}

func decodeChangePasswordRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var r resetPasswordRequest
	err := json.NewDecoder(req.Body).Decode(&r)
	return r, err
}

func decodeListRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return listRequest{}, nil
}

type errorer interface {
	error() error
}

type statuser interface {
	status() int
}

// formatResponse is the uniform response format used throughout the users service,
// for every endpoint response.
type formatResponse struct {
	Data interface{}  `json:"data,omitempty"`
	Meta metaResponse `json:"meta"`
}

// metaResponse is part of response json that tells about basic meta information.
type metaResponse struct {
	Status   int    `json:"status"`
	Error    string `json:"error,omitempty"`
	Previous string `json:"previous,omitempty"`
	Next     string `json:"next,omitempty"`
	Total    int    `json:"total,omitempty"`
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, d interface{}) error {
	if e, ok := d.(errorer); ok && e.error() != nil {
		// Now its a business logic error.
		// Extract base domain error.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	status := http.StatusOK
	if s, ok := d.(statuser); ok && s.status() != 0 {
		status = s.status()
	}

	f := formatResponse{
		Data: d,
		Meta: metaResponse{Status: status},
	}
	return json.NewEncoder(w).Encode(f)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Its important to pass errors.Cause() as we decide status code based on
	// root error which is domain specific
	code := codeFrom(errors.Cause(err))
	w.WriteHeader(code)
	f := formatResponse{Meta: metaResponse{Status: code, Error: err.Error()}}
	json.NewEncoder(w).Encode(f)
}

func codeFrom(err error) int {
	switch err {
	case ErrUserNotFound:
		return http.StatusNotFound
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrInvalidPassword, ErrInvalidResetKey, ErrMissingField, ErrPasswordMismatch:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
