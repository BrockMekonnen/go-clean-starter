package errors

import (
	"github.com/go-playground/validator/v10"
)

// ValidationError represents an error due to failed validation.
type ValidationError struct {
	*BaseError[ValidationErrorProps] // embedding BaseError for reusability
}

// ValidationErrorProps holds additional data associated with a validation error.
type ValidationErrorProps struct {
	Target string                  `json:"target"`
	Error  validator.ValidationErrors `json:"error"`
}

var (
	// The symbol to identify the error type.
	ValidationType = symbol("ValidationError")
	// Default name for the error.
	ValidationName = "ValidationError"
)

// NewValidationError creates a new ValidationError instance.
func NewValidationError(message string, target string, err validator.ValidationErrors) *ValidationError {
	// Create the props object for the error
	props := ValidationErrorProps{
		Target: target,
		Error:  err,
	}

	// Create a new BaseError with the specific type
	baseError := NewBaseError[ValidationErrorProps](ValidationName, message, ValidationName, props, ValidationType)
	return &ValidationError{
		BaseError: baseError,
	}
}

// Is checks if the error is of type ValidationError.
func (e *ValidationError) Is(target error) bool {
	// This checks if the target error is of type ValidationError
	if _, ok := target.(*ValidationError); ok {
		return true
	}
	return false
}