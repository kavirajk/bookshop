package auth

import (
	"crypto/rsa"
	"io/ioutil"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/kavirajk/bookshop/pkg/user"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestGenerateAccessToken(t *testing.T) {
	pubKey, err := loadPublicKey(t, "testdata/public.pem")
	require.NoError(t, err)

	privKey, err := loadPrivateKey(t, "testdata/private.pem")
	require.NoError(t, err)

	logger := log.NewNopLogger()

	ts := NewTokenService("bookshop", privKey, pubKey, 10*time.Second, logger)
	cases := []struct {
		u          *user.User
		shouldFail bool
	}{
		{
			u:          &user.User{ID: 3, Email: "tony@stark.com"},
			shouldFail: false,
		},
		{
			u:          &user.User{ID: 0, Email: "tony@stark.com"},
			shouldFail: false,
		},
	}

	for _, c := range cases {
		token, err := ts.GenerateAccessToken(c.u)
		if c.shouldFail {
			require.Error(t, err)
			continue
		}
		require.Nil(t, err)
		t.Log("token: ", token)
	}

}

func loadPrivateKey(t *testing.T, filepath string) (*rsa.PrivateKey, error) {
	t.Helper()
	c, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "loadPrivateKey.ReadFailed")
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(c)
	if err != nil {
		return nil, errors.Wrap(err, "loadPrivateKey.parseFailed")
	}
	return privKey, err
}

func loadPublicKey(t *testing.T, filepath string) (*rsa.PublicKey, error) {
	t.Helper()
	c, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "loadPublicKey.ReadFailed")
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(c)
	if err != nil {
		return nil, errors.Wrap(err, "loadPublicKey.parseFailed")
	}
	return pubKey, err
}
