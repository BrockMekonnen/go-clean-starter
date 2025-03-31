package errors

// UnauthorizedError represents an "Unauthorized" error.
type UnauthorizedError struct {
	*BaseError[any] // embedding BaseError for reusability
}

var (
	// The symbol to identify the error type.
	UnauthorizedType = symbol("UnauthorizedError")
	// Default name and message for the error.
	UnauthorizedName    = "UnauthorizedError"
	UnauthorizedMessage = "Unauthorized"
)

// NewUnauthorizedError creates a new UnauthorizedError instance.
func NewUnauthorizedError[M any](message string, code string, meta M) *UnauthorizedError {
	// Default to "UnauthorizedError" code and message if none are provided.
	if message == "" {
		message = UnauthorizedMessage
	}
	if code == "" {
		code = UnauthorizedName
	}

	// Create a new BaseError with the specific type.
	baseError := NewBaseError[any](UnauthorizedName, message, code, meta, UnauthorizedType)
	return &UnauthorizedError{
		BaseError: baseError,
	}
}

// Is checks if the error is of type UnauthorizedError.
func (e *UnauthorizedError) Is(target error) bool {
	// This checks if the target error is of type UnauthorizedError
	if _, ok := target.(*UnauthorizedError); ok {
		return true
	}
	return false
}