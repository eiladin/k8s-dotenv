package client

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestAppsV1DaemonSet(t *testing.T) {
	type testCase struct {
		Name           string
		AppsV1         *AppsV1
		Resource       string
		ExpectedResult *Result
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualClient := tc.AppsV1.DaemonSet(tc.Resource)

			assert.Equal(t, tc.ExpectedResult, actualClient.result)

			if tc.ExpectError {
				assert.Error(t, tc.AppsV1.client.Error)
			} else {
				assert.NoError(t, tc.AppsV1.client.Error)
			}
		})
	}

	mockv1 := mock.DaemonSet("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:     "Should return jobs",
		AppsV1:   NewAppsV1(client),
		Resource: "test",
		ExpectedResult: &Result{
			Environment: envValues{"k": "v"},
			Secrets:     map[string]envValues{"test": {"k": "v"}},
			ConfigMaps:  map[string]envValues{"test": {"k": "v"}},
		},
	})

	kubeClient.PrependReactor("get", "daemonsets", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:        "Should return API errors",
		AppsV1:      NewAppsV1(client),
		Resource:    "test",
		ExpectError: true,
	})
}

func TestAppsV1DaemonSets(t *testing.T) {
	type testCase struct {
		Name          string
		AppsV1        *AppsV1
		ExpectedSlice []string
		ExpectError   bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := tc.AppsV1.DaemonSets()

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	mockv1 := mock.DaemonSet("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	kubeClient := mock.NewFakeClient(mockv1)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:          "Should return jobs",
		AppsV1:        NewAppsV1(client),
		ExpectedSlice: []string{"test"},
	})

	kubeClient.PrependReactor("list", "daemonsets", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:        "Should return API errors",
		AppsV1:      NewAppsV1(client),
		ExpectError: true,
	})
}

func TestAppsV1Deployment(t *testing.T) {
	type testCase struct {
		Name           string
		AppsV1         *AppsV1
		Resource       string
		ExpectedResult *Result
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualClient := tc.AppsV1.Deployment(tc.Resource)

			assert.Equal(t, tc.ExpectedResult, actualClient.result)

			if tc.ExpectError {
				assert.Error(t, tc.AppsV1.client.Error)
			} else {
				assert.NoError(t, tc.AppsV1.client.Error)
			}
		})
	}

	mockv1 := mock.Deployment("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:     "Should return jobs",
		AppsV1:   NewAppsV1(client),
		Resource: "test",
		ExpectedResult: &Result{
			Environment: envValues{"k": "v"},
			Secrets:     map[string]envValues{"test": {"k": "v"}},
			ConfigMaps:  map[string]envValues{"test": {"k": "v"}},
		},
	})

	kubeClient.PrependReactor("get", "deployments", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:        "Should return API errors",
		AppsV1:      NewAppsV1(client),
		Resource:    "test",
		ExpectError: true,
	})
}

func TestAppsV1Deployments(t *testing.T) {
	type testCase struct {
		Name          string
		AppsV1        *AppsV1
		ExpectedSlice []string
		ExpectError   bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := tc.AppsV1.Deployments()

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	mockv1 := mock.Deployment("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	kubeClient := mock.NewFakeClient(mockv1)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:          "Should return jobs",
		AppsV1:        NewAppsV1(client),
		ExpectedSlice: []string{"test"},
	})

	kubeClient.PrependReactor("list", "deployments", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:        "Should return API errors",
		AppsV1:      NewAppsV1(client),
		ExpectError: true,
	})
}

func TestAppsV1ReplicaSet(t *testing.T) {
	type testCase struct {
		Name           string
		AppsV1         *AppsV1
		Resource       string
		ExpectedResult *Result
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualClient := tc.AppsV1.ReplicaSet(tc.Resource)

			assert.Equal(t, tc.ExpectedResult, actualClient.result)

			if tc.ExpectError {
				assert.Error(t, tc.AppsV1.client.Error)
			} else {
				assert.NoError(t, tc.AppsV1.client.Error)
			}
		})
	}

	mockv1 := mock.ReplicaSet("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:     "Should return jobs",
		AppsV1:   NewAppsV1(client),
		Resource: "test",
		ExpectedResult: &Result{
			Environment: envValues{"k": "v"},
			Secrets:     map[string]envValues{"test": {"k": "v"}},
			ConfigMaps:  map[string]envValues{"test": {"k": "v"}},
		},
	})

	kubeClient.PrependReactor("get", "replicasets", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:        "Should return API errors",
		AppsV1:      NewAppsV1(client),
		Resource:    "test",
		ExpectError: true,
	})
}

func TestAppsV1ReplicaSets(t *testing.T) {
	type testCase struct {
		Name          string
		AppsV1        *AppsV1
		ExpectedSlice []string
		ExpectError   bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := tc.AppsV1.ReplicaSets()

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	mockv1 := mock.ReplicaSet("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	kubeClient := mock.NewFakeClient(mockv1)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:          "Should return jobs",
		AppsV1:        NewAppsV1(client),
		ExpectedSlice: []string{"test"},
	})

	kubeClient.PrependReactor("list", "replicasets", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:        "Should return API errors",
		AppsV1:      NewAppsV1(client),
		ExpectError: true,
	})
}
