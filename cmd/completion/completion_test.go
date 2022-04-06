package completion

import (
	"io"
	"os"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clioptions"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/spf13/cobra"
)

func TestNewCmd(t *testing.T) {
	opt := clioptions.CLIOptions{}
	got := NewCmd(&opt)

	t.Run("create", func(t *testing.T) {
		if got == nil {
			t.Errorf("NewCmd() = nil, want not nil")
		}
	})

	t.Run("prerun", func(t *testing.T) {
		got.PreRun(got, nil)
		if opt.Writer != os.Stdout {
			t.Errorf("NewCmd().PreRun.Writer = %v, want %v", opt.Writer, os.Stdout)
		}
	})

	opt = clioptions.CLIOptions{
		Writer: mock.NewWriter(),
	}
	got = NewCmd(&opt)

	t.Run("persistent prerun", func(t *testing.T) {
		got.PersistentPreRun(got, nil)
		if opt.Writer != os.Stdout {
			t.Errorf("NewCmd().PersistentPreRunE.Writer = %v, want %v", opt.Writer, os.Stdout)
		}
	})

	opt = clioptions.CLIOptions{
		Writer: mock.NewWriter(),
	}
	got = NewCmd(&opt)

	t.Run("runE", func(t *testing.T) {
		parent := cobra.Command{Use: "test"}
		parent.AddCommand(got)
		_ = got.RunE(got, []string{"zsh"})
	})
}

func Test_completionShells(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{name: "get shells", want: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := completionShells(); len(got) != tt.want {
				t.Errorf("completionShells() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newCompletionGenerationError(t *testing.T) {
	type args struct {
		err error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "wraps error", args: args{err: mock.AnError}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := newCompletionGenerationError(tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("newCompletionGenerationError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_runCompletion(t *testing.T) {
	type args struct {
		opt  *clioptions.CLIOptions
		cmd  *cobra.Command
		args []string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "run",
			args: args{
				opt:  &clioptions.CLIOptions{Writer: mock.NewWriter()},
				cmd:  NewCmd(&clioptions.CLIOptions{Writer: mock.NewWriter()}),
				args: []string{"zsh"},
			},
			wantErr: false,
		},
		{
			name: "error with too many shell types",
			args: args{
				opt:  &clioptions.CLIOptions{Writer: mock.NewWriter()},
				cmd:  NewCmd(&clioptions.CLIOptions{Writer: mock.NewWriter()}),
				args: []string{"zsh", "bash"},
			},
			wantErr: true,
		},
		{
			name: "error with no shells specified",
			args: args{
				opt: &clioptions.CLIOptions{Writer: mock.NewWriter()},
				cmd: NewCmd(&clioptions.CLIOptions{Writer: mock.NewWriter()}),
			},
			wantErr: true,
		},
		{
			name: "error with undefined shell type",
			args: args{
				opt:  &clioptions.CLIOptions{Writer: mock.NewWriter()},
				cmd:  NewCmd(&clioptions.CLIOptions{Writer: mock.NewWriter()}),
				args: []string{"not-a-shell"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parentCmd := &cobra.Command{Use: "test"}
			parentCmd.AddCommand(tt.args.cmd)
			if err := runCompletion(tt.args.opt, tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("runCompletion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_runCompletionBash(t *testing.T) {
	type args struct {
		root *cobra.Command
	}

	tests := []struct {
		name    string
		args    args
		wantOut bool
		wantErr bool
		writer  io.Writer
	}{
		{
			name:    "run",
			args:    args{root: &cobra.Command{Use: "test"}},
			writer:  mock.NewWriter(),
			wantOut: true,
		},
		{
			name:    "error",
			args:    args{root: &cobra.Command{Use: "test"}},
			writer:  mock.NewErrorWriter(),
			wantErr: true,
		},
		{
			name:    "error after 2",
			args:    args{root: &cobra.Command{Use: "test"}},
			writer:  mock.NewErrorWriter().ErrorAfter(2),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runCompletionBash(tt.writer, tt.args.root); (err != nil) != tt.wantErr {
				t.Errorf("runCompletionBash() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}

func Test_runCompletionZsh(t *testing.T) {
	type args struct {
		root *cobra.Command
	}

	tests := []struct {
		name    string
		args    args
		wantOut bool
		wantErr bool
		writer  io.Writer
	}{
		{
			name:    "run",
			args:    args{root: &cobra.Command{Use: "test"}},
			writer:  mock.NewWriter(),
			wantOut: true,
		},
		{
			name:    "error",
			args:    args{root: &cobra.Command{Use: "test"}},
			writer:  mock.NewErrorWriter(),
			wantErr: true,
		},
		{
			name:    "error after 2",
			args:    args{root: &cobra.Command{Use: "test"}},
			writer:  mock.NewErrorWriter().ErrorAfter(2),
			wantErr: true,
		},
		{
			name:    "error after 3",
			args:    args{root: &cobra.Command{Use: "test"}},
			writer:  mock.NewErrorWriter().ErrorAfter(3),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runCompletionZsh(tt.writer, tt.args.root); (err != nil) != tt.wantErr {
				t.Errorf("runCompletionZsh() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}

func Test_runCompletionFish(t *testing.T) {
	type args struct {
		root *cobra.Command
	}

	tests := []struct {
		name    string
		args    args
		wantOut bool
		wantErr bool
		writer  io.Writer
	}{
		{
			name:    "run",
			args:    args{root: &cobra.Command{Use: "test"}},
			writer:  mock.NewWriter(),
			wantOut: true,
		},
		{
			name:    "error",
			args:    args{root: &cobra.Command{Use: "test"}},
			writer:  mock.NewErrorWriter(),
			wantErr: true,
		},
		{
			name:    "error after 2",
			args:    args{root: &cobra.Command{Use: "test"}},
			writer:  mock.NewErrorWriter().ErrorAfter(2),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runCompletionFish(tt.writer, tt.args.root); (err != nil) != tt.wantErr {
				t.Errorf("runCompletionFish() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}

func Test_runCompletionPwsh(t *testing.T) {
	type args struct {
		root *cobra.Command
	}

	tests := []struct {
		name    string
		args    args
		wantOut bool
		wantErr bool
		writer  io.Writer
	}{
		{
			name:    "run",
			args:    args{root: &cobra.Command{Use: "test"}},
			writer:  mock.NewWriter(),
			wantOut: true,
		},
		{
			name:    "error",
			args:    args{root: &cobra.Command{Use: "test"}},
			writer:  mock.NewErrorWriter(),
			wantErr: true,
		},
		{
			name:    "error after 2",
			args:    args{root: &cobra.Command{Use: "test"}},
			writer:  mock.NewErrorWriter().ErrorAfter(2),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runCompletionPwsh(tt.writer, tt.args.root); (err != nil) != tt.wantErr {
				t.Errorf("runCompletionPwsh() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}
