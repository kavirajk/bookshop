package auth

import (
	"crypto/rsa"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

func LoadPrivateKey(filepath string) (*rsa.PrivateKey, error) {
	c, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "loadPrivateKey.ReadFailed")
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(c)
	if err != nil {
		return nil, errors.Wrap(err, "loadPrivateKey.parseFailed")
	}
	return privKey, err
}

func LoadPublicKey(filepath string) (*rsa.PublicKey, error) {
	c, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "loadPublicKey.ReadFailed")
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(c)
	if err != nil {
		return nil, errors.Wrap(err, "loadPublicKey.parseFailed")
	}
	return pubKey, err
}
