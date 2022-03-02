package v1

import (
	"errors"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/environment"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestDaemonSet(t *testing.T) {
	type testCase struct {
		Name string

		Client         *client.Client
		Namespace      string
		Resource       string
		ExpectedResult *environment.Result
		ExpectedError  error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualResult, actualError := DaemonSet(tc.Client, tc.Namespace, tc.Resource)

			assert.Equal(t, tc.ExpectedResult, actualResult)
			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	mockv1 := mock.DaemonSet("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	cl := fake.NewSimpleClientset(mockv1, mockecret, mockConfigMap)
	validate(t, &testCase{
		Name:      "Should return daemonsets",
		Client:    client.NewClient(cl),
		Namespace: "test",
		Resource:  "test",
		ExpectedResult: &environment.Result{
			Environment: map[string]string{"k": "v"},
			Secrets:     []string{"test"},
			ConfigMaps:  []string{"test"},
		},
	})

	cl = fake.NewSimpleClientset()
	cl.PrependReactor("get", "daemonsets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("error getting daemonset")
	})
	validate(t, &testCase{
		Name:          "Should return API errors",
		Client:        client.NewClient(cl),
		Namespace:     "test",
		Resource:      "test",
		ExpectedError: errors.New("error getting daemonset"),
	})
}

func TestDaemonSets(t *testing.T) {
	type testCase struct {
		Name string

		Client    *client.Client
		Namespace string

		ExpectedSlice []string
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := DaemonSets(tc.Client, tc.Namespace)

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	mockv1 := mock.DaemonSet("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	cl := fake.NewSimpleClientset(mockv1)
	validate(t, &testCase{
		Name:          "Should return daemonsets",
		Client:        client.NewClient(cl),
		Namespace:     "test",
		ExpectedSlice: []string{"test"},
	})

	cl = fake.NewSimpleClientset()
	cl.PrependReactor("list", "daemonsets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("error getting daemonset list")
	})
	validate(t, &testCase{
		Name:          "Should return API errors",
		Client:        client.NewClient(cl),
		Namespace:     "test",
		ExpectedError: errors.New("error getting daemonset list"),
	})
}
