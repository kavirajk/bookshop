package testutil

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/pkg/errors"
)

// RandomString returns base64 encoded random string.
func RandomString(len int) (string, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "RandomString failed")
	}
	return base64.URLEncoding.EncodeToString(b), err
}
