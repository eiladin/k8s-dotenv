package testing

import (
	"bytes"
	"errors"
	"io"
)

// ErrWriter is a fixed message returned by `ErrorWriter.Write`.
var ErrWriter = errors.New("mock error")

// NewErrorWriter creates an io.Writer used in tests.
func NewErrorWriter() *ErrorWriter {
	return &ErrorWriter{}
}

// ErrorWriter used in tests.
type ErrorWriter struct {
	writes     int
	errorAfter int
	writer     bytes.Buffer
}

// ErrorAfter will cause the writer to error after a certain number of writes.
func (w *ErrorWriter) ErrorAfter(errorAfter int) io.Writer {
	w.errorAfter = errorAfter

	return w
}

// Write will return 0 for length and after `errorAfter` is hit, will start returning `ErrWriter`.
func (w *ErrorWriter) Write(p []byte) (int, error) {
	w.writer.Write(p)
	w.writes++

	if w.writes >= w.errorAfter {
		return 0, ErrWriter
	}

	return 0, nil
}
