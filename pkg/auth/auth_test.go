package auth

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/kavirajk/bookshop/pkg/user"
	"github.com/kavirajk/bookshop/util/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testTokenService(t *testing.T, logger log.Logger) *TokenService {
	t.Helper()
	pubKey, err := LoadPublicKey("testdata/public.pem")
	require.NoError(t, err)

	privKey, err := LoadPrivateKey("testdata/private.pem")
	require.NoError(t, err)

	return NewTokenService("bookshop", privKey, pubKey, 10*time.Second, logger)
}

func TestRegister(t *testing.T) {
	logger := log.NewLogfmtLogger(os.Stdout)

	db, err := testutil.NewDB(t, logger)
	require.NoError(t, err)
	require.NotNil(t, db)

	defer testutil.FlushDB(t, db, "users")

	ts := testTokenService(t, logger)

	svc := NewService(logger, ts, db)
	ctx := context.Background()

	// Existing user
	require.NoError(t, db.Create(&user.User{Email: "olduser@gmail.com", IsActive: true}).Error)
	assert.Equal(t, svc.Register(ctx, &user.User{Email: "olduser@gmail.com"}), ErrUserExists)

	// Inactive user
	require.NoError(t, db.Create(&user.User{Email: "inactive-user@gmail.com"}).Error)
	assert.Equal(t, svc.Register(ctx, &user.User{Email: "inactive-user@gmail.com"}), ErrUserInactive)

	// New user (Ideal case)
	require.NoError(t, svc.Register(ctx, &user.User{Email: "new-user@gmail.com"}))

}

func TestLogin(t *testing.T) {
	logger := log.NewLogfmtLogger(os.Stdout)

	db, err := testutil.NewDB(t, logger)
	require.NoError(t, err)
	require.NotNil(t, db)

	defer testutil.FlushDB(t, db, "users")

	ts := testTokenService(t, logger)

	svc := NewService(logger, ts, db)
	ctx := context.Background()

	pass, err := user.HashPassword("user1")
	user1 := user.User{Email: "user1@gmail.com", IsActive: false, PasswordHash: pass}

	// Inactive user
	require.NoError(t, err)
	require.NoError(t, db.Create(&user1).Error)
	bundle, err := svc.Login(ctx, "user1@gmail.com", "user1")
	assert.Equal(t, err, ErrUserInactive)
	require.Nil(t, bundle)

	// Non-exists user
	bundle, err = svc.Login(ctx, "non-exists@gmail.com", "non-exists")
	assert.Equal(t, err, ErrInvalidCredentials)
	require.Nil(t, bundle)

	// Active and Invalid Credentials
	require.NoError(t, db.Model(&user1).Update("is_active", true).Error)
	bundle, err = svc.Login(ctx, "user1@gmail.com", "invalid-password")
	assert.Equal(t, ErrInvalidCredentials, err)
	require.Nil(t, bundle)

	// Active and Valid credentials
	bundle, err = svc.Login(ctx, "user1@gmail.com", "user1")
	require.NoError(t, err)
	require.NotNil(t, bundle)

}
