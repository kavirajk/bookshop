package auth

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClaimValid(t *testing.T) {
	cases := []struct {
		cm         Claim
		shouldFail bool
	}{
		{Claim{ID: 0}, true},
		{Claim{ID: 1, Email: "   "}, true},
		{Claim{ID: 1, Email: "\n"}, true},
		{Claim{ID: 1, Email: ""}, true},
		{Claim{ID: 0, Email: ""}, true},
		{Claim{ID: 1, Email: "trinity.com"}, true},
		{Claim{ID: 1, Email: "trinity@matrix.com"}, false},
	}

	for _, c := range cases {
		err := c.cm.Valid()
		if c.shouldFail {
			require.Error(t, err, c.cm.Email, fmt.Sprintf("%d", c.cm.ID))
			continue
		}
		require.Nil(t, err)
	}
}
