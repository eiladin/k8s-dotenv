package client

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestBatchV1CronJob(t *testing.T) {
	type testCase struct {
		Name           string
		BatchV1        *BatchV1
		Resource       string
		ExpectedResult *Result
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualClient := tc.BatchV1.CronJob(tc.Resource)

			assert.Equal(t, tc.ExpectedResult, actualClient.result)

			if tc.ExpectError {
				assert.Error(t, tc.BatchV1.client.Error)
			} else {
				assert.NoError(t, tc.BatchV1.client.Error)
			}
		})
	}

	mockv1 := mock.CronJobv1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:     "Should return cronjobs",
		BatchV1:  client.batchv1,
		Resource: "test",
		ExpectedResult: &Result{
			Environment: map[string]string{"k": "v"},
			Secrets:     []string{"test"},
			ConfigMaps:  []string{"test"},
		},
	})

	kubeClient.PrependReactor("get", "cronjobs", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:        "Should return API errors",
		BatchV1:     NewBatchV1(client),
		Resource:    "test",
		ExpectError: true,
	})
}

func TestBatchV1CronJobs(t *testing.T) {
	type testCase struct {
		Name          string
		BatchV1       *BatchV1
		ExpectedSlice []string
		ExpectError   bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := tc.BatchV1.CronJobs()

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	mockv1 := mock.CronJobv1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	kubeClient := mock.NewFakeClient(mockv1)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:          "Should return cronjobs",
		BatchV1:       NewBatchV1(client),
		ExpectedSlice: []string{"test"},
	})

	kubeClient.PrependReactor("list", "cronjobs", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:        "Should return API errors",
		BatchV1:     NewBatchV1(client),
		ExpectError: true,
	})
}

func TestBatchV1Job(t *testing.T) {
	type testCase struct {
		Name           string
		BatchV1        *BatchV1
		Resource       string
		ExpectedResult *Result
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualClient := tc.BatchV1.Job(tc.Resource)

			assert.Equal(t, tc.ExpectedResult, actualClient.result)

			if tc.ExpectError {
				assert.Error(t, tc.BatchV1.client.Error)
			} else {
				assert.NoError(t, tc.BatchV1.client.Error)
			}
		})
	}

	mockv1 := mock.Job("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:     "Should return jobs",
		BatchV1:  NewBatchV1(client),
		Resource: "test",
		ExpectedResult: &Result{
			Environment: map[string]string{"k": "v"},
			Secrets:     []string{"test"},
			ConfigMaps:  []string{"test"},
		},
	})

	kubeClient.PrependReactor("get", "jobs", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:        "Should return API errors",
		BatchV1:     NewBatchV1(client),
		Resource:    "test",
		ExpectError: true,
	})
}

func TestBatchV1Jobs(t *testing.T) {
	type testCase struct {
		Name          string
		BatchV1       *BatchV1
		ExpectedSlice []string
		ExpectError   bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := tc.BatchV1.Jobs()

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	mockv1 := mock.Job("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	kubeClient := mock.NewFakeClient(mockv1)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:          "Should return jobs",
		BatchV1:       NewBatchV1(client),
		ExpectedSlice: []string{"test"},
	})

	kubeClient.PrependReactor("list", "jobs", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:        "Should return API errors",
		BatchV1:     NewBatchV1(client),
		ExpectError: true,
	})
}
