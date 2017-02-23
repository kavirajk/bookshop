package user

import (
	"errors"
	"fmt"

	"github.com/kavirajk/bookshop/user/email"
)

var (
	ErrUnauthorized = errors.New("Unauthorized")
)

type Service interface {
	Register(user NewUser) (User, error)
	Login(email, password string) (User, error)
	EmailResetPasswordLink(email string) error
	ResetPassword(key, oldpass, newpass string) error
	ChangePassword(oldpass, newpass string) error
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return service{repo: repo}
}

func (s service) Register(nuser NewUser) (User, error) {
	if err := nuser.Validate(); err != nil {
		return User{}, fmt.Errorf("user.register: %v", err)
	}
	user := nuser.User()
	if err := s.repo.Create(&user); err != nil {
		return User{}, fmt.Errorf("user.register: %v", err)
	}
	ctx := map[string]interface{}{
		"Name": user.FirstName + " " + user.LastName,
	}
	go email.Welcome([]string{user.Email}, ctx)
	return user, nil
}

func (s service) Login(email, password string) (User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return User{}, fmt.Errorf("user.login: %v", err)
	}
	if user.Password != calculatePassHash(password, user.Salt) {
		return New(), fmt.Errorf("user.login: %v", ErrUnauthorized)
	}
	return user, nil
}

func (s service) ResetPassword(oldpass, newpass string) error {
	return nil
}

type Middleware func(Service) Service
