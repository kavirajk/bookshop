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

func (r *repo) Get(db *gorm.DB, scopes ...Scope) (*User, error) {
	db = db.Scopes(scopes...)

	var u User
	if err := db.First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRepoUserNotFound
		}
		return nil, errors.Wrapf(err, "user.repo.Get")
	}

	return &u, nil
}

func (r *repo) Find(db *gorm.DB, scopes ...Scope) ([]User, error) {
	db = db.Scopes(scopes...)

	var users []User

	if err := db.Find(&users).Error; err != nil {
		return nil, errors.Wrap(err, "user.repo.Find")
	}
	return users, nil
}

func (r *repo) Save(db *gorm.DB, u *User, scopes ...Scope) error {
	db = db.Scopes(scopes...)

	if err := db.Save(u).Error; err != nil {
		return errors.Wrap(err, "user.repo.Save")
	}
	return nil
}

func (r *repo) Delete(db *gorm.DB, scopes ...Scope) error {
	db = db.Scopes(scopes...)

	if err := db.Delete(&User{}); err != nil {
		return errors.Wrap(err, "user.repo.Delete")
	}
	return nil
}
