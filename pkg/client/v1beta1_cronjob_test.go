package client

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestClientCronJobV1Beta1(t *testing.T) {
	type testCase struct {
		Name           string
		Client         *Client
		Resource       string
		ExpectedResult *Result
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualClient := tc.Client.CronJobV1Beta1(tc.Resource)

			assert.Equal(t, tc.ExpectedResult, actualClient.result)

			if tc.ExpectError {
				assert.Error(t, tc.Client.Error)
			} else {
				assert.NoError(t, tc.Client.Error)
			}
		})
	}

	mockv1beta1 := mock.CronJobv1beta1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1beta1, mockConfigMap, mockSecret)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:     "Should return cronjobs",
		Client:   client,
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
		Client:      client,
		Resource:    "test",
		ExpectError: true,
	})
}

func TestClientCronJobsV1Beta1(t *testing.T) {
	type testCase struct {
		Name          string
		Client        *Client
		ExpectedSlice []string
		ExpectError   bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := tc.Client.CronJobsV1beta1()

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	mockv1beta1 := mock.CronJobv1beta1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	kubeClient := mock.NewFakeClient(mockv1beta1)
	client := NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:          "Should return cronjobs",
		Client:        client,
		ExpectedSlice: []string{"test"},
	})

	kubeClient.PrependReactor("list", "cronjobs", true, nil, assert.AnError)
	client = NewClient(kubeClient, WithNamespace("test"))

	validate(t, &testCase{
		Name:        "Should return API errors",
		Client:      client,
		ExpectError: true,
	})
}
