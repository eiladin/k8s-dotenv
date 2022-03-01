package v1

import (
	"errors"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mocks"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestNamespaces(t *testing.T) {
	type testCase struct {
		Name string

		Opt *options.Options

		ExpectedSlice []string
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := Namespaces(tc.Opt)

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	client := fake.NewSimpleClientset(mocks.Namespace("one"))
	validate(t, &testCase{
		Name:          "Should return a single namespace",
		Opt:           &options.Options{Client: client},
		ExpectedSlice: []string{"one"},
	})

	client = fake.NewSimpleClientset(mocks.Namespace("one"), mocks.Namespace("two"))
	validate(t, &testCase{
		Name:          "Should return multiple namespaces",
		Opt:           &options.Options{Client: client},
		ExpectedSlice: []string{"one", "two"},
	})

	client = fake.NewSimpleClientset()
	client.PrependReactor("list", "namespaces", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &corev1.NamespaceList{}, errors.New("error getting namespaces")
	})
	validate(t, &testCase{
		Name:          "Should return multiple namespaces",
		Opt:           &options.Options{Client: client},
		ExpectedError: errors.New("error getting namespaces"),
	})
}
