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

	registerHandler := httptransport.NewServer(
		ctx,
		e.RegisterEndpoint,
		decodeRegisterRequest,
		encodeResponse,
	)
	loginHandler := httptransport.NewServer(
		ctx,
		e.LoginEndpoint,
		decodeLoginRequest,
		encodeResponse,
	)
	resetPasswordHandler := httptransport.NewServer(
		ctx,
		e.ResetPasswordEndpoint,
		decodeResetPasswordRequest,
		encodeResponse,
	)
	changePasswordHandler := httptransport.NewServer(
		ctx,
		e.ChangePasswordEndpoint,
		decodeChangePasswordRequest,
		encodeResponse,
	)
	listHandler := httptransport.NewServer(
		ctx,
		e.ListEndpoint,
		decodeListRequest,
		encodeResponse,
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

func encodeResponse(ctx context.Context, w http.ResponseWriter, d interface{}) error {
	if e, ok := d.(errorer); ok && e.error() != nil {
		// Now its a business logic error
		// Extract base domain error
		encodeError(ctx, errors.Cause(e.error()), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(d)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
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
