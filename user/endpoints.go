package user

import (
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

// Endpoints combine all the user service endpoints under single type.
type Endpoints struct {
	RegisterEndpoint       endpoint.Endpoint
	LoginEndpoint          endpoint.Endpoint
	ResetPasswordEndpoint  endpoint.Endpoint
	ChangePasswordEndpoint endpoint.Endpoint
	ListEndpoint           endpoint.Endpoint
}

// MakeEndpoints returns Endpoints type which is the combination of
// all the user service endpoints.
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		RegisterEndpoint:       MakeRegisterEndpoint(s),
		LoginEndpoint:          MakeLoginEndpoint(s),
		ResetPasswordEndpoint:  MakeResetPasswordEndpoint(s),
		ChangePasswordEndpoint: MakeChangePasswordEndpoint(s),
		ListEndpoint:           MakeListEndpoint(s),
	}
}

func MakeRegisterEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(registerRequest)
		u, e := s.Register(ctx, req.NewUser)
		if e != nil {
			return registerResponse{User: nil, Error: e}, nil
		}
		return registerResponse{User: &u, Status: http.StatusCreated}, nil
	}
}

func MakeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loginRequest)
		u, e := s.Login(ctx, req.Email, req.Password)
		if e != nil {
			return loginResponse{User: nil, Error: e}, nil
		}
		return loginResponse{User: &u}, nil
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
			return resetPasswordResponse{Error: e}, nil
		}
		return resetPasswordResponse{Message: "password reset success"}, nil
	}
}

func MakeChangePasswordEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(changePasswordRequest)
		if req.Token == "" {
			return nil, ErrUnauthorized
		}
		u, e := s.AuthToken(ctx, req.Token)
		if e != nil {
			return nil, e
		}
		if req.NewPassword != req.ConfirmNewPassword {
			return nil, ErrPasswordMismatch
		}
		e = s.ChangePassword(ctx, u.ID, req.OldPassword, req.NewPassword)
		if e != nil {
			return changePasswordResponse{Error: e}, nil
		}
		return changePasswordResponse{Message: "change password success"}, nil

	}
}

func MakeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		var next, prev string
		users, total, e := s.List(ctx, req.Order, req.Limit, req.Offset)
		if e != nil {
			return listResponse{Error: e}, nil
		}
		limit, offset, err := nextLimitOffset(total, req.Limit, req.Offset)
		if err == nil {
			params := req.URL.Query()
			params = appendLimitOffset(params, limit, offset)
			next = req.URL.Host + req.URL.Path + "?" + params.Encode()
		}
		limit, offset, err = prevLimitOffset(total, req.Limit, req.Offset)
		if err == nil {
			params := req.URL.Query()
			params = appendLimitOffset(params, limit, offset)
			prev = req.URL.Host + req.URL.Path + "?" + params.Encode()
		}

		return listResponse{
			Users: users, Total: total,
			Prev: prev, Next: next,
		}, nil

	}
}

type registerRequest struct {
	NewUser
}

type registerResponse struct {
	Status int   `json:"-"`
	User   *User `json:"user,omitempty"`
	Error  error `json:"error,omitempty"`
}

func (r registerResponse) status() int {
	return r.Status
}

func (r registerResponse) error() error {
	return r.Error
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Status int   `json:"-"`
	User   *User `json:"user,omitempty"`
	Error  error `json:"error,omitempty"`
}

func (l loginResponse) status() int {
	return l.Status
}

func (r loginResponse) error() error {
	return r.Error
}

type resetPasswordRequest struct {
	Key                string `json:"key"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

type resetPasswordResponse struct {
	Status  int    `json:"-"`
	Message string `json:"message,omitempty"`
	Error   error  `json:"error,omitempty"`
}

func (r resetPasswordResponse) status() int {
	return r.Status
}

func (r resetPasswordResponse) error() error {
	return r.Error
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
	Status  int    `json:"-"`
	Message string `json:"message,omitempty"`
	Error   error  `json:"error,omitempty"`
}

func (c changePasswordResponse) status() int {
	return c.Status
}

func (r changePasswordResponse) error() error {
	return r.Error
}

type listRequest struct {
	Order  string `json:"order"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`

	URL *url.URL `json:"-"`
}

type listResponse struct {
	Status int    `json:"-"`
	Users  []User `json:"users"`
	Error  error  `json:"error,omitempty"`

	Total int    `json:"-"`
	Prev  string `json:"-"`
	Next  string `json:"-"`
}

func (r listResponse) status() int {
	return r.Status
}

func (r listResponse) error() error {
	return r.Error
}

func (r listResponse) page() (int, string, string) {
	return r.Total, r.Prev, r.Next
}
