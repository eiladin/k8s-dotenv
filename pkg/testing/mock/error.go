package mock

import "errors"

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

// AnError is a general error for testing.
var AnError = errors.New("mock.AnError general error for testing") //nolint
