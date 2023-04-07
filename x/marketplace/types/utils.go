package types

import (
	"github.com/google/uuid"
	"regexp"
	"strings"
)

var (
	IsAlphaNumeric   = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	IsBeginWithAlpha = regexp.MustCompile(`^[a-zA-Z].*`).MatchString
	IsAlpha          = regexp.MustCompile(`^[a-zA-Z]+`).MatchString
)

func GenUniqueID(prefix string) string {
	return prefix + strings.ReplaceAll(uuid.New().String(), "-", "")
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}