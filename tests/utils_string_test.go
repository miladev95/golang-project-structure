package tests

import (
	"testing"

	"github.com/miladev95/golang-project-structure/pkg/utils"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple string",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "with special characters",
			input:    "Hello! @World #Test",
			expected: "hello-world-test",
		},
		{
			name:     "with multiple spaces",
			input:    "Hello   World",
			expected: "hello-world",
		},
		{
			name:     "with numbers",
			input:    "Hello World 123",
			expected: "hello-world-123",
		},
		{
			name:     "already slug",
			input:    "hello-world",
			expected: "hello-world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.Slugify(tt.input)
			if result != tt.expected {
				t.Errorf("Slugify(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTitleCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase string",
			input:    "hello world",
			expected: "Hello World",
		},
		{
			name:     "already title case",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "uppercase string",
			input:    "HELLO WORLD",
			expected: "Hello World",
		},
		{
			name:     "mixed case",
			input:    "HeLLo WoRLd",
			expected: "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.TitleCase(tt.input)
			if result != tt.expected {
				t.Errorf("TitleCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase string",
			input:    "hello",
			expected: "Hello",
		},
		{
			name:     "already capitalized",
			input:    "Hello",
			expected: "Hello",
		},
		{
			name:     "single character",
			input:    "h",
			expected: "H",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.Capitalize(tt.input)
			if result != tt.expected {
				t.Errorf("Capitalize(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "spaces only",
			input:    "   ",
			expected: true,
		},
		{
			name:     "tabs and newlines",
			input:    "\t\n",
			expected: true,
		},
		{
			name:     "non-empty string",
			input:    "hello",
			expected: false,
		},
		{
			name:     "string with spaces",
			input:    "  hello  ",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IsEmpty(tt.input)
			if result != tt.expected {
				t.Errorf("IsEmpty(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		expected  string
	}{
		{
			name:      "string longer than max",
			input:     "Hello World",
			maxLength: 5,
			expected:  "Hello...",
		},
		{
			name:      "string equal to max",
			input:     "Hello",
			maxLength: 5,
			expected:  "Hello",
		},
		{
			name:      "string shorter than max",
			input:     "Hi",
			maxLength: 5,
			expected:  "Hi",
		},
		{
			name:      "max length zero",
			input:     "Hello",
			maxLength: 0,
			expected:  "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.TruncateString(tt.input, tt.maxLength)
			if result != tt.expected {
				t.Errorf("TruncateString(%q, %d) = %q, want %q", tt.input, tt.maxLength, result, tt.expected)
			}
		})
	}
}

func TestContainsWord(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		word     string
		expected bool
	}{
		{
			name:     "word exists",
			input:    "hello world",
			word:     "world",
			expected: true,
		},
		{
			name:     "word does not exist",
			input:    "hello world",
			word:     "foo",
			expected: false,
		},
		{
			name:     "case insensitive",
			input:    "Hello World",
			word:     "world",
			expected: true,
		},
		{
			name:     "partial word match",
			input:    "hello world",
			word:     "wor",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ContainsWord(tt.input, tt.word)
			if result != tt.expected {
				t.Errorf("ContainsWord(%q, %q) = %v, want %v", tt.input, tt.word, result, tt.expected)
			}
		})
	}
}

func TestReverseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple string",
			input:    "hello",
			expected: "olleh",
		},
		{
			name:     "string with spaces",
			input:    "hello world",
			expected: "dlrow olleh",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "a",
		},
		{
			name:     "unicode characters",
			input:    "café",
			expected: "éfac",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ReverseString(tt.input)
			if result != tt.expected {
				t.Errorf("ReverseString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
