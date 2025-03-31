package errors

// NotFoundError represents a "Not Found" error.
type NotFoundError struct {
	*BaseError[any] // embedding BaseError for reusability
}

var (
	// The symbol to identify the error type.
	NotFoundType = symbol("NotFoundError")
	// Default name and message for the error.
	NotFoundName    = "NotFoundError"
	NotFoundMessage = "Not Found"
)

// NewNotFoundError creates a new NotFoundError instance.
func NewNotFoundError[M any](message string, code string, meta M) *NotFoundError {
	// Default to "NotFoundError" code and message if none are provided.
	if message == "" {
		message = NotFoundMessage
	}
	if code == "" {
		code = NotFoundName
	}

	// Create a new BaseError with the specific type.
	baseError := NewBaseError[any](NotFoundName, message, code, meta, NotFoundType)
	return &NotFoundError{
		BaseError: baseError,
	}
}

// Is checks if the error is of type NotFoundError.
func (e *NotFoundError) Is(target error) bool {
	// This checks if the target error is of type NotFoundError
	if _, ok := target.(*NotFoundError); ok {
		return true
	}
	return false
}