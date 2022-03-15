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

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualConfigureFunc := WithExport(tc.ShouldExport)

			actualConfigureFunc(tc.Client)

			assert.Equal(t, tc.ShouldExport, tc.Client.options.ShouldExport)
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

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualConfigureFunc := WithNamespace(tc.Namespace)

			actualConfigureFunc(tc.Client)

			assert.Equal(t, tc.Namespace, tc.Client.options.Namespace)
		})
	}

	validate(t, &testCase{
		Name:      "Should update Client.namespace",
		Namespace: "test",
		Client:    NewClient(WithKubeClient(mock.NewFakeClient())),
	})
}
