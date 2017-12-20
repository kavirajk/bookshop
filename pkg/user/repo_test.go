package user

import (
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/kavirajk/bookshop/util/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getRepo(t *testing.T) Repo {
	t.Helper()

	return NewRepo(log.NewLogfmtLogger(os.Stdout))
}

func TestAuthenticate(t *testing.T) {
	db, err := testutil.NewDB()
	// TODO(kaviraj): Replace with schema migrations
	db.AutoMigrate(&User{})

	if err != nil {
		t.Error(err)
	}
	repo := getRepo(t)

	hash, err := HashPassword("test")
	require.Nil(t, err)

	user := User{
		Email:        "test@bs.com",
		PasswordHash: hash,
	}

	require.Nil(t, db.Save(&user).Error)

	got, err := repo.Authenticate(db, "test@bs.com", "test")
	require.Nil(t, err)
	assert.Equal(t, "test@bs.com", got.Email)

	got, err = repo.Authenticate(db, "invalid@bs.com", "test")
	assert.Equal(t, err, ErrRepoUserNotFound)
	assert.Nil(t, got)

	got, err = repo.Authenticate(db, "test@bs.com", "invalid")
	assert.Equal(t, err, ErrRepoUserInvalidPassword)
	assert.Nil(t, got)

}
