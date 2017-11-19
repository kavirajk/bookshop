package validator

import (
	"regexp"
	"strings"
)

var (
	validEmailRE = regexp.MustCompile(`(\S+@\S+)`)
)

func IsEmailValid(email string) bool {
	if strings.TrimSpace(email) == "" {
		return false
	}
	return validEmailRE.MatchString(email)
}
