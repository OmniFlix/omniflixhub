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

func GenUniqueID(prefix string) string {
	return prefix + strings.ReplaceAll(uuid.New().String(), "-", "")
}

func IsIBCDenom(denomID string) bool {
	return strings.HasPrefix(denomID, "ibc/")
}
