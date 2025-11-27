package types

import (
	"testing"
)

var (
	commentTests = []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Single line comment",
			input:    "This is a simple comment.\n",
			expected: "# This is a simple comment.\n",
		},
		{
			name:     "Multi-line comment",
			input:    "This is line one.\nThis is line two.\nThis is line three.\n",
			expected: "# This is line one.\n# This is line two.\n# This is line three.\n",
		},
	}
)

func TestCommentToSeclang(t *testing.T) {
	for _, tt := range commentTests {
		t.Run(tt.name, func(t *testing.T) {
			result := commentToSeclang(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%q\nGot:\n%q", tt.expected, result)
			}
		})
	}
}
