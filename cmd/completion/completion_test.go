package completion

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewCmd(t *testing.T) {
	got := NewCmd(nil)
	assert.NotNil(t, got)
}

func TestPreRun(t *testing.T) {
	opt := &options.Options{}
	cmd := NewCmd(opt)
	cmd.PreRun(cmd, []string{})
	assert.Equal(t, os.Stdout, opt.Writer)

	var b bytes.Buffer
	opt = &options.Options{Writer: &b}
	cmd = NewCmd(opt)
	cmd.PreRun(cmd, []string{})
	assert.Equal(t, &b, opt.Writer)
}

func TestRun(t *testing.T) {
	type testCase struct {
		Name string

		Cmd  *cobra.Command
		Args []string

		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			parentCmd := &cobra.Command{Use: "test"}
			parentCmd.AddCommand(tc.Cmd)
			actualError := tc.Cmd.RunE(tc.Cmd, tc.Args)

			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	var b bytes.Buffer
	validate(t, &testCase{Name: "Should run", Cmd: NewCmd(&options.Options{Writer: &b}), Args: []string{"zsh"}})
	validate(t, &testCase{Name: "Should error with too many shell types", Cmd: NewCmd(&options.Options{Writer: &b}), Args: []string{"zsh", "bash"}, ExpectedError: ErrTooManyArguments})
	validate(t, &testCase{Name: "Should error with no arguments", Cmd: NewCmd(&options.Options{Writer: &b}), ExpectedError: ErrShellNotSpecified})
	validate(t, &testCase{Name: "Should error with unsupported shell type", Cmd: NewCmd(&options.Options{Writer: &b}), Args: []string{"not-a-shell"}, ExpectedError: fmt.Errorf(ErrUnsupportedShellType, "not-a-shell")})
}

func TestCompletionShells(t *testing.T) {
	for sh, run := range completionShells {
		testCmd := &cobra.Command{Use: "test"}
		var b bytes.Buffer
		err := run(&b, testCmd)
		assert.NoError(t, err)
		assert.Contains(t, b.String(), sh)

		var errB bytes.Buffer
		errW := newErrorWriter(&errB)
		err = run(errW, testCmd)
		assert.Error(t, err)
	}
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
