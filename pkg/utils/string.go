package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// Slugify converts a string to a URL-friendly slug
// Example: "Hello World" -> "hello-world"
func Slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")

	// Remove non-alphanumeric characters except hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")

	// Remove leading/trailing hyphens
	s = strings.Trim(s, "-")

	// Replace multiple hyphens with single hyphen
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}

	return s
}

// TitleCase converts a string to title case
// Example: "hello world" -> "Hello World"
func TitleCase(s string) string {
	return strings.Title(strings.ToLower(s))
}

// Capitalize capitalizes the first character
// Example: "hello" -> "Hello"
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

// IsEmpty checks if string is empty or contains only whitespace
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// TruncateString truncates a string to max length with ellipsis
// Example: TruncateString("Hello World", 5) -> "Hello..."
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// ContainsWord checks if string contains a word (case-insensitive)
func ContainsWord(s, word string) bool {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	wordLower := strings.ToLower(word)
	for _, w := range words {
		if strings.ToLower(w) == wordLower {
			return true
		}
	}
	return false
}

// ReverseString reverses a string
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}