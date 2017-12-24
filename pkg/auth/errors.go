package auth

import "github.com/pkg/errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenNotFound      = errors.New("token not found")
	ErrTTLNotFound        = errors.New("ttl not found")
	ErrClaimIDInvalid     = errors.New("claim: invalid id")
	ErrClaimEmailInvalid  = errors.New("claim: invalid email")
)
