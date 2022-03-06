package client

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestBatchV1Beta1CronJob(t *testing.T) {
	type testCase struct {
		Name           string
		BatchV1Beta1   *BatchV1Beta1
		Resource       string
		ExpectedResult *Result
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualClient := tc.BatchV1Beta1.CronJob(tc.Resource)

			assert.Equal(t, tc.ExpectedResult, actualClient.result)

			if tc.ExpectError {
				assert.Error(t, tc.BatchV1Beta1.client.Error)
			} else {
				assert.NoError(t, tc.BatchV1Beta1.client.Error)
			}
		})
	}

	mockv1 := mock.CronJobv1beta1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:         "Should return cronjobs",
		BatchV1Beta1: NewBatchV1Beta1(client),
		Resource:     "test",
		ExpectedResult: &Result{
			Environment: envValues{"k": "v"},
			Secrets:     map[string]envValues{"test": {"k": "v"}},
			ConfigMaps:  map[string]envValues{"test": {"k": "v"}},
		},
	})

	kubeClient.PrependReactor("get", "cronjobs", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:         "Should return API errors",
		BatchV1Beta1: NewBatchV1Beta1(client),
		Resource:     "test",
		ExpectError:  true,
	})
}

func TestBatchV1Beta1CronJobs(t *testing.T) {
	type testCase struct {
		Name          string
		BatchV1Beta1  *BatchV1Beta1
		ExpectedSlice []string
		ExpectError   bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := tc.BatchV1Beta1.CronJobs()

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	mockv1 := mock.CronJobv1beta1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	kubeClient := mock.NewFakeClient(mockv1)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:          "Should return cronjobs",
		BatchV1Beta1:  NewBatchV1Beta1(client),
		ExpectedSlice: []string{"test"},
	})

	kubeClient.PrependReactor("list", "cronjobs", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:         "Should return API errors",
		BatchV1Beta1: NewBatchV1Beta1(client),
		ExpectError:  true,
	})
}
