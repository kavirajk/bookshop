package auth

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"github.com/kavirajk/bookshop/pkg/notification"
	"github.com/kavirajk/bookshop/pkg/user"
)

// Service provides an auth service.
type Service interface {
	// Register takes new user, validates it and save it in user store.
	Register(ctx context.Context, user *user.User) error

	// Login returns auth Bundle containing tokens.
	Login(ctx context.Context, email, password string) (*Bundle, error)

	// DecodeAccessToken decodes access token to claims.
	DecodeAccessToken(ctx context.Context, token string) (*Claim, error)

	// Impersonate takes claims of privileged user and userID to impersonate.
	// And returns auth bundle of userID if validation is success.
	Impersonate(ctx context.Context, claim *Claim, userID string) (*Bundle, error)
}

// Bundle is collection of auth related data, usually returned to client.
type Bundle struct {
	AccessToken, RefreshToken string
}

// service is a basic service that implement auth Service interface.
type service struct {
	logger   log.Logger
	ts       *TokenService
	userRepo user.Repo
	db       *gorm.DB
	notifier notification.Notifier
}

// NewService returns new auth Service.
func NewService(logger log.Logger, ts *TokenService, db *gorm.DB) Service {
	return &service{
		logger:   logger,
		ts:       ts,
		userRepo: user.NewRepo(logger),
		db:       db,
	}
}

// Register service persists the new user in the data store and
// notifies in the given notifier. Register assumes all the user data
// is already sanatized and validated.
func (s *service) Register(ctx context.Context, user *user.User) error {
	olduser, err := s.userRepo.GetByEmail(s.db, user.Email)
	if err == nil {
		if !olduser.IsActive {
			return ErrUserInactive
		}
		return ErrUserExists
	}

	return s.userRepo.Save(s.db, user)
}

// Login service assumes all its arguments are already validated.
// Its developer work to validate user payload before calling the service login.
func (s *service) Login(ctx context.Context, email, password string) (*Bundle, error) {
	u, err := s.userRepo.Authenticate(s.db, email, password)
	if err == user.ErrRepoUserNotFound || err == user.ErrRepoUserInvalidPassword {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if !u.IsActive {
		return nil, ErrUserInactive
	}
	token, err := s.ts.GenerateAccessToken(u)
	if err != nil {
		return nil, err
	}
	return &Bundle{AccessToken: token}, nil
}
func (s *service) DecodeAccessToken(ctx context.Context, accessToken string) (*Claim, error) {
	return nil, nil
}
func (s *service) Impersonate(ctx context.Context, claim *Claim, userID string) (*Bundle, error) {
	return nil, nil
}
