package parser

import (
	"strings"
	"unicode"
)

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
