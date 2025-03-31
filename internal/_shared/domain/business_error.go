package domain

import (
	"fmt"
	"github.com/google/uuid"
)

// BaseError structure
type BaseError[M any] struct {
	Name    string
	Type    uuid.UUID
	Code    string
	Message string
	Meta    *M
}

// BusinessError is the domain-specific error for business logic violations.
type BusinessError struct {
	BaseError[any]
}

// NewBusinessError creates a new BusinessError instance.
func NewBusinessError(message, code string) *BusinessError {
	if code == "" {
		code = "BusinessError"
	}
	return &BusinessError{
		BaseError: BaseError[any]{
			Name:    "BusinessError",
			Type:    uuid.New(),
			Code:    code,
			Message: message,
		},
	}
}

// Error implements the error interface for BusinessError
func (e *BusinessError) Error() string {
	return fmt.Sprintf("BusinessError - Code: %s, Message: %s", e.Code, e.Message)
}

// Is checks if the error is of type BusinessError.
func (e *BusinessError) Is(target error) bool {
	// Check if the target error is a BusinessError
	if be, ok := target.(*BusinessError); ok {
		return be.Code == e.Code
	}
	return false
}

// Create a new instance of BusinessError with a custom message and code
func CreateBusinessError(message, code string) *BusinessError {
	return NewBusinessError(message, code)
}