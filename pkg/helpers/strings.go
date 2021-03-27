package helpers

import "strings"

func IsEmptyOrWhitespace(s string) bool {
	s = strings.ReplaceAll(s, " ", "")
	return len(s) == 0
}

func IsNilEmptyOrWhitespace(s *string) bool {
	if s == nil {
		return true
	}

	return IsEmptyOrWhitespace(*s)
}
