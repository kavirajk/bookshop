package auth

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints store all the auth service endpoints
type Endpoints struct {
	LoginEndpoint endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		LoginEndpoint: makeLoginEndpoint(s),
	}
}

func makeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(loginRequest)
		bundle, err := s.Login(ctx, req.Email, req.Password)
		if err != nil {
			return loginResponse{Err: err}, nil
		}
		return loginResponse{AccessToken: bundle.AccessToken, RefreshToken: bundle.RefreshToken}, nil
	}
}

// loginRequest is the one client send via transport.
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginResponse is the returned reponse to client on login endpoint.
type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Err          error  `json:"err, omitempty"`
}
