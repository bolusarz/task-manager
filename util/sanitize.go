package util

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func SanitizeInput(s string) string {
	// Trim spaces
	s = strings.TrimSpace(s)

	// Normalize Unicode (e.g. é -> single code point)
	s = norm.NFC.String(s)

	// Remove control chars except newline/tab
	var b strings.Builder
	for _, r := range s {
		if unicode.IsControl(r) {
			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}
