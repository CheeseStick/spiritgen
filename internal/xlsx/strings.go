package xlsx

import (
	"strings"
	"unicode"
)

// NormalizeString trims whitespace and removes non-space Unicode space runes,
// preserving only the regular ASCII space character.
func NormalizeString(s string) string {
	var b strings.Builder
	trimmed := strings.TrimSpace(s)
	for _, r := range trimmed {
		if !unicode.IsSpace(r) || r == ' ' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// IsBlank returns true if the normalized form of s is empty.
func IsBlank(s string) bool {
	return NormalizeString(s) == ""
}
