package v1

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"github.com/eiladin/k8s-dotenv/pkg/result"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestCronJob(t *testing.T) {
	type testCase struct {
		Name           string
		BatchV1        *BatchV1
		Resource       string
		ExpectedResult *result.Result
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualResult := tc.BatchV1.CronJob(tc.Resource)

			assert.Equal(t, tc.ExpectedResult, actualResult)
		})
	}

	mockv1 := mock.CronJobv1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)

	validate(t, &testCase{
		Name:     "Should return cronjobs",
		BatchV1:  NewBatchV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
		Resource: "test",
		ExpectedResult: &result.Result{
			Environment: result.EnvValues{"k": "v"},
			Secrets:     map[string]result.EnvValues{"test": {"k": "v"}},
			ConfigMaps:  map[string]result.EnvValues{"test": {"k": "v"}},
		},
	})

	kubeClient.PrependReactor("get", "cronjobs", true, nil, assert.AnError)

	validate(t, &testCase{
		Name:           "Should return API errors",
		BatchV1:        NewBatchV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
		Resource:       "test",
		ExpectedResult: result.NewFromError(NewResourceLoadError("CronJob", assert.AnError)),
	})
}

func TestCronJobList(t *testing.T) {
	type testCase struct {
		Name          string
		BatchV1       *BatchV1
		ExpectedSlice []string
		ExpectError   bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := tc.BatchV1.CronJobList()

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

	validate(t, &testCase{
		Name:          "Should return cronjobs",
		BatchV1:       NewBatchV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
		ExpectedSlice: []string{"test"},
	})

	kubeClient.PrependReactor("list", "cronjobs", true, nil, assert.AnError)

	validate(t, &testCase{
		Name:        "Should return API errors",
		BatchV1:     NewBatchV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
		ExpectError: true,
	})
}
