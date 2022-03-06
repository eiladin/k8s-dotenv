package kubeclient

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
		ExpectedError  error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := CurrentNamespace(tc.Namespace, tc.ConfigPath)

			assert.Equal(t, tc.ExpectedString, actualString)
			assert.Equal(t, tc.ExpectedError, actualError)
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
		Name:          "Should throw an error on invalid config",
		ConfigPath:    "error.config",
		ExpectedError: ErrNamespaceResolution,
	})
}
