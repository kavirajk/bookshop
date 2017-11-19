package auth

import (
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kavirajk/bookshop/util/validator"
)

var (
	ErrClaimIDInvalid    = errors.New("claim: invalid id")
	ErrClaimEmailInvalid = errors.New("claim: invalid email")
)

// Claim represents a JWT claim.
type Claim struct {
	ID       int
	Email    string
	IsActive bool
	jwt.StandardClaims
}

// Valid implements jwt.Claims interface extending the
// validation of custom claims.
func (c *Claim) Valid() error {
	if c.ID == 0 {
		return ErrClaimIDInvalid
	}
	if !validator.IsEmailValid(c.Email) {
		return ErrClaimEmailInvalid
	}
	return c.StandardClaims.Valid()
}

// ClaimOption can be used to append any claim to the
// existing one.
type ClaimOption func(*Claim)
