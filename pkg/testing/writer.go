package testing

import (
	"bytes"
)

// NewWriter creates an io.Writer used in tests.
func NewWriter() *Writer {
	return &Writer{}
}

// Writer used in tests.
type Writer struct {
	writer bytes.Buffer
}

func (w *Writer) String() string {
	return w.writer.String()
}

// Write bytes.
func (w *Writer) Write(p []byte) (int, error) {
	return w.writer.Write(p)
}
