package v1

import (
	"errors"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/environment"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestDeployment(t *testing.T) {
	type testCase struct {
		Name string

		Opt *options.Options

		ExpectedResult *environment.Result
		ExpectedError  error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualResult, actualError := Deployment(tc.Opt)

			assert.Equal(t, tc.ExpectedResult, actualResult)
			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	mockv1 := mock.Deployment("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	client := fake.NewSimpleClientset(mockv1, mockecret, mockConfigMap)
	validate(t, &testCase{
		Name: "Should return deployments",
		Opt: &options.Options{
			Client:    client,
			Namespace: "test",
			Name:      "test",
		},
		ExpectedResult: &environment.Result{
			Environment: map[string]string{"k": "v"},
			Secrets:     []string{"test"},
			ConfigMaps:  []string{"test"},
		},
	})

	client = fake.NewSimpleClientset()
	client.PrependReactor("get", "deployments", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("error getting deployment")
	})
	validate(t, &testCase{
		Name: "Should return API errors",
		Opt: &options.Options{
			Client:    client,
			Namespace: "test",
			Name:      "test",
		},
		ExpectedError: errors.New("error getting deployment"),
	})
}

func TestDeployments(t *testing.T) {
	type testCase struct {
		Name string

		Opt *options.Options

		ExpectedSlice []string
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := Deployments(tc.Opt)

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	mockv1 := mock.Deployment("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	client := fake.NewSimpleClientset(mockv1)
	validate(t, &testCase{
		Name: "Should return deployments",
		Opt: &options.Options{
			Client:    client,
			Namespace: "test",
			Name:      "test",
		},
		ExpectedSlice: []string{"test"},
	})

	client = fake.NewSimpleClientset()
	client.PrependReactor("list", "deployments", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("error getting deployment list")
	})
	validate(t, &testCase{
		Name: "Should return API errors",
		Opt: &options.Options{
			Client:    client,
			Namespace: "test",
			Name:      "test",
		},
		ExpectedError: errors.New("error getting deployment list"),
	})
}
