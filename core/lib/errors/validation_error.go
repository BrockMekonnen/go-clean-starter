package errors

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError represents an error due to failed validation.
type ValidationError struct {
	*BaseError[ValidationErrorProps]
}

type ValidationErrorProps struct {
	Target string                     `json:"target"`
	Error  validator.ValidationErrors `json:"error"`
}

var (
	ValidationType = symbol("ValidationError")
	ValidationName = "ValidationError"
)

func NewValidationError(message string, target string, verr validator.ValidationErrors) *ValidationError {
	var messages []string
	for _, fieldErr := range verr {
		messages = append(messages, formatFieldError(fieldErr))
	}
	combinedMessage := strings.Join(messages, ", ")

	props := ValidationErrorProps{
		Target: target,
		Error:  verr,
	}

	baseError := NewBaseError(
		ValidationName,
		combinedMessage, // Now uses clean formatted messages
		ValidationName,
		props,
		ValidationType,
	)

	return &ValidationError{
		BaseError: baseError,
	}
}

func formatFieldError(fe validator.FieldError) string {
	return fmt.Sprintf("Field validation for '%s' failed on the '%s' rule", fe.Field(), fe.Tag())
}

func (e *ValidationError) Is(target error) bool {
	_, ok := target.(*ValidationError)
	return ok
}