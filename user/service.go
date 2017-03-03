package user

import (
	"errors"

	"golang.org/x/net/context"
)

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidResetKey = errors.New("invalid resetkey")
	ErrUserNotFound    = errors.New("user not found")
)

// Service defines all the services provided user package.
type Service interface {
	Register(ctx context.Context, user NewUser) (User, error)
	Login(ctx context.Context, email, password string) (User, error)

	// Used to authenticate via token
	AuthToken(ctx context.Context, token string) (User, error)

	// Used to change user's password without old password (e.g: Forget Password)
	ResetPassword(ctx context.Context, key, newpass string) error

	// Used to change user's password with old password (e.g: Profile settings)
	ChangePassword(ctx context.Context, userID string, oldpass, newpass string) error

	// List available user based on limit and offset.
	// order takes string in the format "username asc" or " username desc"
	// or in combination of multiple fields like "username asc, email desc"
	List(ctx context.Context, order string, limit, offset int) (users []User, total int, err error)
}

// service is a simple implementation of Service interface.
type service struct {
	repo Repo
}

// NewService takes User Repo and returns new User Service.
func NewService(repo Repo) Service {
	return service{repo: repo}
}

// Register registers the new user.
// in case of non-nil error return User is always empty
func (s service) Register(_ context.Context, nuser NewUser) (User, error) {
	if err := nuser.Validate(); err != nil {
		return User{}, err
	}
	user := nuser.User()
	if err := s.repo.Create(&user); err != nil {
		return User{}, err
	}
	return user, nil
}

// Login is used to authenticate any user with email and password.
func (s service) Login(_ context.Context, email, password string) (User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return User{}, err
	}
	if user.Password != calculatePassHash(password, user.Salt) {
		return User{}, ErrUnauthorized
	}
	return user, nil
}

// AuthToken is used to get user associated with token.
func (s service) AuthToken(_ context.Context, token string) (User, error) {
	user, err := s.repo.GetByToken(token)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// ResetPassword is used to change the users' password with key and newPass.
// Typical use-case would be forgot password.
func (s service) ResetPassword(ctx context.Context, key, newPass string) error {
	user, err := s.repo.GetByResetKey(key)
	if err != nil {
		return err
	}
	if user.ResetKey != key {
		return ErrInvalidResetKey
	}
	return s.changePassword(ctx, user, newPass)
}

// ChangePassword is used to change the user's password with oldpassword.
// Typical use-case would be to use it in profile page
func (s service) ChangePassword(ctx context.Context, userID, oldPass, newPass string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}
	if user.Password != calculatePassHash(oldPass, user.Salt) {
		return ErrInvalidPassword
	}
	return s.changePassword(ctx, user, newPass)
}

// ListUser lists all the available users in the system.
func (s service) List(ctx context.Context, order string, limit, offset int) ([]User, int, error) {
	return s.repo.List(order, limit, offset)
}

// changePassword is an unexpoted helper function to change the password of the user.
func (s service) changePassword(_ context.Context, user User, newPass string) error {
	user.Password = calculatePassHash(newPass, user.Salt)
	if err := s.repo.Save(&user); err != nil {
		return err
	}
	return nil
}

// Middleware is a Service middleware for user Service
type Middleware func(Service) Service
