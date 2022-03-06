package client

import (
	"bytes"
	"os"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestClientWrite(t *testing.T) {
	type testCase struct {
		Name           string
		Client         *Client
		Result         *Result
		Reader         func() string
		ExpectedResult string
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Client.result = tc.Result
			actualError := tc.Client.Write()

			if tc.Reader != nil {
				assert.Equal(t, tc.ExpectedResult, tc.Reader())
			}

			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	objs := []runtime.Object{}
	objs = append(objs, mock.ConfigMap("test", "test", map[string]string{"cm1": "val", "cm2": "val2"}))
	objs = append(objs, mock.Secret("test", "test", map[string][]byte{"sec1": []byte("val"), "sec2": []byte("val2")}))
	kubeClient := mock.NewFakeClient(objs...)
	envMap := map[string]string{"env1": "val", "env2": "val2"}

	r1 := resultFromContainers([]v1.Container{mock.Container(envMap, []string{"test"}, []string{"test"})})
	r2 := resultFromContainers([]v1.Container{mock.Container(envMap, []string{"test"}, []string{"test1"})})

	envResult := "env1=\"val\"\nenv2=\"val2\"\n"
	secResult := "##### SECRET - test #####\nsec1=\"val\"\nsec2=\"val2\"\n"
	cmResult := "##### CONFIGMAP - test #####\ncm1=\"val\"\ncm2=\"val2\"\n"

	var b bytes.Buffer

	validate(t, &testCase{
		Name:           "Should work",
		Client:         NewClient(kubeClient, WithNamespace("test"), WithWriter(&b)),
		Result:         r1,
		Reader:         b.String,
		ExpectedResult: envResult + secResult + cmResult,
	})

	var b2 bytes.Buffer

	validate(t, &testCase{
		Name:        "Should Error with missing secret",
		Result:      r2,
		Client:      NewClient(kubeClient, WithNamespace("test"), WithWriter(&b)),
		Reader:      b2.String,
		ExpectError: true,
	})

	defer os.Remove("./test.out")
	validate(t, &testCase{
		Name:        "Should Error with missing writer",
		Result:      r1,
		Client:      NewClient(kubeClient, WithNamespace("test")),
		ExpectError: true,
	})

	validate(t, &testCase{
		Name:        "Should return writer errors",
		Result:      r1,
		Client:      NewClient(kubeClient, WithNamespace("test"), WithWriter(mock.NewErrorWriter().ErrorAfter(1))),
		ExpectError: true,
	})

	cl := NewClient(kubeClient)
	cl.Error = assert.AnError

	validate(t, &testCase{
		Name:        "Should return client error",
		Result:      r1,
		Client:      cl,
		ExpectError: true,
	})
}

func TestResultOutput(t *testing.T) {
	type testCase struct {
		Name           string
		Result         *Result
		Client         *Client
		ShouldExport   bool
		ExpectedString string
		ErrorChecker   func(err error) bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := tc.Result.output(tc.Client)

			assert.Equal(t, tc.ExpectedString, actualString)
			if tc.ErrorChecker != nil {
				assert.Equal(t, true, tc.ErrorChecker(actualError))
			}
		})
	}

	kubeClient := mock.NewFakeClient(
		mock.ConfigMap("test", "test", map[string]string{"cm1": "val", "cm2": "val2"}),
		mock.Secret("test", "test", map[string][]byte{"sec1": []byte("val"), "sec2": []byte("val2")}),
	)

	envMap := map[string]string{"env1": "val", "env2": "val2"}

	envResult := "env1=\"val\"\nenv2=\"val2\"\n"
	secResult := "##### SECRET - test #####\nsec1=\"val\"\nsec2=\"val2\"\n"
	cmResult := "##### CONFIGMAP - test #####\ncm1=\"val\"\ncm2=\"val2\"\n"

	cl := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:           "Should get env configmaps and secrets",
		Result:         &Result{Environment: envMap, ConfigMaps: []string{"test"}, Secrets: []string{"test"}},
		Client:         cl,
		ExpectedString: envResult + secResult + cmResult,
	})

	validate(t, &testCase{
		Name:           "Should get env and configmaps with no secrets",
		Result:         &Result{Environment: envMap, ConfigMaps: []string{"test"}},
		Client:         cl,
		ExpectedString: envResult + cmResult,
	})

	validate(t, &testCase{
		Name:           "Should get env and secrets with no configmaps",
		Result:         &Result{Environment: envMap, Secrets: []string{"test"}},
		Client:         cl,
		ExpectedString: envResult + secResult,
	})

	validate(t, &testCase{
		Name:           "Should get env with no secrets or configmaps",
		Result:         &Result{Environment: envMap},
		Client:         cl,
		ExpectedString: envResult,
	})

	validate(t, &testCase{
		Name:           "Should get configmaps with no env or secrets",
		Result:         &Result{ConfigMaps: []string{"test"}},
		Client:         cl,
		ExpectedString: cmResult,
	})

	validate(t, &testCase{
		Name:           "Should get secrets with no env or configmaps",
		Result:         &Result{Secrets: []string{"test"}},
		Client:         cl,
		ExpectedString: secResult,
	})

	validate(t, &testCase{
		Name:   "Should error with missing secret",
		Result: &Result{Secrets: []string{"test1"}},
		Client: cl,
		ErrorChecker: func(err error) bool {
			return assert.ErrorIs(t, err, ErrMissingResource)
		},
	})

	validate(t, &testCase{
		Name:   "Should error with missing configmap",
		Result: &Result{ConfigMaps: []string{"test1"}},
		Client: cl,
		ErrorChecker: func(err error) bool {
			return assert.ErrorIs(t, err, ErrMissingResource)
		},
	})
}

func TestNewResult(t *testing.T) {
	type testCase struct {
		Name           string
		ExpectedResult *Result
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualResult := NewResult()

			assert.Equal(t, tc.ExpectedResult, actualResult)
		})
	}

	validate(t, &testCase{
		Name: "Should return a new Result",
		ExpectedResult: &Result{
			Environment: map[string]string{},
			Secrets:     []string{},
			ConfigMaps:  []string{},
		},
	})
}

func TestResultFromContainers(t *testing.T) {
	type testCase struct {
		Name           string
		Containers     []v1.Container
		ExpectedResult *Result
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualResult := resultFromContainers(tc.Containers)

			assert.Equal(t, tc.ExpectedResult, actualResult)
		})
	}

	validate(t, &testCase{
		Name: "Should return results from containers",
		Containers: []v1.Container{
			mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"configmap"}, []string{"secret"}),
		},
		ExpectedResult: &Result{
			Environment: map[string]string{"env1": "val", "env2": "val2"},
			ConfigMaps:  []string{"configmap"},
			Secrets:     []string{"secret"},
		},
	})
}
