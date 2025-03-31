package errors

// ForbiddenError represents a "Forbidden" error.
type ForbiddenError struct {
	*BaseError[any] // embedding BaseError for reusability
}

var (
	// The symbol to identify the error type.
	ForbiddenType = symbol("ForbiddenError")
	// Default name and message for the error.
	ForbiddenName    = "ForbiddenError"
	ForbiddenMessage = "Forbidden"
)

// NewForbiddenError creates a new ForbiddenError instance.
func NewForbiddenError[M any](message string, code string, meta M) *ForbiddenError {
	// Default to "ForbiddenError" code and message if none are provided.
	if message == "" {
		message = ForbiddenMessage
	}
	if code == "" {
		code = ForbiddenName
	}

	// Create a new BaseError with the specific type.
	baseError := NewBaseError[any](ForbiddenName, message, code, meta, ForbiddenType)
	return &ForbiddenError{
		BaseError: baseError,
	}
}

// Is checks if the error is of type ForbiddenError.
func (e *ForbiddenError) Is(target error) bool {
	// This checks if the target error is of type ForbiddenError
	if _, ok := target.(*ForbiddenError); ok {
		return true
	}
	return false
}