package result

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func TestResultWrite(t *testing.T) {
	type testCase struct {
		Name           string
		Result         *Result
		Reader         func() string
		ExpectedResult string
		ExpectError    bool
		Writer         io.Writer
	}

	validate := func(t *testing.T, testCase *testCase) {
		t.Run(testCase.Name, func(t *testing.T) {
			actualError := testCase.Result.Write(testCase.Writer)

			if testCase.Reader != nil {
				assert.Equal(t, testCase.ExpectedResult, testCase.Reader())
			}

			if testCase.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	result1 := &Result{
		Environment: EnvValues{"env1": "val", "env2": "val2"},
		ConfigMaps:  map[string]EnvValues{"test": {"cm1": "val", "cm2": "val2"}},
		Secrets:     map[string]EnvValues{"test": {"sec1": "val", "sec2": "val2"}},
	}

	envResult := "env1=\"val\"\nenv2=\"val2\"\n"
	secResult := "##### SECRET - test #####\nsec1=\"val\"\nsec2=\"val2\"\n"
	cmResult := "##### CONFIGMAP - test #####\ncm1=\"val\"\ncm2=\"val2\"\n"

	var buffer bytes.Buffer

	validate(t, &testCase{
		Name:           "Should write results",
		Result:         result1,
		Reader:         buffer.String,
		Writer:         &buffer,
		ExpectedResult: envResult + cmResult + secResult,
	})

	defer os.Remove("./test.out")
	validate(t, &testCase{
		Name:        "Should Error with missing writer",
		Result:      result1,
		ExpectError: true,
	})

	validate(t, &testCase{
		Name:        "Should return writer errors",
		Result:      result1,
		Writer:      mock.NewErrorWriter().ErrorAfter(1),
		ExpectError: true,
	})

	validate(t, &testCase{
		Name:        "Should return client error",
		Result:      NewFromError(assert.AnError),
		ExpectError: true,
	})
}

func TestResultParse(t *testing.T) {
	type testCase struct {
		Name           string
		Result         *Result
		ShouldExport   bool
		ExpectedString string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString := tc.Result.parse()

			assert.Equal(t, tc.ExpectedString, actualString)
		})
	}

	envMap := EnvValues{"env1": "val", "env2": "val2"}
	cmMap := map[string]EnvValues{"test": {"cm1": "val", "cm2": "val2"}}
	secretMap := map[string]EnvValues{"test": {"sec1": "val", "sec2": "val2"}}

	envResult := "env1=\"val\"\nenv2=\"val2\"\n"
	secResult := "##### SECRET - test #####\nsec1=\"val\"\nsec2=\"val2\"\n"
	cmResult := "##### CONFIGMAP - test #####\ncm1=\"val\"\ncm2=\"val2\"\n"

	validate(t, &testCase{
		Name:           "Should get env configmaps and secrets",
		Result:         &Result{Environment: envMap, ConfigMaps: cmMap, Secrets: secretMap},
		ExpectedString: envResult + cmResult + secResult,
	})

	validate(t, &testCase{
		Name:           "Should get env and configmaps with no secrets",
		Result:         &Result{Environment: envMap, ConfigMaps: cmMap},
		ExpectedString: envResult + cmResult,
	})

	validate(t, &testCase{
		Name:           "Should get env and secrets with no configmaps",
		Result:         &Result{Environment: envMap, Secrets: secretMap},
		ExpectedString: envResult + secResult,
	})

	validate(t, &testCase{
		Name:           "Should get env with no secrets or configmaps",
		Result:         &Result{Environment: envMap},
		ExpectedString: envResult,
	})

	validate(t, &testCase{
		Name:           "Should get configmaps with no env or secrets",
		Result:         &Result{ConfigMaps: cmMap},
		ExpectedString: cmResult,
	})

	validate(t, &testCase{
		Name:           "Should get secrets with no env or configmaps",
		Result:         &Result{Secrets: secretMap},
		ExpectedString: secResult,
	})
}

func TestNewResult(t *testing.T) {
	type testCase struct {
		Name           string
		ExpectedResult *Result
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualResult := newResult()

			assert.Equal(t, tc.ExpectedResult, actualResult)
		})
	}

	validate(t, &testCase{
		Name: "Should return a new Result",
		ExpectedResult: &Result{
			Environment: EnvValues{},
			Secrets:     map[string]EnvValues{},
			ConfigMaps:  map[string]EnvValues{},
		},
	})
}

func TestNewFromContainers(t *testing.T) {
	type testCase struct {
		Name           string
		Client         kubernetes.Interface
		Namespace      string
		ShouldExport   bool
		Containers     []v1.Container
		ExpectedResult *Result
		ExpectedError  error
	}

	validate := func(t *testing.T, testCase *testCase) {
		t.Run(testCase.Name, func(t *testing.T) {
			actualResult := NewFromContainers(testCase.Client, testCase.Namespace, testCase.ShouldExport, testCase.Containers)

			assert.Equal(t, testCase.ExpectedResult, actualResult)

			assert.ErrorIs(t, actualResult.Error, testCase.ExpectedError)
		})
	}

	kubeClient := mock.NewFakeClient(
		mock.ConfigMap("test", "test", map[string]string{"cm1": "val", "cm2": "val2"}),
		mock.Secret("test", "test", map[string][]byte{"sec1": []byte("val"), "sec2": []byte("val2")}),
	)

	validate(t, &testCase{
		Name:      "Should return results from containers",
		Client:    kubeClient,
		Namespace: "test",
		Containers: []v1.Container{
			mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test"}, []string{"test"}),
		},
		ExpectedResult: &Result{
			shouldExport: false,
			Environment:  EnvValues{"env1": "val", "env2": "val2"},
			ConfigMaps:   map[string]EnvValues{"test": {"cm1": "val", "cm2": "val2"}},
			Secrets:      map[string]EnvValues{"test": {"sec1": "val", "sec2": "val2"}},
		},
	})

	validate(t, &testCase{
		Name:   "Should set Error with missing secret",
		Client: kubeClient,
		Containers: []v1.Container{
			mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test"}, []string{"test1"}),
		},
		Namespace:      "test",
		ExpectedResult: NewFromError(ErrMissingResource),
		ExpectedError:  ErrMissingResource,
	})

	validate(t, &testCase{
		Name:   "Should set Error with missing configmap",
		Client: kubeClient,
		Containers: []v1.Container{
			mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test1"}, []string{"test"}),
		},
		Namespace:      "test",
		ExpectedResult: NewFromError(ErrMissingResource),
		ExpectedError:  ErrMissingResource,
	})
}
