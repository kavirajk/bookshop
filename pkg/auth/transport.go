package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// MakeHTTPHandler expose auth service over http transport.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := httprouter.New()
	e := MakeEndpoints(s)
	r.Handler("POST", "/api/v1/login/", httptransport.NewServer(
		e.LoginEndpoint,
		decodeLoginRequest,
		encodeLoginResponse,
	))
	return r
}

// decodeLoginRequest decode user login payload to `loginRequest`.
func decodeLoginRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	r := loginRequest{}
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		return r, err
	}
	return r, nil
}

// errorre interface is used to encode domain related errors.
type errorer interface {
	error() error
}

// encodeResponse is a generic helper to encode response to JSON.
func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, data interface{}) error {
	// Not transport error. business-logic error.
	if e, ok := data.(errorer); ok && e.error() != nil {
		encodeError(e.error(), w)
		return nil
	}
	w.Header().Set("ContentType", "application/json; charset:utf-8")
	return json.NewEncoder(w).Encode(&data)
}

// encodeError sets proper header for business-logic error and encodes
// the response to JSON.
func encodeError(err error, w http.ResponseWriter) {
	if err == nil {
		panic("called with non-nil error")
	}
	status := codeFromError(err)
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset:utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  fmt.Sprintf("%s", err),
		"status": status,
	})
}

// codeFromError returns matching http status code for specific (business-logic)errors.
func codeFromError(err error) int {
	switch errors.Cause(err) {
	case ErrTokenNotFound, ErrTTLNotFound:
		return http.StatusNotFound
	case ErrClaimIDInvalid, ErrClaimEmailInvalid, ErrInvalidCredentials:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
