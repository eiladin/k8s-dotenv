package client

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestCoreV1ConfigMapData(t *testing.T) {
	type testCase struct {
		Name           string
		Client         *CoreV1
		Configmap      string
		ExpectedValues map[string]string
		ShouldExport   bool
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualValues, actualError := tc.Client.ConfigMapData(tc.Configmap, true)

			assert.Equal(t, tc.ExpectedValues, actualValues)
			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	cm := mock.ConfigMap("test", "test", map[string]string{"n": "v"})
	kubeClient := mock.NewFakeClient(cm)
	client := NewClient(
		WithKubeClient(kubeClient),
		WithNamespace("test"),
	)

	validate(t, &testCase{
		Name:           "Should find test.test",
		Client:         NewCoreV1(client),
		Configmap:      "test",
		ExpectedValues: map[string]string{"n": "v"},
	})

	validate(t, &testCase{
		Name:        "Should not find test.test1",
		Client:      NewCoreV1(client),
		Configmap:   "test1",
		ExpectError: true,
	})

	client = NewClient(WithKubeClient(kubeClient), WithNamespace("test2"))

	validate(t, &testCase{
		Name:        "Should not find test2.test",
		Client:      NewCoreV1(client),
		Configmap:   "test",
		ExpectError: true,
	})
}

func TestCoreV1SecretData(t *testing.T) {
	type testCase struct {
		Name           string
		Client         *CoreV1
		Secret         string
		ExpectedValues map[string]string
		ShouldExport   bool
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualValues, actualError := tc.Client.SecretData(tc.Secret, tc.ShouldExport)

			assert.Equal(t, tc.ExpectedValues, actualValues)

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
		Client:         NewCoreV1(NewClient(WithKubeClient(kubeClient), WithNamespace("test"))),
		ExpectedValues: map[string]string{"n": "v"},
	})

	validate(t, &testCase{
		Name:        "Should not find test.test1",
		Secret:      "test1",
		Client:      NewCoreV1(NewClient(WithKubeClient(kubeClient), WithNamespace("test"))),
		ExpectError: true,
	})

	validate(t, &testCase{
		Name:        "Should not find test2.test",
		Secret:      "test",
		Client:      NewCoreV1(NewClient(WithKubeClient(kubeClient), WithNamespace("test2"))),
		ExpectError: true,
	})
}

func TestCoreV1Namespaces(t *testing.T) {
	type testCase struct {
		Name          string
		CoreV1        *CoreV1
		ExpectedSlice []string
		ExpectError   bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := tc.CoreV1.NamespaceList()

			assert.Equal(t, tc.ExpectedSlice, actualSlice)

			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	kubeClient := mock.NewFakeClient(mock.Namespace("one"))

	validate(t, &testCase{
		Name:          "Should return a single namespace",
		CoreV1:        NewCoreV1(NewClient(WithKubeClient(kubeClient))),
		ExpectedSlice: []string{"one"},
	})

	kubeClient = mock.NewFakeClient(mock.Namespace("one"), mock.Namespace("two"))

	validate(t, &testCase{
		Name:          "Should return multiple namespaces",
		CoreV1:        NewCoreV1(NewClient(WithKubeClient(kubeClient))),
		ExpectedSlice: []string{"one", "two"},
	})

	kubeClient = mock.NewFakeClient().
		PrependReactor("list", "namespaces", true, &corev1.NamespaceList{}, assert.AnError)

	validate(t, &testCase{
		Name:        "Should return multiple namespaces",
		CoreV1:      NewCoreV1(NewClient(WithKubeClient(kubeClient))),
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
	client := NewClient(
		WithKubeClient(kubeClient),
		WithNamespace("test"),
	)

	validate(t, &testCase{
		Name:     "Should return pods",
		CoreV1:   NewCoreV1(client),
		Resource: "test",
		ExpectedResult: &Result{
			Environment: envValues{"k": "v"},
			Secrets:     map[string]envValues{"test": {"k": "v"}},
			ConfigMaps:  map[string]envValues{"test": {"k": "v"}},
		},
	})

	kubeClient.PrependReactor("get", "pods", true, nil, assert.AnError)
	client = NewClient(
		WithKubeClient(kubeClient),
		WithNamespace("test"),
	)

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
			actualSlice, actualError := tc.CoreV1.PodList()

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
	client := NewClient(
		WithKubeClient(kubeClient),
		WithNamespace("test"),
	)

	validate(t, &testCase{
		Name:          "Should return pods",
		CoreV1:        NewCoreV1(client),
		ExpectedSlice: []string{"test"},
	})

	kubeClient.PrependReactor("list", "pods", true, nil, assert.AnError)
	client = NewClient(
		WithKubeClient(kubeClient),
		WithNamespace("test"),
	)

	validate(t, &testCase{
		Name:        "Should return API errors",
		CoreV1:      NewCoreV1(client),
		ExpectError: true,
	})
}
