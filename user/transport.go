package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	jsonEncodingError = "json encoding error %v"
	jsonDecodingError = "json decoding error %v"
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
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf(jsonDecodingError, err)
	}
	return r, nil
}

func decodeLoginRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var r loginRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf(jsonDecodingError, err)
	}
	return r, nil
}

func decodeResetPasswordRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var r resetPasswordRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf(jsonDecodingError, err)
	}
	return r, nil
}

func decodeChangePasswordRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var r resetPasswordRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf(jsonDecodingError, err)
	}
	return r, nil
}

func decodeListRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return listRequest{}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, d interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(d); err != nil {
		return fmt.Errorf(jsonEncodingError, err)
	}
	return nil
}
