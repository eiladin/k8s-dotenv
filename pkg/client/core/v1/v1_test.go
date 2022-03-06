package corev1

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func TestCoreV1SetNamespace(t *testing.T) {
	type testCase struct {
		Name              string
		CoreV1            *CoreV1
		Namespace         string
		ExpectedNamespace string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			tc.CoreV1.SetNamespace(tc.Namespace)

			assert.Equal(t, tc.ExpectedNamespace, tc.CoreV1.namespace)
		})
	}

	validate(t, &testCase{
		Name:              "Should set namespace",
		CoreV1:            NewCoreV1(mock.NewFakeClient().CoreV1(), ""),
		Namespace:         "test",
		ExpectedNamespace: "test",
	})
}

func TestNewCoreV1(t *testing.T) {
	type testCase struct {
		Name            string
		CoreV1Interface v1.CoreV1Interface
		Namespace       string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualCoreV1 := NewCoreV1(tc.CoreV1Interface, tc.Namespace)

			assert.NotNil(t, actualCoreV1.CoreV1Interface)
		})
	}

	validate(t, &testCase{
		Name:            "Should set the internal interface",
		CoreV1Interface: mock.NewFakeClient().CoreV1(),
		Namespace:       "test",
	})
}
