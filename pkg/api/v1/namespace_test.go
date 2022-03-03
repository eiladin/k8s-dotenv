package v1

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestNamespaces(t *testing.T) {
	type testCase struct {
		Name          string
		Client        *client.Client
		ExpectedSlice []string
		ExpectError   bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualSlice, actualError := Namespaces(tc.Client)

			assert.Equal(t, tc.ExpectedSlice, actualSlice)

			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	cl := mock.NewFakeClient(mock.Namespace("one"))

	validate(t, &testCase{
		Name:          "Should return a single namespace",
		Client:        client.NewClient(cl),
		ExpectedSlice: []string{"one"},
	})

	cl = mock.NewFakeClient(mock.Namespace("one"), mock.Namespace("two"))

	validate(t, &testCase{
		Name:          "Should return multiple namespaces",
		Client:        client.NewClient(cl),
		ExpectedSlice: []string{"one", "two"},
	})

	cl = mock.NewFakeClient().
		PrependReactor("list", "namespaces", true, &corev1.NamespaceList{}, assert.AnError)

	validate(t, &testCase{
		Name:        "Should return multiple namespaces",
		Client:      client.NewClient(cl),
		ExpectError: true,
	})
}
