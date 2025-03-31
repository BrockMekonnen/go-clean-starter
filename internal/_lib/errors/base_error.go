package errors

import (
	"fmt"
)

// Exception defines a general structure for errors.
type Exception[M any] struct {
	Name    string
	Type    symbol
	Message string
	Code    string
	Meta    M
}

// symbol is used to simulate the TypeScript symbol.
type symbol string

// BaseError is the main struct for custom errors.
type BaseError[M any] struct {
	Exception[M]
}

// NewBaseError creates a new BaseError instance.
func NewBaseError[M any](name string, message string, code string, meta M, errorType symbol) *BaseError[M] {
	return &BaseError[M]{
		Exception: Exception[M]{
			Name:    name,
			Type:    errorType,
			Message: message,
			Code:    code,
			Meta:    meta,
		},
	}
}

// Error implements the error interface for BaseError.
func (e *BaseError[M]) Error() string {
	return fmt.Sprintf("Error %s: %s", e.Code, e.Message)
}