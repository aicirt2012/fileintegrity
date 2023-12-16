package path

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsIgnoredFile(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Ignored win file",
			input:    "~$",
			expected: true,
		},
		{
			name:     "Ignored win file",
			input:    "~$any name",
			expected: true,
		},
		{
			name:     "Ignored macos file",
			input:    "._",
			expected: true,
		},
		{
			name:     "Ignored macos file",
			input:    "._any name",
			expected: true,
		},
		{
			name:     "Ignored macos file",
			input:    ".DS_Store",
			expected: true,
		},
		{
			name:     "Valid hidden file",
			input:    ".",
			expected: false,
		},
		{
			name:     "Valid hidden file",
			input:    ".some hidden file",
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, isIgnoredFile(c.input))
		})
	}
}
