package deployment

import (
	"errors"
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clioptions"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
)

func TestNewCmd(t *testing.T) {
	kubeClient := mock.NewFakeClient(mock.Deployment("test", "test", nil, nil, nil))

	t.Run("create", func(t *testing.T) {
		got := NewCmd(&clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test"})
		if got == nil {
			t.Errorf("NewCmd() is nil want not nil")
		}
	})

	t.Run("valid args", func(t *testing.T) {
		got := NewCmd(&clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test"})
		resources, _ := got.ValidArgsFunction(got, []string{}, "")
		if resources[0] != "test" {
			t.Errorf("NewCmd().ValidArgs = %v, want %v", resources, []string{"test"})
		}
	})

	t.Run("runE", func(t *testing.T) {
		got := NewCmd(&clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test"})
		err := got.RunE(got, []string{})
		if !errors.Is(err, ErrResourceNameRequired) {
			t.Errorf("NewCmd().RunE = %v, want %v", err, ErrResourceNameRequired)
		}
	})
}

func Test_runError(t *testing.T) {
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
			if err := runError(tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("runError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validArgs(t *testing.T) {
	v1mock := mock.Deployment("my-deployment", "test", nil, nil, nil)
	kubeClient := mock.NewFakeClient(v1mock)

	type args struct {
		opt *clioptions.CLIOptions
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "find v1 deployments",
			args: args{
				opt: &clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test"},
			},
			want: []string{"my-deployment"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validArgs(tt.args.opt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_run(t *testing.T) {
	kubeClient := mock.NewFakeClient(mock.Deployment("test", "test", map[string]string{"k": "v", "k2": "v2"}, nil, nil))
	writer := mock.NewWriter()

	type args struct {
		opt  *clioptions.CLIOptions
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "error with no args",
			wantErr: true,
		},
		{
			name: "find deployments",
			args: args{
				opt:  &clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test", Writer: writer},
				args: []string{"test"},
			},
			wantErr: false,
		},
		{
			name: "return writer errors",
			args: args{
				opt:  &clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test", Writer: mock.NewErrorWriter().ErrorAfter(1)},
				args: []string{"test"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args.opt, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
