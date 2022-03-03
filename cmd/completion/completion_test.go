package completion

import (
	"bytes"
	"os"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	tests "github.com/eiladin/k8s-dotenv/pkg/testing"
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

	validate(t, &testCase{
		Name: "Should run",
		Cmd:  NewCmd(&options.Options{Writer: &b}),
		Args: []string{"zsh"},
	})

	validate(t, &testCase{
		Name: "Should error with too many shell types",
		Cmd:  NewCmd(&options.Options{Writer: &b}),
		Args: []string{"zsh", "bash"}, ExpectedError: ErrTooManyArguments,
	})

	validate(t, &testCase{
		Name:          "Should error with no arguments",
		Cmd:           NewCmd(&options.Options{Writer: &b}),
		ExpectedError: ErrShellNotSpecified,
	})

	validate(t, &testCase{
		Name:          "Should error with unsupported shell type",
		Cmd:           NewCmd(&options.Options{Writer: &b}),
		Args:          []string{"not-a-shell"},
		ExpectedError: ErrUnsupportedShell,
	})
}

func TestCompletionShells(t *testing.T) {
	for sh, run := range completionShells() {
		var b bytes.Buffer

		var errB bytes.Buffer

		testCmd := &cobra.Command{Use: "test"}
		err := run(&b, testCmd)
		assert.NoError(t, err)
		assert.Contains(t, b.String(), sh)

		errorAfter := 0

		if sh == "zsh" {
			errorAfter = 1
		}

		errW := tests.NewErrorWriter(&errB).ErrorAfter(errorAfter)
		err = run(errW, testCmd)
		assert.Error(t, err)

		errB.Reset()

		errorAfter += 2

		errW2 := tests.NewErrorWriter(&errB).ErrorAfter(errorAfter)
		err = run(errW2, testCmd)
		assert.Error(t, err)
	}
}
