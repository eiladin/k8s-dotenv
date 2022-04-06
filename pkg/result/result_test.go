package result

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func Test_newWriteError(t *testing.T) {
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
			if err := newWriteError(tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("newWriteError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newResult(t *testing.T) {
	tests := []struct {
		name string
		want *Result
	}{
		{
			name: "create",
			want: &Result{
				Environment: EnvValues{},
				ConfigMaps:  map[string]EnvValues{},
				Secrets:     map[string]EnvValues{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newResult(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvValues_sortedKeys(t *testing.T) {
	tests := []struct {
		name string
		env  EnvValues
		want []string
	}{
		{name: "sort keys", env: EnvValues{"b": "v", "a": "v", "z": "v", "f": "v"}, want: []string{"a", "b", "f", "z"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.env.sortedKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnvValues.sortedKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_configMapData(t *testing.T) {
	kubeClient := mock.NewFakeClient(mock.ConfigMap("test", "test", map[string]string{"cm1": "val", "cm2": "val2"}))

	type args struct {
		client    kubernetes.Interface
		namespace string
		resource  string
	}

	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "gets config map data",
			args: args{
				client:    kubeClient,
				namespace: "test",
				resource:  "test",
			},
			want: map[string]string{"cm1": "val", "cm2": "val2"},
		},
		{
			name: "returns error when resource not found",
			args: args{
				client:    kubeClient,
				namespace: "test",
				resource:  "not-found",
			},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := configMapData(testCase.args.client, testCase.args.namespace, testCase.args.resource)
			if (err != nil) != testCase.wantErr {
				t.Errorf("configMapData() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("configMapData() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func Test_secretData(t *testing.T) {
	kubeClient := mock.NewFakeClient(
		mock.Secret("test", "test", map[string][]byte{"sec1": []byte("val"), "sec2": []byte("val2")}),
	)

	type args struct {
		client    kubernetes.Interface
		namespace string
		resource  string
	}

	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "gets secret data",
			args: args{
				client:    kubeClient,
				namespace: "test",
				resource:  "test",
			},
			want: map[string]string{"sec1": "val", "sec2": "val2"},
		},
		{
			name: "returns error when resource not found",
			args: args{
				client:    kubeClient,
				namespace: "test",
				resource:  "not-found",
			},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := secretData(testCase.args.client, testCase.args.namespace, testCase.args.resource)
			if (err != nil) != testCase.wantErr {
				t.Errorf("secretData() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("secretData() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestNewFromError(t *testing.T) {
	err := mock.AnError

	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want *Result
	}{
		{name: "create", args: args{err: err}, want: &Result{Error: err}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFromError(tt.args.err)
			if got.Error != tt.want.Error {
				t.Errorf("NewFromError() = %v, want %v", got.Error, tt.want.Error)
			}
		})
	}
}

func TestNewFromContainers(t *testing.T) {
	kubeClient := mock.NewFakeClient(
		mock.ConfigMap("test", "test", map[string]string{"cm1": "val", "cm2": "val2"}),
		mock.Secret("test", "test", map[string][]byte{"sec1": []byte("val"), "sec2": []byte("val2")}),
	)

	type args struct {
		client       kubernetes.Interface
		namespace    string
		shouldExport bool
		containers   []v1.Container
	}

	tests := []struct {
		name    string
		args    args
		want    *Result
		wantErr bool
	}{
		{
			name: "create",
			args: args{
				client:    kubeClient,
				namespace: "test",
				containers: []v1.Container{
					mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test"}, []string{"test"}),
				},
				shouldExport: true,
			},
			want: &Result{
				shouldExport: true,
				Environment:  EnvValues{"env1": "val", "env2": "val2"},
				ConfigMaps:   map[string]EnvValues{"test": {"cm1": "val", "cm2": "val2"}},
				Secrets:      map[string]EnvValues{"test": {"sec1": "val", "sec2": "val2"}},
			},
		},
		{
			name: "error on missing configmap",
			args: args{
				client:    kubeClient,
				namespace: "test",
				containers: []v1.Container{
					mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test1"}, []string{"test"}),
				},
				shouldExport: true,
			},
			want:    &Result{Error: ErrMissingResource},
			wantErr: true,
		},
		{
			name: "error on missing secret",
			args: args{
				client:    kubeClient,
				namespace: "test",
				containers: []v1.Container{
					mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test"}, []string{"test1"}),
				},
				shouldExport: true,
			},
			want:    &Result{Error: ErrMissingResource},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.wantErr {
				got := NewFromContainers(
					testCase.args.client,
					testCase.args.namespace,
					testCase.args.shouldExport,
					testCase.args.containers,
				)

				if got.Error.Error() != testCase.want.Error.Error() {
					t.Errorf("NewFromContainers() = %v want %v", got.Error, testCase.want.Error)
				}
			} else {
				if got := NewFromContainers(testCase.args.client, testCase.args.namespace, testCase.args.shouldExport, testCase.args.containers); !reflect.DeepEqual(got, testCase.want) {
					t.Errorf("NewFromContainers() = %v, want %v", got, testCase.want)
				}
			}
		})
	}
}

func TestResult_parse(t *testing.T) {
	tests := []struct {
		name string
		r    *Result
		want string
	}{
		{
			name: "parse",
			r: &Result{
				Environment: EnvValues{"env": "val"},
				ConfigMaps:  map[string]EnvValues{"test": {"cm": "val"}},
				Secrets:     map[string]EnvValues{"test": {"sec": "val"}},
			},
			want: `env="val"
##### CONFIGMAP - test #####
cm="val"
##### SECRET - test #####
sec="val"
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.parse(); got != tt.want {
				t.Errorf("Result.parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_Write(t *testing.T) {
	tests := []struct {
		name       string
		r          *Result
		wantWriter string
		writer     io.Writer
		wantErr    bool
		nilWriter  bool
	}{
		{
			name:    "return error",
			r:       &Result{Error: mock.AnError},
			wantErr: true,
		},
		{
			name:      "error on missing writer",
			r:         &Result{},
			wantErr:   true,
			nilWriter: true,
		},
		{
			name:    "return writer error",
			r:       &Result{},
			wantErr: true,
			writer:  mock.NewErrorWriter().ErrorAfter(1),
		},
		{
			name: "parse",
			r: &Result{
				Environment: EnvValues{"env": "val"},
				ConfigMaps:  map[string]EnvValues{"test": {"cm": "val"}},
				Secrets:     map[string]EnvValues{"test": {"sec": "val"}},
			},
			wantWriter: `env="val"
##### CONFIGMAP - test #####
cm="val"
##### SECRET - test #####
sec="val"
`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			switch {
			case testCase.nilWriter:
				if err := testCase.r.Write(nil); (err != nil) != testCase.wantErr {
					t.Errorf("Result.Write() err = %v, wantErr %v", err, testCase.wantErr)
				}
			case testCase.writer == nil:
				writer := &bytes.Buffer{}
				if err := testCase.r.Write(writer); (err != nil) != testCase.wantErr {
					t.Errorf("Result.Write() error = %v, wantErr %v", err, testCase.wantErr)

					return
				}
				if gotWriter := writer.String(); gotWriter != testCase.wantWriter {
					t.Errorf("Result.Write() = %v, want %v", gotWriter, testCase.wantWriter)
				}
			default:
				if err := testCase.r.Write(testCase.writer); (err != nil) != testCase.wantErr {
					t.Errorf("Result.Write() error = %v, wantErr %v", err, testCase.wantErr)
				}
			}
		})
	}
}
