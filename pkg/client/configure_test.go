package client

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestWithExport(t *testing.T) {
	type testCase struct {
		Name         string
		ShouldExport bool
		Client       *Client
	}

	validate := func(t *testing.T, testCase *testCase) {
		t.Run(testCase.Name, func(t *testing.T) {
			actualConfigureFunc := WithExport(testCase.ShouldExport)

			actualConfigureFunc(testCase.Client)

			assert.Equal(t, testCase.ShouldExport, testCase.Client.options.ShouldExport)
		})
	}

	validate(t, &testCase{
		Name:         "Should update Client.shouldExport",
		ShouldExport: true,
		Client:       NewClient(WithKubeClient(mock.NewFakeClient())),
	})
}

func TestWithNamespace(t *testing.T) {
	type testCase struct {
		Name      string
		Namespace string
		Client    *Client
	}

	validate := func(t *testing.T, testCase *testCase) {
		t.Run(testCase.Name, func(t *testing.T) {
			actualConfigureFunc := WithNamespace(testCase.Namespace)

			actualConfigureFunc(testCase.Client)

			assert.Equal(t, testCase.Namespace, testCase.Client.options.Namespace)
		})
	}

	validate(t, &testCase{
		Name:      "Should update Client.namespace",
		Namespace: "test",
		Client:    NewClient(WithKubeClient(mock.NewFakeClient())),
	})
}
