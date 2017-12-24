package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeHTTPHandler expose auth service over http transport.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeEndpoints(s)
	r.Path("/login/").Methods("POST").Handler(httptransport.NewServer(
		e.LoginEndpoint,
		decodeLoginRequest,
		encodeLoginResponse,
	))
	r.Path("/hello/").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello")
	})
	return r
}

// decodeLoginRequest decode user login payload to `loginRequest`.
func decodeLoginRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	fmt.Println("decode login request")
	res := loginResponse{}
	if err := json.NewDecoder(req.Body).Decode(&res); err != nil {
		return res, err
	}
	return res, nil
}

// encodeLoginResponse encode service response to `loginResponse`.
func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, data interface{}) error {
	if err := json.NewEncoder(w).Encode(&data); err != nil {
		return err
	}
	return nil
}
