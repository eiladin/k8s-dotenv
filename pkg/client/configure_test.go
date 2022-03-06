package client

import (
	"bytes"
	"io"
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

			assert.Equal(t, tc.ShouldExport, tc.Client.shouldExport)
		})
	}

	validate(t, &testCase{
		Name:         "Should update Client.shouldExport",
		ShouldExport: true,
		Client:       NewClient(mock.NewFakeClient()),
	})
}

func TestWithFilename(t *testing.T) {
	type testCase struct {
		Name     string
		Filename string
		Client   *Client
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualConfigureFunc := WithFilename(tc.Filename)

			actualConfigureFunc(tc.Client)

			assert.Equal(t, tc.Filename, tc.Client.filename)
		})
	}

	validate(t, &testCase{
		Name:     "Should update Client.filename",
		Filename: "test",
		Client:   NewClient(mock.NewFakeClient()),
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

			assert.Equal(t, tc.Namespace, tc.Client.namespace)
		})
	}

	validate(t, &testCase{
		Name:      "Should update Client.namespace",
		Namespace: "test",
		Client:    NewClient(mock.NewFakeClient()),
	})
}

func TestWithWriter(t *testing.T) {
	type testCase struct {
		Name   string
		Writer io.Writer
		Client *Client
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualConfigureFunc := WithWriter(tc.Writer)

			actualConfigureFunc(tc.Client)

			assert.Equal(t, tc.Writer, tc.Client.writer)
		})
	}

	var b bytes.Buffer

	validate(t, &testCase{
		Name:   "Should update Client.writer",
		Writer: &b,
		Client: NewClient(mock.NewFakeClient()),
	})
}
