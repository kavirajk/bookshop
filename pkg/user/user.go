package user

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrMissingField     = errors.New("missing field")
	ErrPasswordMismatch = errors.New("passwords didn't match")
)

// User represents domain model of user service.
type User struct {
	ID           int    `json:"id" sql:"primary_key"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	IsActive     bool   `json:"-"`
	PasswordHash string `json:"-"`
}

func (_ User) TableName() string {
	return "bsuser"
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

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
