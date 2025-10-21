package tests

import (
	"errors"
	"strings"
	"testing"

	"github.com/yourusername/yourproject/pkg/utils"
)

func TestAppError(t *testing.T) {
	t.Run("new app error", func(t *testing.T) {
		err := utils.NewAppError("INVALID_INPUT", "Input is invalid")
		if err.Code != "INVALID_INPUT" {
			t.Errorf("Code: got %s, want INVALID_INPUT", err.Code)
		}

		if err.Message != "Input is invalid" {
			t.Errorf("Message: got %s, want 'Input is invalid'", err.Message)
		}

		if err.Err != nil {
			t.Error("Expected Err to be nil")
		}
	})

	t.Run("new app error with cause", func(t *testing.T) {
		cause := errors.New("database connection failed")
		err := utils.NewAppErrorWithCause("DB_ERROR", "Failed to connect", cause)

		if err.Code != "DB_ERROR" {
			t.Errorf("Code: got %s, want DB_ERROR", err.Code)
		}

		if err.Err != cause {
			t.Error("Expected Err to match cause")
		}
	})

	t.Run("set details", func(t *testing.T) {
		err := utils.NewAppError("CODE", "message").SetDetails(map[string]string{"field": "value"})

		if err.Details == nil {
			t.Error("Expected Details to be set")
		}
	})

	t.Run("error string without cause", func(t *testing.T) {
		err := utils.NewAppError("CODE", "message")
		expected := "CODE: message"

		if err.Error() != expected {
			t.Errorf("Error(): got %s, want %s", err.Error(), expected)
		}
	})

	t.Run("error string with cause", func(t *testing.T) {
		cause := errors.New("cause")
		err := utils.NewAppErrorWithCause("CODE", "message", cause)
		str := err.Error()

		if !strings.Contains(str, "CODE") || !strings.Contains(str, "message") || !strings.Contains(str, "cause") {
			t.Errorf("Error() should contain code, message, and cause: %s", str)
		}
	})
}

func TestNotFoundError(t *testing.T) {
	err := utils.NewNotFoundError("User", 123)

	if err.Resource != "User" {
		t.Errorf("Resource: got %s, want User", err.Resource)
	}

	if err.ID != 123 {
		t.Errorf("ID: got %v, want 123", err.ID)
	}

	expected := "User with id 123 not found"
	if err.Error() != expected {
		t.Errorf("Error(): got %s, want %s", err.Error(), expected)
	}
}

func TestConflictError(t *testing.T) {
	err := utils.NewConflictError("Email already exists")

	if err.Message != "Email already exists" {
		t.Errorf("Message: got %s, want 'Email already exists'", err.Message)
	}

	expected := "conflict: Email already exists"
	if err.Error() != expected {
		t.Errorf("Error(): got %s, want %s", err.Error(), expected)
	}
}

func TestUnauthorizedError(t *testing.T) {
	err := utils.NewUnauthorizedError("Invalid token")

	if err.Message != "Invalid token" {
		t.Errorf("Message: got %s, want 'Invalid token'", err.Message)
	}

	expected := "unauthorized: Invalid token"
	if err.Error() != expected {
		t.Errorf("Error(): got %s, want %s", err.Error(), expected)
	}
}

func TestForbiddenError(t *testing.T) {
	err := utils.NewForbiddenError("You don't have permission")

	if err.Message != "You don't have permission" {
		t.Errorf("Message: got %s, want 'You don't have permission'", err.Message)
	}

	expected := "forbidden: You don't have permission"
	if err.Error() != expected {
		t.Errorf("Error(): got %s, want %s", err.Error(), expected)
	}
}

func TestInternalServerError(t *testing.T) {
	t.Run("without cause", func(t *testing.T) {
		err := utils.NewInternalServerError("Processing failed", nil)

		expected := "internal server error: Processing failed"
		if err.Error() != expected {
			t.Errorf("Error(): got %s, want %s", err.Error(), expected)
		}
	})

	t.Run("with cause", func(t *testing.T) {
		cause := errors.New("database error")
		err := utils.NewInternalServerError("Processing failed", cause)
		str := err.Error()

		if !strings.Contains(str, "internal server error") || !strings.Contains(str, "Processing failed") || !strings.Contains(str, "database error") {
			t.Errorf("Error() should contain all parts: %s", str)
		}
	})
}

func TestValidationErrors(t *testing.T) {
	t.Run("new validation errors", func(t *testing.T) {
		ve := utils.NewValidationErrors()

		if ve.Code != "VALIDATION_ERROR" {
			t.Errorf("Code: got %s, want VALIDATION_ERROR", ve.Code)
		}

		if len(ve.Errors) != 0 {
			t.Error("Expected no errors initially")
		}
	})

	t.Run("add single error", func(t *testing.T) {
		ve := utils.NewValidationErrors()
		result := ve.Add("email", "Invalid format")

		if result != ve {
			t.Error("Add should return the same validator for chaining")
		}

		if len(ve.Errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(ve.Errors))
		}

		if ve.Errors[0].Field != "email" {
			t.Errorf("Field: got %s, want email", ve.Errors[0].Field)
		}

		if ve.Errors[0].Message != "Invalid format" {
			t.Errorf("Message: got %s, want 'Invalid format'", ve.Errors[0].Message)
		}
	})

	t.Run("add multiple errors", func(t *testing.T) {
		ve := utils.NewValidationErrors()
		ve.Add("email", "Invalid format").
			Add("password", "Too short").
			Add("name", "Required")

		if len(ve.Errors) != 3 {
			t.Errorf("Expected 3 errors, got %d", len(ve.Errors))
		}
	})

	t.Run("add with value", func(t *testing.T) {
		ve := utils.NewValidationErrors()
		ve.AddWithValue("age", "Must be over 18", 16)

		if len(ve.Errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(ve.Errors))
		}

		if ve.Errors[0].Value != 16 {
			t.Errorf("Value: got %v, want 16", ve.Errors[0].Value)
		}
	})

	t.Run("has errors", func(t *testing.T) {
		ve := utils.NewValidationErrors()

		if ve.HasErrors() {
			t.Error("Should not have errors initially")
		}

		ve.Add("field", "error")

		if !ve.HasErrors() {
			t.Error("Should have errors after adding one")
		}
	})

	t.Run("error message", func(t *testing.T) {
		ve := utils.NewValidationErrors()
		ve.Add("field1", "error1")
		ve.Add("field2", "error2")
		ve.Add("field3", "error3")

		expected := "validation error: 3 field(s) failed validation"
		if ve.Error() != expected {
			t.Errorf("Error(): got %s, want %s", ve.Error(), expected)
		}
	})
}

func TestValidationError(t *testing.T) {
	err := utils.ValidationError{
		Field:   "email",
		Message: "Invalid email format",
	}

	expected := "validation error on field 'email': Invalid email format"
	if err.Error() != expected {
		t.Errorf("Error(): got %s, want %s", err.Error(), expected)
	}
}