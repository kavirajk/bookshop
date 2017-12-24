package auth

import (
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/kavirajk/bookshop/pkg/user"
	"github.com/stretchr/testify/require"
)

func TestGenerateAccessToken(t *testing.T) {
	pubKey, err := LoadPublicKey("testdata/public.pem")
	require.NoError(t, err)

	privKey, err := LoadPrivateKey("testdata/private.pem")
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
