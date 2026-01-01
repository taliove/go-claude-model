package errors

// Error represents an error with a suggestion for the user
type Error struct {
	Message    string
	Suggestion string
}

// NewError creates a new error with a suggestion
func NewError(message, suggestion string) *Error {
	return &Error{
		Message:    message,
		Suggestion: suggestion,
	}
}

// Error returns the error message
func (e *Error) Error() string {
	return e.Message
}
