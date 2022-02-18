package completion

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type CompletionCmdSuite struct {
	suite.Suite
}

func newErrorWriter(w io.Writer) io.Writer {
	return &errorWriter{w}
}

type errorWriter struct {
	w io.Writer
}

func (w *errorWriter) Write(p []byte) (int, error) {
	return 0, errors.New("error")
}

func (suite CompletionCmdSuite) TestNewCmd() {
	got := NewCmd()
	suite.NotNil(got)
}

func (suite CompletionCmdSuite) TestRunCompletion() {
	cases := []struct {
		args      []string
		shouldErr bool
	}{
		{args: []string{"zsh"}, shouldErr: false},
		{args: []string{}, shouldErr: true},
		{args: []string{"zsh", "bash"}, shouldErr: true},
		{args: []string{"not-a-shell"}, shouldErr: true},
	}

	for _, c := range cases {
		cmd := NewCmd()
		testCmd := &cobra.Command{Use: "test"}
		testCmd.AddCommand(cmd)
		var b bytes.Buffer
		err := RunCompletion(&b, cmd, c.args)
		if c.shouldErr {
			suite.Error(err)
		} else {
			suite.NoError(err)
			suite.NotEmpty(b.String())
		}
	}
}

func (suite CompletionCmdSuite) TestCompletionShells() {
	for sh, run := range completionShells {
		testCmd := &cobra.Command{Use: "test"}
		var b bytes.Buffer
		err := run(&b, testCmd)
		suite.NoError(err)
		suite.Contains(b.String(), sh)

		var errB bytes.Buffer
		errW := newErrorWriter(&errB)
		err = run(errW, testCmd)
		suite.Error(err)
	}
}

func TestCompletionCmdSuite(t *testing.T) {
	suite.Run(t, new(CompletionCmdSuite))
}