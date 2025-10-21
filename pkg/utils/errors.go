package utils

import (
	"fmt"
)

// AppError represents a custom application error
type AppError struct {
	Code    string
	Message string
	Details interface{}
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (error: %v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAppError creates a new app error
func NewAppError(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewAppErrorWithCause creates a new app error with cause
func NewAppErrorWithCause(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// SetDetails sets additional error details
func (e *AppError) SetDetails(details interface{}) *AppError {
	e.Details = details
	return e
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors
type ValidationErrors struct {
	Code   string             `json:"code"`
	Errors []ValidationError  `json:"errors"`
}

func (ve ValidationErrors) Error() string {
	return fmt.Sprintf("validation error: %d field(s) failed validation", len(ve.Errors))
}

// NewValidationErrors creates a new validation errors collection
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Code:   "VALIDATION_ERROR",
		Errors: []ValidationError{},
	}
}

// Add adds a validation error
func (ve *ValidationErrors) Add(field, message string) *ValidationErrors {
	ve.Errors = append(ve.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
	return ve
}

// AddWithValue adds a validation error with the invalid value
func (ve *ValidationErrors) AddWithValue(field, message string, value interface{}) *ValidationErrors {
	ve.Errors = append(ve.Errors, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
	return ve
}

// HasErrors checks if there are validation errors
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// NotFoundError represents a not found error
type NotFoundError struct {
	Resource string
	ID       interface{}
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with id %v not found", e.Resource, e.ID)
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource string, id interface{}) NotFoundError {
	return NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// ConflictError represents a conflict error (e.g., duplicate entry)
type ConflictError struct {
	Message string
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s", e.Message)
}

// NewConflictError creates a new conflict error
func NewConflictError(message string) ConflictError {
	return ConflictError{
		Message: message,
	}
}

// UnauthorizedError represents an unauthorized error
type UnauthorizedError struct {
	Message string
}

func (e UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %s", e.Message)
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) UnauthorizedError {
	return UnauthorizedError{
		Message: message,
	}
}

// ForbiddenError represents a forbidden error
type ForbiddenError struct {
	Message string
}

func (e ForbiddenError) Error() string {
	return fmt.Sprintf("forbidden: %s", e.Message)
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string) ForbiddenError {
	return ForbiddenError{
		Message: message,
	}
}

// InternalServerError represents an internal server error
type InternalServerError struct {
	Message string
	Err     error
}

func (e InternalServerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("internal server error: %s (cause: %v)", e.Message, e.Err)
	}
	return fmt.Sprintf("internal server error: %s", e.Message)
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(message string, err error) InternalServerError {
	return InternalServerError{
		Message: message,
		Err:     err,
	}
}