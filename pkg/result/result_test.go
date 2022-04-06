package result

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// func TestResultWrite(t *testing.T) {
// 	type testCase struct {
// 		Name           string
// 		Result         *Result
// 		Reader         func() string
// 		ExpectedResult string
// 		ExpectError    bool
// 		Writer         io.Writer
// 	}

// 	validate := func(t *testing.T, testCase *testCase) {
// 		t.Run(testCase.Name, func(t *testing.T) {
// 			actualError := testCase.Result.Write(testCase.Writer)

// 			if testCase.Reader != nil {
// 				assert.Equal(t, testCase.ExpectedResult, testCase.Reader())
// 			}

// 			if testCase.ExpectError {
// 				assert.Error(t, actualError)
// 			} else {
// 				assert.NoError(t, actualError)
// 			}
// 		})
// 	}

// 	result1 := &Result{
// 		Environment: EnvValues{"env1": "val", "env2": "val2"},
// 		ConfigMaps:  map[string]EnvValues{"test": {"cm1": "val", "cm2": "val2"}},
// 		Secrets:     map[string]EnvValues{"test": {"sec1": "val", "sec2": "val2"}},
// 	}

// 	envResult := "env1=\"val\"\nenv2=\"val2\"\n"
// 	secResult := "##### SECRET - test #####\nsec1=\"val\"\nsec2=\"val2\"\n"
// 	cmResult := "##### CONFIGMAP - test #####\ncm1=\"val\"\ncm2=\"val2\"\n"

// 	var buffer bytes.Buffer

// 	validate(t, &testCase{
// 		Name:           "Should write results",
// 		Result:         result1,
// 		Reader:         buffer.String,
// 		Writer:         &buffer,
// 		ExpectedResult: envResult + cmResult + secResult,
// 	})

// 	defer os.Remove("./test.out")
// 	validate(t, &testCase{
// 		Name:        "Should Error with missing writer",
// 		Result:      result1,
// 		ExpectError: true,
// 	})

// 	validate(t, &testCase{
// 		Name:        "Should return writer errors",
// 		Result:      result1,
// 		Writer:      mock.NewErrorWriter().ErrorAfter(1),
// 		ExpectError: true,
// 	})

// 	validate(t, &testCase{
// 		Name:        "Should return client error",
// 		Result:      NewFromError(assert.AnError),
// 		ExpectError: true,
// 	})
// }

// func TestResultParse(t *testing.T) {
// 	type testCase struct {
// 		Name           string
// 		Result         *Result
// 		ShouldExport   bool
// 		ExpectedString string
// 	}

// 	validate := func(t *testing.T, tc *testCase) {
// 		t.Run(tc.Name, func(t *testing.T) {
// 			actualString := tc.Result.parse()

// 			assert.Equal(t, tc.ExpectedString, actualString)
// 		})
// 	}

// 	envMap := EnvValues{"env1": "val", "env2": "val2"}
// 	cmMap := map[string]EnvValues{"test": {"cm1": "val", "cm2": "val2"}}
// 	secretMap := map[string]EnvValues{"test": {"sec1": "val", "sec2": "val2"}}

// 	envResult := "env1=\"val\"\nenv2=\"val2\"\n"
// 	secResult := "##### SECRET - test #####\nsec1=\"val\"\nsec2=\"val2\"\n"
// 	cmResult := "##### CONFIGMAP - test #####\ncm1=\"val\"\ncm2=\"val2\"\n"

// 	validate(t, &testCase{
// 		Name:           "Should get env configmaps and secrets",
// 		Result:         &Result{Environment: envMap, ConfigMaps: cmMap, Secrets: secretMap},
// 		ExpectedString: envResult + cmResult + secResult,
// 	})

// 	validate(t, &testCase{
// 		Name:           "Should get env and configmaps with no secrets",
// 		Result:         &Result{Environment: envMap, ConfigMaps: cmMap},
// 		ExpectedString: envResult + cmResult,
// 	})

// 	validate(t, &testCase{
// 		Name:           "Should get env and secrets with no configmaps",
// 		Result:         &Result{Environment: envMap, Secrets: secretMap},
// 		ExpectedString: envResult + secResult,
// 	})

// 	validate(t, &testCase{
// 		Name:           "Should get env with no secrets or configmaps",
// 		Result:         &Result{Environment: envMap},
// 		ExpectedString: envResult,
// 	})

// 	validate(t, &testCase{
// 		Name:           "Should get configmaps with no env or secrets",
// 		Result:         &Result{ConfigMaps: cmMap},
// 		ExpectedString: cmResult,
// 	})

// 	validate(t, &testCase{
// 		Name:           "Should get secrets with no env or configmaps",
// 		Result:         &Result{Secrets: secretMap},
// 		ExpectedString: secResult,
// 	})
// }

