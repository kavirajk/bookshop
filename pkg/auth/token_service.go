package auth

import (
	"crypto/rsa"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/kavirajk/bookshop/pkg/user"
	"github.com/pkg/errors"
)

const (
	defaultNounceLength = 30
)

// TokenService is JWT token service, responsible for generating
// and validating JWT token.
type TokenService struct {
	pubKey        *rsa.PublicKey
	privKey       *rsa.PrivateKey
	nounceLength  int
	tokenLifetime time.Duration
	issuer        string
	logger        log.Logger
}

// NewTokenService returns new TokenService with default nounceLength
func NewTokenService(issuer string, privKey *rsa.PrivateKey, pubKey *rsa.PublicKey, tokenLifetime time.Duration, logger log.Logger) *TokenService {
	return &TokenService{
		issuer:        issuer,
		privKey:       privKey,
		pubKey:        pubKey,
		tokenLifetime: tokenLifetime,
		nounceLength:  defaultNounceLength,
		logger:        logger,
	}
}

// GenerateAccessToken returns a JWT encoded with user detail.
func (ts *TokenService) GenerateAccessToken(u *user.User) (string, error) {
	utcNow := time.Now().UTC()

	claim := &Claim{
		ID:       u.ID,
		Email:    u.Email,
		IsActive: u.IsActive,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: utcNow.Add(ts.tokenLifetime).UnixNano(),
			Issuer:    ts.issuer,
			IssuedAt:  utcNow.UnixNano(),
			NotBefore: utcNow.UnixNano(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claim)

	signed, err := token.SignedString(ts.privKey)
	if err != nil {
		return "", errors.Wrap(err, "token_service.generate_access_token.signing_token_failed")
	}

	return signed, nil
}
