package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "trims spaces",
			input:    "   hello world   ",
			expected: "hello world",
		},
		{
			name:     "removes control characters",
			input:    "hi\x00there\x1F",
			expected: "hithere",
		},
		{
			name:     "removes carriage return",
			input:    "hello\r\nworld",
			expected: "helloworld",
		},
		{
			name: "unicode normalization NFC",
			// "é" can be decomposed (e + ◌́). Normalization makes it one rune (é).
			input:    "Cafe\u0301",
			expected: "Café",
		},
		{
			name:     "no cleanup needed",
			input:    "simple text",
			expected: "simple text",
		},
		{
			name:     "only control characters",
			input:    "\x00\x07\r\x1F",
			expected: "",
		},
		{
			name:     "mixed safe and unsafe",
			input:    "keep this but remove\x00that",
			expected: "keep this but removethat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, SanitizeInput(tt.input), tt.expected)
		})
	}
}
