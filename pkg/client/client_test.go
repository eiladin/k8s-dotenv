package client

import (
	"bytes"
	"os"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/batch/v1"
)

func TestClientGetAPIGroup(t *testing.T) {
	type testCase struct {
		Name           string
		Client         *Client
		Resource       string
		ExpectedString string
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := tc.Client.GetAPIGroup(tc.Resource)

			assert.Equal(t, tc.ExpectedString, actualString)

			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	kubeClient := mock.NewFakeClient(&v1.Job{}).WithResources(mock.Jobv1Resource())

	validate(t, &testCase{
		Name:           "Should detect resource group",
		Client:         NewClient(WithKubeClient(kubeClient)),
		Resource:       "Job",
		ExpectedString: "v1",
	})

	kubeClient = mock.NewFakeClient(&v1.Job{})

	validate(t, &testCase{
		Name:        "Should error if the resource is not found",
		Client:      NewClient(WithKubeClient(kubeClient)),
		Resource:    "Job",
		ExpectError: true,
	})

	kubeClient = mock.NewFakeClient(&v1.Job{}).WithResources(mock.InvalidGroupResource())

	validate(t, &testCase{
		Name:        "Should return API errors",
		Client:      NewClient(WithKubeClient(kubeClient)),
		Resource:    "Job",
		ExpectError: true,
	})
}

func TestClientSetDefaultWriter(t *testing.T) {
	type testCase struct {
		Name          string
		Client        *Client
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := tc.Client.setDefaultWriter()

			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	var b bytes.Buffer

	defer os.Remove("./out.test")

	validate(t, &testCase{
		Name:   "Should use the passed in writer",
		Client: NewClient(WithKubeClient(mock.NewFakeClient()), WithWriter(&b)),
	})

	validate(t, &testCase{
		Name:          "Should Error given no filename or writer",
		Client:        NewClient(WithKubeClient(mock.NewFakeClient())),
		ExpectedError: ErrNoFilename,
	})

	validate(t, &testCase{
		Name:   "Should not error given a filename",
		Client: NewClient(WithKubeClient(mock.NewFakeClient()), WithFilename("./out.test")),
	})
}

func TestNewClient(t *testing.T) {
	type testCase struct {
		Name           string
		Configures     []ConfigureFunc
		ExpectedClient *Client
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualClient := NewClient(tc.Configures...)

			assert.Equal(t, tc.ExpectedClient.shouldExport, actualClient.shouldExport)
			assert.Equal(t, tc.ExpectedClient.namespace, actualClient.namespace)
			assert.Equal(t, tc.ExpectedClient.filename, actualClient.filename)
			assert.Equal(t, tc.ExpectedClient.writer, actualClient.writer)
			assert.Equal(t, tc.ExpectedClient.Error, actualClient.Error)
		})
	}

	validate(t, &testCase{
		Name: "Should run configures",
		Configures: []ConfigureFunc{
			WithKubeClient(mock.NewFakeClient()),
			WithFilename("string"),
		},
		ExpectedClient: &Client{
			Interface:    mock.NewFakeClient(),
			shouldExport: false,
			namespace:    "",
			filename:     "string",
			writer:       nil,
			Error:        nil,
		},
	})
}

func TestAPIClients(t *testing.T) {
	client := NewClient(WithKubeClient(mock.NewFakeClient()))

	assert.NotNil(t, client.AppsV1())
	assert.NotNil(t, client.BatchV1())
	assert.NotNil(t, client.BatchV1Beta1())
	assert.NotNil(t, client.CoreV1())
}

func TestAPIClientsPanics(t *testing.T) {
	client := NewClient()

	assert.Panics(t, func() { client.AppsV1() })
	assert.Panics(t, func() { client.BatchV1() })
	assert.Panics(t, func() { client.BatchV1Beta1() })
	assert.Panics(t, func() { client.CoreV1() })
}
