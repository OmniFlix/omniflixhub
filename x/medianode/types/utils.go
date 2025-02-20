package types

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var (
	IsAlphaNumeric   = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	IsBeginWithAlpha = regexp.MustCompile(`^[a-zA-Z].*`).MatchString
	IsAlpha          = regexp.MustCompile(`^[a-zA-Z]+`).MatchString
)

func GenUniqueID(prefix string) (string, error) {
	randomId, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return prefix + strings.ReplaceAll(randomId.String(), "-", ""), nil
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
