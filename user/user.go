package user

import (
	"crypto/sha1"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrMissingField     = errors.New("missing field")
	ErrPasswordMismatch = errors.New("passwords didn't match")
)

type User struct {
	ID        string `json:"id" sql:"primary_key"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Salt      string `json:"-"`
	ResetKey  string `json:"-"`
	AuthToken string `json:"-"`
}

func New() User {
	u := User{}
	u.NewSalt()
	return u
}

func (u *User) NewSalt() {
	h := sha1.New()
	io.WriteString(h, strconv.Itoa(int(time.Now().UnixNano())))
	u.Salt = fmt.Sprintf("%x", h.Sum(nil))
}

type NewUser struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (n *NewUser) Validate() error {
	if n.FirstName == "" {
		return errors.Wrap(ErrMissingField, ": first_name")
	}
	if n.LastName == "" {
		return errors.Wrap(ErrMissingField, ": last_name")
	}
	if n.Email == "" {
		return errors.Wrap(ErrMissingField, ": email")
	}
	if n.Password == "" {
		return errors.Wrap(ErrMissingField, ": password")
	}
	if n.ConfirmPassword == "" {
		return errors.Wrap(ErrMissingField, ": confirm_password")
	}
	if n.Password != n.ConfirmPassword {
		return ErrPasswordMismatch
	}
	return nil
}

func (n *NewUser) User() User {
	u := New()
	u.NewSalt()
	u.FirstName = n.FirstName
	u.LastName = n.LastName
	u.Email = n.Email
	u.Password = calculatePassHash(n.Password, u.Salt)
	u.Username = strings.Split(n.Email, "@")[0]
	return u
}

func calculatePassHash(pass, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}
