package v1beta1

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/environment"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestCronJob(t *testing.T) {
	type testCase struct {
		Name           string
		Client         *client.Client
		Namespace      string
		Resource       string
		ExpectedResult *environment.Result
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualResult, actualError := CronJob(tc.Client, tc.Namespace, tc.Resource)

			assert.Equal(t, tc.ExpectedResult, actualResult)

			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	mockv1 := mock.CronJobv1beta1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	cl := mock.NewFakeClient(mockv1, mockecret, mockConfigMap)

	validate(t, &testCase{
		Name:      "Should return cronjobs",
		Client:    client.NewClient(cl),
		Namespace: "test",
		Resource:  "test",
		ExpectedResult: &environment.Result{
			Environment: map[string]string{"k": "v"},
			Secrets:     []string{"test"},
			ConfigMaps:  []string{"test"},
		},
	})

	cl = mock.NewFakeClient().
		PrependReactor("get", "cronjobs", true, nil, assert.AnError)

	validate(t, &testCase{
		Name:        "Should return API errors",
		Client:      client.NewClient(cl),
		Namespace:   "test",
		Resource:    "test",
		ExpectError: true,
	})
}

func TestCronJobs(t *testing.T) {
	type testCase struct {
		Name          string
		Client        *client.Client
		Namespace     string
		ExpectedSlice []string
		ExpectError   bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := CronJobs(tc.Client, tc.Namespace)

			assert.Equal(t, tc.ExpectedSlice, actualSlice)

			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	mockv1 := mock.CronJobv1beta1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	cl := mock.NewFakeClient(mockv1)

	validate(t, &testCase{
		Name:          "Should return cronjobs",
		Client:        client.NewClient(cl),
		Namespace:     "test",
		ExpectedSlice: []string{"test"},
	})

	cl = mock.NewFakeClient().
		PrependReactor("list", "cronjobs", true, nil, assert.AnError)

	validate(t, &testCase{
		Name:        "Should return API errors",
		Client:      client.NewClient(cl),
		Namespace:   "test",
		ExpectError: true,
	})
}
