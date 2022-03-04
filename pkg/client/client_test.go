package client

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/batch/v1"
)

const defaultNamespaceConfig = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: 9VSmmKMhYNKBoxopdbbgiw==
    server: https://not-a-real-cluster
  name: dev
contexts:
- context:
    cluster: dev
    user: dev
  name: dev
current-context: dev
kind: Config
preferences: {}
users:
- name: dev
  user:
    token: not-a-real-token
`

const devNamespaceConfig = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: 9VSmmKMhYNKBoxopdbbgiw==
    server: https://not-a-real-cluster
  name: dev
contexts:
- context:
    cluster: dev
    namespace: dev
    user: dev
  name: dev
current-context: dev
kind: Config
preferences: {}
users:
- name: dev
  user:
    token: not-a-real-token
`

const errorConfig = `
	apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: 9VSmmKMhYNKBoxopdbbgiw==
    server: https://not-a-real-cluster
  name: dev
contexts:
- context:
    cluster: dev
    namespace: dev
    user: dev
  name: dev
current-context: dev
kind: Config
preferences: {}
users:
- name: dev
  user:
    token: not-a-real-token
`

func TestCurrentNamespace(t *testing.T) {
	type testCase struct {
		Name           string
		Namespace      string
		ConfigPath     string
		ExpectedString string
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := CurrentNamespace(tc.Namespace, tc.ConfigPath)

			assert.Equal(t, tc.ExpectedString, actualString)

			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	validate(t, &testCase{
		Name:           "Should return passed in namespace",
		Namespace:      "test",
		ExpectedString: "test",
	})

	err := ioutil.WriteFile("./default.config", []byte(defaultNamespaceConfig), 0600)
	assert.NoError(t, err)

	defer os.Remove("./default.config")

	validate(t, &testCase{
		Name:           "Should resolve default",
		ConfigPath:     "default.config",
		ExpectedString: "default",
	})

	err = ioutil.WriteFile("./dev.config", []byte(devNamespaceConfig), 0600)
	assert.NoError(t, err)

	defer os.Remove("./dev.config")

	validate(t, &testCase{
		Name:           "Should resolve dev",
		ConfigPath:     "dev.config",
		ExpectedString: "dev",
	})

	err = ioutil.WriteFile("./error.config", []byte(errorConfig), 0600)
	assert.NoError(t, err)

	defer os.Remove("./error.config")

	validate(t, &testCase{
		Name:        "Should throw an error on invalid config",
		ConfigPath:  "error.config",
		ExpectError: true,
	})
}

func TestGetAPIGroup(t *testing.T) {
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

	cl := mock.NewFakeClient(&v1.Job{}).WithResources(mock.Jobv1Resource())

	validate(t, &testCase{
		Name:           "Should detect resource group",
		Client:         NewClient(cl),
		Resource:       "Job",
		ExpectedString: "v1",
	})

	cl = mock.NewFakeClient(&v1.Job{})

	validate(t, &testCase{
		Name:        "Should error if the resource is not found",
		Client:      NewClient(cl),
		Resource:    "Job",
		ExpectError: true,
	})

	cl = mock.NewFakeClient(&v1.Job{}).WithResources(mock.InvalidGroupResource())

	validate(t, &testCase{
		Name:        "Should return API errors",
		Client:      NewClient(cl),
		Resource:    "Job",
		ExpectError: true,
	})
}
