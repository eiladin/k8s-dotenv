package v1

import (
	"errors"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/environment"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mocks"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestCronJob(t *testing.T) {
	type testCase struct {
		Name string

		Opt *options.Options

		ExpectedResult *environment.Result
		ExpectedError  error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualResult, actualError := CronJob(tc.Opt)

			assert.Equal(t, tc.ExpectedResult, actualResult)
			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	mockv1 := mocks.CronJobv1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mocks.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mocks.ConfigMap("test", "test", map[string]string{"k": "v"})
	client := fake.NewSimpleClientset(mockv1, mockSecret, mockConfigMap)
	validate(t, &testCase{
		Name: "Should return cronjobs",
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
	client.PrependReactor("get", "cronjobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("error getting cronjob")
	})
	validate(t, &testCase{
		Name: "Should return API errors",
		Opt: &options.Options{
			Client:    client,
			Namespace: "test",
			Name:      "test",
		},
		ExpectedError: errors.New("error getting cronjob"),
	})
}

func TestCronJobs(t *testing.T) {
	type testCase struct {
		Name string

		Opt *options.Options

		ExpectedSlice []string
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := CronJobs(tc.Opt)

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	mockv1 := mocks.CronJobv1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	client := fake.NewSimpleClientset(mockv1)
	validate(t, &testCase{
		Name: "Should return cronjobs",
		Opt: &options.Options{
			Client:    client,
			Namespace: "test",
			Name:      "test",
		},
		ExpectedSlice: []string{"test"},
	})

	client = fake.NewSimpleClientset()
	client.PrependReactor("list", "cronjobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, errors.New("error getting cronjob list")
	})
	validate(t, &testCase{
		Name: "Should return API errors",
		Opt: &options.Options{
			Client:    client,
			Namespace: "test",
			Name:      "test",
		},
		ExpectedError: errors.New("error getting cronjob list"),
	})
}