// func TestNewResult(t *testing.T) {
// 	type testCase struct {
// 		Name           string
// 		ExpectedResult *Result
// 	}

// 	validate := func(t *testing.T, tc *testCase) {
// 		t.Run(tc.Name, func(t *testing.T) {
// 			actualResult := newResult()

// 			assert.Equal(t, tc.ExpectedResult, actualResult)
// 		})
// 	}

// 	validate(t, &testCase{
// 		Name: "Should return a new Result",
// 		ExpectedResult: &Result{
// 			Environment: EnvValues{},
// 			Secrets:     map[string]EnvValues{},
// 			ConfigMaps:  map[string]EnvValues{},
// 		},
// 	})
// }

// func TestNewFromContainers(t *testing.T) {
// 	type testCase struct {
// 		Name           string
// 		Client         kubernetes.Interface
// 		Namespace      string
// 		ShouldExport   bool
// 		Containers     []v1.Container
// 		ExpectedResult *Result
// 		ExpectedError  error
// 	}

// 	validate := func(t *testing.T, testCase *testCase) {
// 		t.Run(testCase.Name, func(t *testing.T) {
// 			actualResult := NewFromContainers(testCase.Client, testCase.Namespace, testCase.ShouldExport, testCase.Containers)

// 			assert.Equal(t, testCase.ExpectedResult, actualResult)

// 			assert.ErrorIs(t, actualResult.Error, testCase.ExpectedError)
// 		})
// 	}

// 	kubeClient := mock.NewFakeClient(
// 		mock.ConfigMap("test", "test", map[string]string{"cm1": "val", "cm2": "val2"}),
// 		mock.Secret("test", "test", map[string][]byte{"sec1": []byte("val"), "sec2": []byte("val2")}),
// 	)

// 	validate(t, &testCase{
// 		Name:      "Should return results from containers",
// 		Client:    kubeClient,
// 		Namespace: "test",
// 		Containers: []v1.Container{
// 			mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test"}, []string{"test"}),
// 		},
// 		ExpectedResult: &Result{
// 			shouldExport: false,
// 			Environment:  EnvValues{"env1": "val", "env2": "val2"},
// 			ConfigMaps:   map[string]EnvValues{"test": {"cm1": "val", "cm2": "val2"}},
// 			Secrets:      map[string]EnvValues{"test": {"sec1": "val", "sec2": "val2"}},
// 		},
// 	})

// 	validate(t, &testCase{
// 		Name:   "Should set Error with missing secret",
// 		Client: kubeClient,
// 		Containers: []v1.Container{
// 			mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test"}, []string{"test1"}),
// 		},
// 		Namespace:      "test",
// 		ExpectedResult: NewFromError(ErrMissingResource),
// 		ExpectedError:  ErrMissingResource,
// 	})

// 	validate(t, &testCase{
// 		Name:   "Should set Error with missing configmap",
// 		Client: kubeClient,
// 		Containers: []v1.Container{
// 			mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test1"}, []string{"test"}),
// 		},
// 		Namespace:      "test",
// 		ExpectedResult: NewFromError(ErrMissingResource),
// 		ExpectedError:  ErrMissingResource,
// 	})
// }

func Test_newWriteError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "wraps error", args: args{err: assert.AnError}, wantErr: true},
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := configMapData(tt.args.client, tt.args.namespace, tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("configMapData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("configMapData() = %v, want %v", got, tt.want)
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := secretData(tt.args.client, tt.args.namespace, tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("secretData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("secretData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFromError(t *testing.T) {
	err := assert.AnError
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
		containers   []corev1.Container
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				got := NewFromContainers(tt.args.client, tt.args.namespace, tt.args.shouldExport, tt.args.containers)
				if got.Error.Error() != tt.want.Error.Error() {
					t.Errorf("NewFromContainers() = %v want %v", got.Error, tt.want.Error)
				}
			} else {
				if got := NewFromContainers(tt.args.client, tt.args.namespace, tt.args.shouldExport, tt.args.containers); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewFromContainers() = %v, want %v", got, tt.want)
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
		wantErr    bool
		writer     io.Writer
		nilWriter  bool
	}{
		{
			name:    "return error",
			r:       &Result{Error: assert.AnError},
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.nilWriter {
				if err := tt.r.Write(nil); (err != nil) != tt.wantErr {
					t.Errorf("Result.Write() err = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else if tt.writer == nil {
				writer := &bytes.Buffer{}
				if err := tt.r.Write(writer); (err != nil) != tt.wantErr {
					t.Errorf("Result.Write() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotWriter := writer.String(); gotWriter != tt.wantWriter {
					t.Errorf("Result.Write() = %v, want %v", gotWriter, tt.wantWriter)
					return
				}
			} else {
				if err := tt.r.Write(tt.writer); (err != nil) != tt.wantErr {
					t.Errorf("Result.Write() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}
