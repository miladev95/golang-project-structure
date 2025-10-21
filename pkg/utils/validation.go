package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailPattern)
	return re.MatchString(email)
}

// IsValidPhoneNumber validates phone number format (simple)
// Accepts: +1234567890, 1234567890, 123-456-7890
func IsValidPhoneNumber(phone string) bool {
	// Remove common separators
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "+", "")

	// Check if it's between 7 and 15 digits
	if !regexp.MustCompile(`^\d{7,15}$`).MatchString(phone) {
		return false
	}
	return true
}

// IsValidUsername validates username (alphanumeric, underscore, hyphen, 3-20 chars)
func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(username)
}

// IsValidPassword validates password strength
// Minimum 8 chars, at least one uppercase, one lowercase, one digit
func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)

	return hasUpper && hasLower && hasDigit
}

// IsValidURL validates URL format
func IsValidURL(url string) bool {
	const urlPattern = `^https?://[a-zA-Z0-9.-]+(:[0-9]+)?(/.*)?$`
	return regexp.MustCompile(urlPattern).MatchString(url)
}

// IsValidUUID validates UUID format (v4)
func IsValidUUID(uuid string) bool {
	const uuidPattern = `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
	return regexp.MustCompile(uuidPattern).MatchString(strings.ToLower(uuid))
}

// IsValidIP validates IPv4 address
func IsValidIP(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if !regexp.MustCompile(`^\d{1,3}$`).MatchString(part) {
			return false
		}
		// Check if value is between 0-255
		var num int
		fmt.Sscanf(part, "%d", &num)
		if num < 0 || num > 255 {
			return false
		}
	}
	return true
}

// IsStringInSlice checks if string exists in slice
func IsStringInSlice(str string, slice []string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// IsNumberBetween checks if number is between min and max
func IsNumberBetween(num, min, max int64) bool {
	return num >= min && num <= max
}

