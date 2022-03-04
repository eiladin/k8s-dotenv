package completion

import (
	"os"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
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

	opt = &options.Options{Writer: mock.NewWriter()}
	cmd = NewCmd(opt)
	cmd.PreRun(cmd, []string{})
	assert.Equal(t, mock.NewWriter(), opt.Writer)
}

func TestRun(t *testing.T) {
	type testCase struct {
		Name        string
		Cmd         *cobra.Command
		Args        []string
		ExpectError bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			parentCmd := &cobra.Command{Use: "test"}
			parentCmd.AddCommand(tc.Cmd)
			actualError := tc.Cmd.RunE(tc.Cmd, tc.Args)

			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	validate(t, &testCase{
		Name: "Should run",
		Cmd:  NewCmd(&options.Options{Writer: mock.NewWriter()}),
		Args: []string{"zsh"},
	})

	validate(t, &testCase{
		Name:        "Should error with too many shell types",
		Cmd:         NewCmd(&options.Options{Writer: mock.NewWriter()}),
		Args:        []string{"zsh", "bash"},
		ExpectError: true,
	})

	validate(t, &testCase{
		Name:        "Should error with no arguments",
		Cmd:         NewCmd(&options.Options{Writer: mock.NewWriter()}),
		ExpectError: true,
	})

	validate(t, &testCase{
		Name:        "Should error with unsupported shell type",
		Cmd:         NewCmd(&options.Options{Writer: mock.NewWriter()}),
		Args:        []string{"not-a-shell"},
		ExpectError: true,
	})
}

func TestCompletionShells(t *testing.T) {
	for sh, run := range completionShells() {
		testCmd := &cobra.Command{Use: "test"}
		wr := mock.NewWriter()
		err := run(wr, testCmd)
		assert.NoError(t, err)
		assert.Contains(t, wr.String(), sh)

		errorAfter := 0

		if sh == "zsh" {
			errorAfter = 1
		}

		errW := mock.NewErrorWriter().ErrorAfter(errorAfter)
		err = run(errW, testCmd)
		assert.Error(t, err)

		errorAfter += 2
		errW2 := mock.NewErrorWriter().ErrorAfter(errorAfter)
		err = run(errW2, testCmd)
		assert.Error(t, err)
	}
}
