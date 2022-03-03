package testing

import (
	"io"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
)

// NewErrorWriter creates an io.Writer used in tests.
func NewErrorWriter(w io.Writer) *ErrorWriter {
	return &ErrorWriter{w: w}
}

// ErrorWriter used in tests.
type ErrorWriter struct {
	writes     int
	errorAfter int
	w          io.Writer
}

// ErrorAfter will cause the writer to error after a certain number of writes.
func (w *ErrorWriter) ErrorAfter(errorAfter int) io.Writer {
	w.errorAfter = errorAfter

	return w
}

// Write will return 0 for length and after `errorAfter` is hit, will start returning mock.NewError("error").
func (w *ErrorWriter) Write(p []byte) (int, error) {
	w.writes++
	if w.writes >= w.errorAfter {
		return 0, mock.NewError("error")
	}

	return 0, nil
}
