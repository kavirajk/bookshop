package user

import (
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var (
	ErrRepoUserNotFound = errors.New("user.repo: not found")
)

type Repo interface {
	// Get returns a single user matching scopes.
	// Returns ErrRepoUserNotFound if scopes doesn't match any.
	Get(db *gorm.DB, scopes ...Scope) (*User, error)

	// Find returns slice of users matching the scopes.
	Find(db *gorm.DB, scopes ...Scope) ([]User, error)

	// Save either creates/update the user matching the scope.
	Save(db *gorm.DB, u *User, scopes ...Scope) error

	// Delete remove the users matching the scopes.
	Delete(db *gorm.DB, scopes ...Scope) error
}

// repo implements simple Repo.
type repo struct {
	logger log.Logger
}

func NewRepo(logger log.Logger) Repo {
	return &repo{logger: logger}
}

// TODO
func (r *repo) Get(db *gorm.DB, scopes ...Scope) (*User, error) {
	return nil, nil
}

func (r *repo) Find(db *gorm.DB, scopes ...Scope) ([]User, error) {
	return nil, nil
}

func (r *repo) Save(db *gorm.DB, u *User, scopes ...Scope) error {
	return nil
}

func (r *repo) Delete(db *gorm.DB, scopes ...Scope) error {
	return nil
}
