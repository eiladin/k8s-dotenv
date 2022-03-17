package v1

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestNamespaceList(t *testing.T) {
	type testCase struct {
		Name          string
		CoreV1        *CoreV1
		ExpectedSlice []string
		ExpectError   bool
	}

	validate := func(t *testing.T, testCase *testCase) {
		t.Run(testCase.Name, func(t *testing.T) {
			actualSlice, actualError := testCase.CoreV1.NamespaceList()

			assert.Equal(t, testCase.ExpectedSlice, actualSlice)

			if testCase.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	kubeClient := mock.NewFakeClient(mock.Namespace("one"))

	validate(t, &testCase{
		Name:          "Should return a single namespace",
		CoreV1:        NewCoreV1(kubeClient, clientoptions.New()),
		ExpectedSlice: []string{"one"},
	})

	kubeClient = mock.NewFakeClient(mock.Namespace("one"), mock.Namespace("two"))

	validate(t, &testCase{
		Name:          "Should return multiple namespaces",
		CoreV1:        NewCoreV1(kubeClient, clientoptions.New()),
		ExpectedSlice: []string{"one", "two"},
	})

	kubeClient = mock.NewFakeClient().
		PrependReactor("list", "namespaces", true, &corev1.NamespaceList{}, assert.AnError)

	validate(t, &testCase{
		Name:        "Should return multiple namespaces",
		CoreV1:      NewCoreV1(kubeClient, clientoptions.New()),
		ExpectError: true,
	})
}
