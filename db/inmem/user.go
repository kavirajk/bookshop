package inmem

import (
	"errors"
	"fmt"

	"github.com/kavirajk/bookshop/user"
)

var (
	ErrNotFound = errors.New("not found")
)

type userRepo map[string]user.User

func (r userRepo) GetByID(ID string) (user.User, error) {
	u, ok := r[ID]
	if !ok {
		return user.User{}, fmt.Errorf("user %v", ErrNotFound)
	}
	return u, nil
}

func (r userRepo) GetByUserName(username string) (user.User, error) {
	for _, v := range r {
		if v.Username == username {
			return v, nil
		}
	}
	return user.User{}, fmt.Errorf("user %v", ErrNotFound)
}

func (r userRepo) List() ([]user.User, error) {
	var users []user.User
	for _, v := range r {
		users = append(users, v)
	}
	return users, nil
}

func NewUserRepo() user.Repo {
	return userRepo{}
}
