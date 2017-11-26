package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/kavirajk/bookshop/util/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreToken(t *testing.T) {
	rp, err := testutil.SetupRedis(t)
	require.NotNil(t, rp)
	require.NoError(t, err)

	defer testutil.TearDownRedis(t)

	logger := log.NewNopLogger()

	repo := NewRepo(logger)

	var (
		token1, token2 string
	)

	token1, err = testutil.RandomString(10)
	require.NoError(t, err)

	token2, err = testutil.RandomString(10)
	require.NoError(t, err)

	cases := []struct {
		token string
		ttl   time.Duration
	}{
		{
			token: token1,
			ttl:   time.Second,
		},
		{
			token: token2,
			ttl:   time.Minute,
		},
	}

	for i, c := range cases {
		require.NoError(t, repo.SaveToken(rp, c.token, c.ttl))
		ttlGot, err := repo.GetTTL(rp, c.token)
		require.NoError(t, err)
		assert.Equal(t, c.ttl, ttlGot, fmt.Sprintf("failed %d\n", i))
	}
}
