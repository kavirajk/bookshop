package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmailValid(t *testing.T) {
	cases := []struct {
		email   string
		isValid bool
	}{
		{"neo@matrix.com", true},
		{"morpheus.com", false},
		{"", false},
		{"2", false},
		{"mr-anderson@matrix.com", true},
	}

	for _, c := range cases {
		assert.Equal(t, IsEmailValid(c.email), c.isValid)
	}
}
