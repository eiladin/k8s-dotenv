package client

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestCoreV1ConfigMapValues(t *testing.T) {
	type testCase struct {
		Name           string
		Client         *CoreV1
		Configmap      string
		ExpectedString string
		ShouldExport   bool
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := tc.Client.ConfigMapValues(tc.Configmap, true)

			assert.Equal(t, tc.ExpectedString, actualString)
			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	cm := mock.ConfigMap("test", "test", map[string]string{"n": "v"})
	kubeClient := mock.NewFakeClient(cm)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:           "Should find test.test",
		Client:         NewCoreV1(client),
		Configmap:      "test",
		ExpectedString: "##### CONFIGMAP - test #####\nexport n=\"v\"\n",
	})

	validate(t, &testCase{
		Name:        "Should not find test.test1",
		Client:      NewCoreV1(client),
		Configmap:   "test1",
		ExpectError: true,
	})

	client = NewClient(kubeClient, WithNamespace("test2"))
	validate(t, &testCase{
		Name:        "Should not find test2.test",
		Client:      NewCoreV1(client),
		Configmap:   "test",
		ExpectError: true,
	})
}

func TestCoreV1SecretValues(t *testing.T) {
	type testCase struct {
		Name           string
		Client         *CoreV1
		Secret         string
		ExpectedString string
		ShouldExport   bool
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := tc.Client.SecretValues(tc.Secret, tc.ShouldExport)

			assert.Equal(t, tc.ExpectedString, actualString)

			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	cm := mock.Secret("test", "test", map[string][]byte{"n": []byte("v")})
	kubeClient := mock.NewFakeClient(cm)

	validate(t, &testCase{
		Name:           "Should find test.test",
		Secret:         "test",
		Client:         NewCoreV1(NewClient(kubeClient, WithNamespace("test"))),
		ExpectedString: "##### SECRET - test #####\nn=\"v\"\n",
	})

	validate(t, &testCase{
		Name:        "Should not find test.test1",
		Secret:      "test1",
		Client:      NewCoreV1(NewClient(kubeClient, WithNamespace("test"))),
		ExpectError: true,
	})

	validate(t, &testCase{
		Name:        "Should not find test2.test",
		Secret:      "test",
		Client:      NewCoreV1(NewClient(kubeClient, WithNamespace("test2"))),
		ExpectError: true,
	})
}

func TestCoreV1Pod(t *testing.T) {
	type testCase struct {
		Name           string
		CoreV1         *CoreV1
		Resource       string
		ExpectedResult *Result
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualClient := tc.CoreV1.Pod(tc.Resource)

			assert.Equal(t, tc.ExpectedResult, actualClient.result)

			if tc.ExpectError {
				assert.Error(t, tc.CoreV1.client.Error)
			} else {
				assert.NoError(t, tc.CoreV1.client.Error)
			}
		})
	}

	mockv1 := mock.Pod("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:     "Should return pods",
		CoreV1:   NewCoreV1(client),
		Resource: "test",
		ExpectedResult: &Result{
			Environment: map[string]string{"k": "v"},
			Secrets:     []string{"test"},
			ConfigMaps:  []string{"test"},
		},
	})

	kubeClient.PrependReactor("get", "pods", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:        "Should return API errors",
		CoreV1:      NewCoreV1(client),
		Resource:    "test",
		ExpectError: true,
	})
}

func TestCoreV1Pods(t *testing.T) {
	type testCase struct {
		Name          string
		CoreV1        *CoreV1
		ExpectedSlice []string
		ExpectError   bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := tc.CoreV1.Pods()

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	mockv1 := mock.Pod("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	kubeClient := mock.NewFakeClient(mockv1)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:          "Should return pods",
		CoreV1:        NewCoreV1(client),
		ExpectedSlice: []string{"test"},
	})

	kubeClient.PrependReactor("list", "pods", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:        "Should return API errors",
		CoreV1:      NewCoreV1(client),
		ExpectError: true,
	})
}
