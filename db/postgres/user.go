package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/kavirajk/bookshop/user"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(source string) (user.Repo, error) {
	db, err := gorm.Open("postgres", source)
	if err != nil {
		return nil, err
	}
	return &userRepo{db: db}, nil
}

func (r *userRepo) GetByID(id string) (user.User, error) {
	return user.New(), nil
}

func (r *userRepo) GetByUserName(username string) (user.User, error) {
	return user.New(), nil
}

func (r *userRepo) GetByEmail(email string) (user.User, error) {
	return user.New(), nil
}

func (r *userRepo) GetByToken(token string) (user.User, error) {
	return user.New(), nil
}

func (r *userRepo) GetByResetKey(token string) (user.User, error) {
	return user.New(), nil
}

func (r *userRepo) List() ([]user.User, error) {
	return nil, nil
}

func (r *userRepo) Create(u *user.User) error {
	return nil
}

func (r *userRepo) Save(u *user.User) error {
	return nil
}
