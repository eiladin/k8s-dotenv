package mock

// Error is used in tests.
type Error struct {
	message string
}

// NewError constructor.
func NewError(message string) *Error {
	return &Error{
		message: message,
	}
}

// Error implements the error interface.
func (e *Error) Error() string {
	return e.message
}

var AnError = NewError("mock.AnError general error for testing")
