package inmem

import (
	"errors"
	"fmt"
	"sync"

	"github.com/kavirajk/bookshop/user"
	"github.com/twinj/uuid"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("entity already exists")
)

type userRepo map[string]user.User

type userID struct {
	ID string
	sync.Mutex
}

var global userID

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

func (r userRepo) GetByEmail(email string) (user.User, error) {
	for _, v := range r {
		if v.Email == email {
			return v, nil
		}
	}
	return user.User{}, fmt.Errorf("user %v", ErrNotFound)
}

func (r userRepo) GetByToken(token string) (user.User, error) {
	for _, v := range r {
		if v.AuthToken == token {
			return v, nil
		}
	}
	return user.User{}, fmt.Errorf("user %v", ErrNotFound)
}

func (r userRepo) GetByResetKey(key string) (user.User, error) {
	for _, v := range r {
		if v.ResetKey == key {
			return v, nil
		}
	}
	return user.User{}, fmt.Errorf("user %v", ErrNotFound)
}

func (r userRepo) List() ([]user.User, error) {
	users := make([]user.User, 0)
	for _, v := range r {
		users = append(users, v)
	}
	return users, nil
}

func (r userRepo) Create(user *user.User) error {
	global.Lock()
	global.ID = uuid.NewV4().String()
	user.ID = global.ID
	global.Unlock()

	if _, ok := r[user.ID]; ok {
		return ErrAlreadyExists
	}
	r[user.ID] = *user
	return nil
}

func (r userRepo) Save(user *user.User) error {
	r[user.ID] = *user
	return nil
}

func NewUserRepo() user.Repo {
	return userRepo{}
}
