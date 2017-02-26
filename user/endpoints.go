package user

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

func MakeRegisterEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(registerRequest)
		u, e := s.Register(ctx, req.NewUser)
		return registerResponse{User: u, Error: e}, nil
	}
}

func MakeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loginRequest)
		u, e := s.Login(ctx, req.Email, req.Password)
		return loginResponse{User: u, Error: e}, nil
	}
}

func MakeResetPasswordEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(resetPasswordRequest)
		if req.NewPassword != req.ConfirmNewPassword {
			return nil, ErrPasswordMismatch
		}
		e := s.ResetPassword(ctx, req.Key, req.NewPassword)
		if e != nil {
			return resetPasswordResponse{Message: "password reset success"}, nil
		}
		return resetPasswordResponse{Error: e}, nil
	}
}

func MakeChangePasswordEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(changePasswordRequest)
		if req.Token == "" {
			return nil, ErrUnauthorized
		}
		u, e := s.AuthToken(req.Token)
		if e != nil {
			return nil, e
		}
		if req.NewPassword != req.ConfirmNewPassword {
			return nil, ErrPasswordMismatch
		}
		e = s.ChangePassword(ctx, u.ID, req.OldPassword, req.NewPassword)
		if e != nil {
			return changePasswordResponse{Message: "change password success"}, nil
		}
		return changePasswordResponse{Error: e}, nil
	}
}

type registerRequest struct {
	NewUser
}

type registerResponse struct {
	User
	Error error `json:"error"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	User
	Error error `json:"error"`
}

type resetPasswordRequest struct {
	Key                string `json:"key"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

type resetPasswordResponse struct {
	Message string `json:"message,omitempty"`
	Error   error  `json:"error"`
}

type authTokenRequest struct {
	Token string `json:"-"` // We get from header
}

type changePasswordRequest struct {
	UserID             string `json:"-"`
	Token              string `json:"-"`
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

type changePasswordResponse struct {
	Message string `json:"message,omitempty"`
	Error   error  `json:"error"`
}
