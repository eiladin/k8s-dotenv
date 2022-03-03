package v1

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestNamespaces(t *testing.T) {
	type testCase struct {
		Name string

		Client *client.Client

		ExpectedSlice []string
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := Namespaces(tc.Client)

			assert.Equal(t, tc.ExpectedSlice, actualSlice)
			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	cl := fake.NewSimpleClientset(mock.Namespace("one"))

	validate(t, &testCase{
		Name:          "Should return a single namespace",
		Client:        client.NewClient(cl),
		ExpectedSlice: []string{"one"},
	})

	cl = fake.NewSimpleClientset(mock.Namespace("one"), mock.Namespace("two"))

	validate(t, &testCase{
		Name:          "Should return multiple namespaces",
		Client:        client.NewClient(cl),
		ExpectedSlice: []string{"one", "two"},
	})

	cl = fake.NewSimpleClientset()
	cl.PrependReactor("list", "namespaces", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &corev1.NamespaceList{}, mock.NewError("error getting namespaces")
	})

	validate(t, &testCase{
		Name:          "Should return multiple namespaces",
		Client:        client.NewClient(cl),
		ExpectedError: NewResourceLoadError(mock.NewError("error getting namespaces")),
	})
}
