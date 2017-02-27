package user

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

var (
	ErrMissingField     = errors.New("missing field")
	ErrPasswordMismatch = errors.New("passwords didn't match")
)

type User struct {
	ID        string `json:"id"`
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
		return fmt.Errorf("newuser: %v. firstname", ErrMissingField)
	}
	if n.LastName == "" {
		return fmt.Errorf("newuser: %v. lastname", ErrMissingField)
	}
	if n.Email == "" {
		return fmt.Errorf("newuser: %v. email", ErrMissingField)
	}
	if n.Password == "" {
		return fmt.Errorf("newuser: %v. password", ErrMissingField)
	}
	if n.ConfirmPassword == "" {
		return fmt.Errorf("newuser: %v. confirm_password", ErrMissingField)
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
	u.Password = calculatePassHash(u.Password, u.Salt)
	return u
}

func calculatePassHash(pass, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%s", h.Sum(nil))
}
