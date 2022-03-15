package v1

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"github.com/eiladin/k8s-dotenv/pkg/result"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestReplicaSet(t *testing.T) {
	type testCase struct {
		Name           string
		AppsV1         *AppsV1
		Resource       string
		ExpectedResult *result.Result
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualResult := tc.AppsV1.ReplicaSet(tc.Resource)

			assert.Equal(t, tc.ExpectedResult, actualResult)
		})
	}

	mockv1 := mock.ReplicaSet("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)

	validate(t, &testCase{
		Name:     "Should return replicasets",
		AppsV1:   NewAppsV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
		Resource: "test",
		ExpectedResult: &result.Result{
			Environment: result.EnvValues{"k": "v"},
			Secrets:     map[string]result.EnvValues{"test": {"k": "v"}},
			ConfigMaps:  map[string]result.EnvValues{"test": {"k": "v"}},
		},
	})

	kubeClient.PrependReactor("get", "replicasets", true, nil, assert.AnError)

	validate(t, &testCase{
		Name:           "Should return API errors",
		AppsV1:         NewAppsV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
		Resource:       "test",
		ExpectedResult: result.NewFromError(NewResourceLoadError("ReplicaSet", assert.AnError)),
	})
}

func TestReplicaSetList(t *testing.T) {
	type testCase struct {
		Name          string
		AppsV1        *AppsV1
		ExpectedSlice []string
		ExpectError   bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := tc.AppsV1.ReplicaSetList()

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

	validate(t, &testCase{
		Name:          "Should return jobs",
		AppsV1:        NewAppsV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
		ExpectedSlice: []string{"test"},
	})

	kubeClient.PrependReactor("list", "replicasets", true, nil, assert.AnError)

	validate(t, &testCase{
		Name:        "Should return API errors",
		AppsV1:      NewAppsV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
		ExpectError: true,
	})
}
