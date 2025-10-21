package tests

import (
	"testing"

	"github.com/miladev95/golang-project-structure/pkg/utils"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "valid email",
			email:    "test@example.com",
			expected: true,
		},
		{
			name:     "valid email with numbers",
			email:    "test123@example.com",
			expected: true,
		},
		{
			name:     "valid email with dots",
			email:    "test.user@example.com",
			expected: true,
		},
		{
			name:     "invalid - no @",
			email:    "testexample.com",
			expected: false,
		},
		{
			name:     "invalid - no domain",
			email:    "test@",
			expected: false,
		},
		{
			name:     "invalid - no local part",
			email:    "@example.com",
			expected: false,
		},
		{
			name:     "invalid - no TLD",
			email:    "test@example",
			expected: false,
		},
		{
			name:     "empty string",
			email:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IsValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("IsValidEmail(%q) = %v, want %v", tt.email, result, tt.expected)
			}
		})
	}
}

func TestIsValidPhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		{
			name:     "valid phone",
			phone:    "1234567890",
			expected: true,
		},
		{
			name:     "valid with country code",
			phone:    "+1234567890",
			expected: true,
		},
		{
			name:     "valid with dashes",
			phone:    "123-456-7890",
			expected: true,
		},
		{
			name:     "valid with spaces",
			phone:    "123 456 7890",
			expected: true,
		},
		{
			name:     "too short",
			phone:    "123456",
			expected: false,
		},
		{
			name:     "invalid characters",
			phone:    "123-456-789a",
			expected: false,
		},
		{
			name:     "empty string",
			phone:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IsValidPhoneNumber(tt.phone)
			if result != tt.expected {
				t.Errorf("IsValidPhoneNumber(%q) = %v, want %v", tt.phone, result, tt.expected)
			}
		})
	}
}

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expected bool
	}{
		{
			name:     "valid username",
			username: "john_doe",
			expected: true,
		},
		{
			name:     "valid with numbers",
			username: "user123",
			expected: true,
		},
		{
			name:     "valid with hyphen",
			username: "user-name",
			expected: true,
		},
		{
			name:     "too short",
			username: "ab",
			expected: false,
		},
		{
			name:     "too long",
			username: "abcdefghijklmnopqrstu",
			expected: false,
		},
		{
			name:     "invalid characters",
			username: "user@name",
			expected: false,
		},
		{
			name:     "with space",
			username: "user name",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IsValidUsername(tt.username)
			if result != tt.expected {
				t.Errorf("IsValidUsername(%q) = %v, want %v", tt.username, result, tt.expected)
			}
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "valid password",
			password: "Secure123",
			expected: true,
		},
		{
			name:     "valid with special chars",
			password: "Secure123!@#",
			expected: true,
		},
		{
			name:     "too short",
			password: "Pass1",
			expected: false,
		},
		{
			name:     "no uppercase",
			password: "password123",
			expected: false,
		},
		{
			name:     "no lowercase",
			password: "PASSWORD123",
			expected: false,
		},
		{
			name:     "no digit",
			password: "PasswordTest",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IsValidPassword(tt.password)
			if result != tt.expected {
				t.Errorf("IsValidPassword(%q) = %v, want %v", tt.password, result, tt.expected)
			}
		})
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "valid http URL",
			url:      "http://example.com",
			expected: true,
		},
		{
			name:     "valid https URL",
			url:      "https://example.com",
			expected: true,
		},
		{
			name:     "with path",
			url:      "https://example.com/path",
			expected: true,
		},
		{
			name:     "with port",
			url:      "https://example.com:8080",
			expected: true,
		},
		{
			name:     "with path and port",
			url:      "https://example.com:8080/api/users",
			expected: true,
		},
		{
			name:     "invalid - no protocol",
			url:      "example.com",
			expected: false,
		},
		{
			name:     "invalid - ftp protocol",
			url:      "ftp://example.com",
			expected: false,
		},
		{
			name:     "empty string",
			url:      "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IsValidURL(tt.url)
			if result != tt.expected {
				t.Errorf("IsValidURL(%q) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestIsValidIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{
			name:     "valid IP",
			ip:       "192.168.1.1",
			expected: true,
		},
		{
			name:     "valid localhost",
			ip:       "127.0.0.1",
			expected: true,
		},
		{
			name:     "valid 0.0.0.0",
			ip:       "0.0.0.0",
			expected: true,
		},
		{
			name:     "valid 255.255.255.255",
			ip:       "255.255.255.255",
			expected: true,
		},
		{
			name:     "too many parts",
			ip:       "192.168.1.1.1",
			expected: false,
		},
		{
			name:     "too few parts",
			ip:       "192.168.1",
			expected: false,
		},
		{
			name:     "out of range",
			ip:       "256.1.1.1",
			expected: false,
		},
		{
			name:     "invalid characters",
			ip:       "192.168.1.a",
			expected: false,
		},
		{
			name:     "empty string",
			ip:       "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IsValidIP(tt.ip)
			if result != tt.expected {
				t.Errorf("IsValidIP(%q) = %v, want %v", tt.ip, result, tt.expected)
			}
		})
	}
}

func TestIsStringInSlice(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		slice    []string
		expected bool
	}{
		{
			name:     "string exists",
			str:      "apple",
			slice:    []string{"apple", "banana", "orange"},
			expected: true,
		},
		{
			name:     "string not exists",
			str:      "grape",
			slice:    []string{"apple", "banana", "orange"},
			expected: false,
		},
		{
			name:     "empty slice",
			str:      "apple",
			slice:    []string{},
			expected: false,
		},
		{
			name:     "empty string",
			str:      "",
			slice:    []string{"", "apple"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IsStringInSlice(tt.str, tt.slice)
			if result != tt.expected {
				t.Errorf("IsStringInSlice(%q, %v) = %v, want %v", tt.str, tt.slice, result, tt.expected)
			}
		})
	}
}

func TestIsNumberBetween(t *testing.T) {
	tests := []struct {
		name     string
		num      int64
		min      int64
		max      int64
		expected bool
	}{
		{
			name:     "number in range",
			num:      5,
			min:      1,
			max:      10,
			expected: true,
		},
		{
			name:     "number at min",
			num:      1,
			min:      1,
			max:      10,
			expected: true,
		},
		{
			name:     "number at max",
			num:      10,
			min:      1,
			max:      10,
			expected: true,
		},
		{
			name:     "number below min",
			num:      0,
			min:      1,
			max:      10,
			expected: false,
		},
		{
			name:     "number above max",
			num:      11,
			min:      1,
			max:      10,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IsNumberBetween(tt.num, tt.min, tt.max)
			if result != tt.expected {
				t.Errorf("IsNumberBetween(%d, %d, %d) = %v, want %v", tt.num, tt.min, tt.max, result, tt.expected)
			}
		})
	}
}
