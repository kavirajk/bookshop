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

// User represents domain model of user service.
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

// New create empty user with random salt.
func New() User {
	u := User{}
	u.NewSalt()
	return u
}

// NewSalt creates random salt value based on current time.
// usefull to hash password in a secure way.
func (u *User) NewSalt() {
	h := sha1.New()
	io.WriteString(h, strconv.Itoa(int(time.Now().UnixNano())))
	u.Salt = fmt.Sprintf("%x", h.Sum(nil))
}

// NewUser represents user who is about to register.
type NewUser struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// Validate does basic validation before saving into db.
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

// User map NewUser with domain User.
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
