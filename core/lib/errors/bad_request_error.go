
package errors

// BadRequestError represents a "BadRequest" error.
type BadRequestError struct {
	*BaseError[any] // embedding BaseError for reusability
}

var (
	// The symbol to identify the error type.
	BadRequestType = symbol("BadRequestError")
	// Default name and message for the error.
	BadRequestName    = "BadRequestError"
	BadRequestMessage = "Bad Request"
)

// NewBadRequestError creates a new BadRequestError instance.
func NewBadRequestError[M any](message string, code string, meta M) *BadRequestError {
	// Default to "BadRequestError" code and message if none are provided.
	if message == "" {
		message = BadRequestMessage
	}
	if code == "" {
		code = BadRequestName
	}

	// Create a new BaseError with the specific type.
	baseError := NewBaseError[any](BadRequestName, message, code, meta, BadRequestType)
	return &BadRequestError{
		BaseError: baseError,
	}
}

// Is checks if the error is of type BadRequestError.
func (e *BadRequestError) Is(target error) bool {
	// This checks if the target error is of type BadRequestError
	if _, ok := target.(*BadRequestError); ok {
		return true
	}
	return false
}