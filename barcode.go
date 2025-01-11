package simutils

import (
	"strings"
	"unicode"
)

func IsDigit(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}

	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
