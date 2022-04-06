package cronjob

import (
	"errors"
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clioptions"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewCmd(t *testing.T) {
	v1mock := mock.CronJobv1("my-cronjob", "test", nil, nil, nil)
	kubeClient := mock.NewFakeClient(v1mock).WithResources(mock.CronJobv1Resource())

	t.Run("create", func(t *testing.T) {
		got := NewCmd(&clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test"})
		if got == nil {
			t.Errorf("NewCmd() is nil want not nil")
		}
	})

	t.Run("valid args", func(t *testing.T) {
		got := NewCmd(&clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test"})
		cronjobs, _ := got.ValidArgsFunction(got, []string{}, "")
		if cronjobs[0] != "my-cronjob" {
			t.Errorf("NewCmd().ValidArgs = %v, want %v", cronjobs, []string{"my-cronjob"})
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

func Test_clientError(t *testing.T) {
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
			if err := clientError(tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("clientError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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
	v1mock := mock.CronJobv1("my-cronjob", "test", nil, nil, nil)
	v1beta1mock := mock.CronJobv1beta1("my-beta-cronjob", "test", nil, nil, nil)
	kubeClient := mock.NewFakeClient(v1mock, v1beta1mock)

	type args struct {
		opt         *clioptions.CLIOptions
		apiresource *metav1.APIResourceList
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "find v1 cronjobs",
			args: args{
				opt:         &clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test"},
				apiresource: mock.CronJobv1Resource(),
			},
			want: []string{"my-cronjob"},
		},
		{
			name: "find v1beta1 cronjobs",
			args: args{
				opt:         &clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test"},
				apiresource: mock.CronJobv1beta1Resource(),
			},
			want: []string{"my-beta-cronjob"},
		},
		{
			name: "don't find non-existent groups",
			args: args{
				opt:         &clioptions.CLIOptions{KubeClient: kubeClient, Namespace: "test"},
				apiresource: mock.InvalidGroupResource(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kubeClient.Fake.Resources = []*metav1.APIResourceList{tt.args.apiresource}
			if got := validArgs(tt.args.opt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_run(t *testing.T) {
	v1mock := mock.CronJobv1("my-cronjob", "test", map[string]string{"k1": "v1", "k2": "v2"}, nil, nil)
	v1beta1mock := mock.CronJobv1beta1("my-beta-cronjob", "test", map[string]string{"k1": "v1", "k2": "v2"}, nil, nil)

	errorClient := mock.NewFakeClient().WithResources(mock.InvalidGroupResource())

	writer := mock.NewWriter()
	v1Client := mock.NewFakeClient(v1mock).WithResources(mock.CronJobv1Resource())
	v1beta1Client := mock.NewFakeClient(v1beta1mock).WithResources(mock.CronJobv1beta1Resource())
	groupClient := mock.NewFakeClient().WithResources(mock.UnsupportedGroupResource())

	type args struct {
		opt  *clioptions.CLIOptions
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "error with no args", wantErr: true},
		{
			name: "return client errors",
			args: args{
				opt:  &clioptions.CLIOptions{KubeClient: errorClient},
				args: []string{"test"},
			},
			wantErr: true,
		},
		{
			name: "write v1 cronjobs",
			args: args{
				opt:  &clioptions.CLIOptions{KubeClient: v1Client, Namespace: "test", Writer: writer},
				args: []string{"my-cronjob"},
			},
		},
		{
			name: "write v1beta1 cronjobs",
			args: args{
				opt:  &clioptions.CLIOptions{KubeClient: v1beta1Client, Namespace: "test", Writer: writer},
				args: []string{"my-beta-cronjob"},
			},
		},
		{
			name: "error on unsupported group",
			args: args{
				opt:  &clioptions.CLIOptions{KubeClient: groupClient, Namespace: "test", Writer: writer},
				args: []string{"test"},
			},
			wantErr: true,
		},
		{
			name: "return writer errors",
			args: args{
				opt:  &clioptions.CLIOptions{KubeClient: v1Client, Namespace: "test", Writer: mock.NewErrorWriter().ErrorAfter(1)},
				args: []string{"my-cronjob"},
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
