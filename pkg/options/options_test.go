package options

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

func TestOptionsResolveNamespace(t *testing.T) {
	type testCase struct {
		Name         string
		Options      *Options
		ConfigPath   string
		ErrorChecker func(err error) bool
		ValueChecker func(opt *Options) bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := tc.Options.ResolveNamespace(tc.ConfigPath)

			if tc.ErrorChecker != nil {
				assert.True(t, tc.ErrorChecker(actualError))
			}

			if tc.ValueChecker != nil {
				assert.True(t, tc.ValueChecker(tc.Options))
			}
		})
	}

	validate(t, &testCase{
		Name:    "Should resolve test",
		Options: &Options{Namespace: "test"},
		ValueChecker: func(opt *Options) bool {
			return opt.Namespace == "test"
		},
	})

	err := ioutil.WriteFile("./default.config", []byte(defaultNamespaceConfig), 0600)
	assert.NoError(t, err)

	defer os.Remove("./default.config")

	validate(t, &testCase{
		Name:       "Should resolve default",
		Options:    &Options{},
		ConfigPath: "default.config",
		ValueChecker: func(opt *Options) bool {
			return opt.Namespace == "default"
		},
	})

	err = ioutil.WriteFile("./dev.config", []byte(devNamespaceConfig), 0600)
	assert.NoError(t, err)

	defer os.Remove("./dev.config")

	validate(t, &testCase{
		Name:       "Should resolve dev",
		Options:    &Options{},
		ConfigPath: "dev.config",
		ValueChecker: func(opt *Options) bool {
			return opt.Namespace == "dev"
		},
	})

	err = ioutil.WriteFile("./error.config", []byte(errorConfig), 0600)
	assert.NoError(t, err)

	defer os.Remove("./error.config")

	validate(t, &testCase{
		Name:       "Should throw an error on invalid config",
		Options:    &Options{},
		ConfigPath: "error.config",
		ErrorChecker: func(err error) bool {
			return err != nil
		},
	})
}
