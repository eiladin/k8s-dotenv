package client

import (
	"bytes"
	"os"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestResultWrite(t *testing.T) {
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

	kubeClient := mock.NewFakeClient()

	r1 := &Result{
		Environment: envValues{"env1": "val", "env2": "val2"},
		ConfigMaps:  map[string]envValues{"test": {"cm1": "val", "cm2": "val2"}},
		Secrets:     map[string]envValues{"test": {"sec1": "val", "sec2": "val2"}},
	}

	envResult := "env1=\"val\"\nenv2=\"val2\"\n"
	secResult := "##### SECRET - test #####\nsec1=\"val\"\nsec2=\"val2\"\n"
	cmResult := "##### CONFIGMAP - test #####\ncm1=\"val\"\ncm2=\"val2\"\n"

	var b bytes.Buffer

	validate(t, &testCase{
		Name:           "Should write results",
		Client:         NewClient(kubeClient, WithNamespace("test"), WithWriter(&b)),
		Result:         r1,
		Reader:         b.String,
		ExpectedResult: envResult + cmResult + secResult,
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

	client := NewClient(kubeClient)
	client.Error = assert.AnError

	validate(t, &testCase{
		Name:        "Should return client error",
		Result:      r1,
		Client:      client,
		ExpectError: true,
	})
}

func TestResultParse(t *testing.T) {
	type testCase struct {
		Name           string
		Result         *Result
		Client         *Client
		ShouldExport   bool
		ExpectedString string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString := tc.Result.parse(tc.Client)

			assert.Equal(t, tc.ExpectedString, actualString)
		})
	}

	envMap := envValues{"env1": "val", "env2": "val2"}
	cmMap := map[string]envValues{"test": {"cm1": "val", "cm2": "val2"}}
	secretMap := map[string]envValues{"test": {"sec1": "val", "sec2": "val2"}}

	envResult := "env1=\"val\"\nenv2=\"val2\"\n"
	secResult := "##### SECRET - test #####\nsec1=\"val\"\nsec2=\"val2\"\n"
	cmResult := "##### CONFIGMAP - test #####\ncm1=\"val\"\ncm2=\"val2\"\n"

	client := NewClient(mock.NewFakeClient(), WithNamespace("test"))

	validate(t, &testCase{
		Name:           "Should get env configmaps and secrets",
		Result:         &Result{Environment: envMap, ConfigMaps: cmMap, Secrets: secretMap},
		Client:         client,
		ExpectedString: envResult + cmResult + secResult,
	})

	validate(t, &testCase{
		Name:           "Should get env and configmaps with no secrets",
		Result:         &Result{Environment: envMap, ConfigMaps: cmMap},
		Client:         client,
		ExpectedString: envResult + cmResult,
	})

	validate(t, &testCase{
		Name:           "Should get env and secrets with no configmaps",
		Result:         &Result{Environment: envMap, Secrets: secretMap},
		Client:         client,
		ExpectedString: envResult + secResult,
	})

	validate(t, &testCase{
		Name:           "Should get env with no secrets or configmaps",
		Result:         &Result{Environment: envMap},
		Client:         client,
		ExpectedString: envResult,
	})

	validate(t, &testCase{
		Name:           "Should get configmaps with no env or secrets",
		Result:         &Result{ConfigMaps: cmMap},
		Client:         client,
		ExpectedString: cmResult,
	})

	validate(t, &testCase{
		Name:           "Should get secrets with no env or configmaps",
		Result:         &Result{Secrets: secretMap},
		Client:         client,
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
			actualResult := NewResult()

			assert.Equal(t, tc.ExpectedResult, actualResult)
		})
	}

	validate(t, &testCase{
		Name: "Should return a new Result",
		ExpectedResult: &Result{
			Environment: envValues{},
			Secrets:     map[string]envValues{},
			ConfigMaps:  map[string]envValues{},
		},
	})
}

func TestResultFromContainers(t *testing.T) {
	type testCase struct {
		Name           string
		Client         *Client
		Containers     []v1.Container
		ExpectedResult *Result
		ExpectedError  error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualResult := tc.Client.resultFromContainers(tc.Containers)

			assert.Equal(t, tc.ExpectedResult, actualResult.result)

			assert.ErrorIs(t, tc.Client.Error, tc.ExpectedError)
		})
	}

	kubeClient := mock.NewFakeClient(
		mock.ConfigMap("test", "test", map[string]string{"cm1": "val", "cm2": "val2"}),
		mock.Secret("test", "test", map[string][]byte{"sec1": []byte("val"), "sec2": []byte("val2")}),
	)

	validate(t, &testCase{
		Name:   "Should return results from containers",
		Client: NewClient(kubeClient, WithNamespace("test")),
		Containers: []v1.Container{
			mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test"}, []string{"test"}),
		},
		ExpectedResult: &Result{
			Environment: envValues{"env1": "val", "env2": "val2"},
			ConfigMaps:  map[string]envValues{"test": {"cm1": "val", "cm2": "val2"}},
			Secrets:     map[string]envValues{"test": {"sec1": "val", "sec2": "val2"}},
		},
	})

	validate(t, &testCase{
		Name:   "Should set client.Error with missing secret",
		Client: NewClient(kubeClient, WithNamespace("test")),
		Containers: []v1.Container{
			mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test"}, []string{"test1"}),
		},
		ExpectedError: ErrMissingResource,
	})

	validate(t, &testCase{
		Name:   "Should set client.Error with missing configmap",
		Client: NewClient(kubeClient, WithNamespace("test")),
		Containers: []v1.Container{
			mock.Container(map[string]string{"env1": "val", "env2": "val2"}, []string{"test1"}, []string{"test"}),
		},
		ExpectedError: ErrMissingResource,
	})
}
