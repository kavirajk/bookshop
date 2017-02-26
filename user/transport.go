package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	jsonEncodingError = "json encoding error %v"
	jsonDecodingError = "json decoding error %v"
)

func decodeRegisterRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var r registerRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf(jsonDecodingError, err)
	}
	return r, nil
}

func encodeRegisterResponse(ctx context.Context, w http.ResponseWriter, d interface{}) error {
	if err := json.NewEncoder(w).Encode(d); err != nil {
		return fmt.Errorf(jsonEncodingError, err)
	}
	return nil
}

func decodeLoginRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var r loginRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf(jsonDecodingError, err)
	}
	return r, nil
}

func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, d interface{}) error {
	if err := json.NewEncoder(w).Encode(d); err != nil {
		return fmt.Errorf(jsonEncodingError, err)
	}
	return nil
}

func decodeResetPasswordRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var r resetPasswordRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf(jsonDecodingError, err)
	}
	return r, nil
}

func encodeResetPasswordResponse(ctx context.Context, w http.ResponseWriter, d interface{}) error {
	if err := json.NewEncoder(w).Encode(d); err != nil {
		return fmt.Errorf(jsonEncodingError, err)
	}
	return nil
}

func decodeChangePasswordRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var r resetPasswordRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf(jsonDecodingError, err)
	}
	return r, nil
}

func encodeChangePasswordResponse(ctx context.Context, w http.ResponseWriter, d interface{}) error {
	if err := json.NewEncoder(w).Encode(d); err != nil {
		return fmt.Errorf(jsonEncodingError, err)
	}
	return nil
}
