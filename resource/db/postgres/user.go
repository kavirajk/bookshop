package postgres

import (
	"bytes"
	"encoding/base32"

	"github.com/jinzhu/gorm"
	"github.com/kavirajk/bookshop/db"
	"github.com/kavirajk/bookshop/user"
	_ "github.com/lib/pq"
	"github.com/pborman/uuid"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(driver, source string) (user.Repo, error) {
	db, err := gorm.Open(driver, source)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&user.User{})
	return &userRepo{db: db}, nil
}

func (r *userRepo) get(where ...interface{}) (user.User, error) {
	var u user.User
	d := r.db.New()

	if err := d.First(&u, where...).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user.New(), db.ErrNotFound
		}
		return user.New(), err
	}
	return u, nil
}

func (r *userRepo) GetByID(id string) (user.User, error) {
	return r.get("id=?", id)
}

func (r *userRepo) GetByUserName(username string) (user.User, error) {
	return r.get("username=?", username)
}

func (r *userRepo) GetByEmail(email string) (user.User, error) {
	return r.get("email=?", email)
}

func (r *userRepo) GetByToken(token string) (user.User, error) {
	return r.get("auth_token=?", token)
}

func (r *userRepo) GetByResetKey(key string) (user.User, error) {
	return r.get("reset_key=?", key)
}

func (r *userRepo) List(order string, limit, offset int) ([]user.User, int, error) {
	users := make([]user.User, 0)
	db := r.db.New()

	var total int
	if err := db.Model(&user.User{}).Order(order).Count(&total).Error; err != nil {
		return users, 0, err
	}

	err := db.Order(order).Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *userRepo) Create(u *user.User) error {
	d := r.db.New()

	if u.ID == "" {
		u.ID = NewID()
	}

	if err := d.Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepo) Save(u *user.User) error {
	d := r.db.New()

	if err := d.Save(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepo) Drop() error {
	return r.db.Exec("DELETE FROM USERS").Error
}

// Helpers

var encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769")

// NewID return global uniq indentifier that will be used as
// primary key.
func NewID() string {
	var b bytes.Buffer
	encoder := base32.NewEncoder(encoding, &b)
	encoder.Write(uuid.NewRandom())
	encoder.Close()
	b.Truncate(26)
	return b.String()
}
