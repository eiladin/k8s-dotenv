package options

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type OptionsSuite struct {
	suite.Suite
}

var defaultNamespaceConfig = `
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

var devNamespaceConfig = `
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

var errorConfig = `
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

func (suite OptionsSuite) TestResolveNamespace() {
	cases := []struct {
		namespace     string
		config        string
		expected      string
		expectedError bool
	}{
		{namespace: "test", expected: "test"},
		{expected: "default", config: defaultNamespaceConfig},
		{expected: "dev", config: devNamespaceConfig},
		{expectedError: true, config: errorConfig},
	}

	for _, c := range cases {
		configPath := ""
		if c.config != "" {
			configPath = "./test.config"
			err := ioutil.WriteFile(configPath, []byte(c.config), 0644)
			suite.NoError(err)
		}
		opt := NewOptions()
		opt.Namespace = c.namespace
		err := opt.ResolveNamespace(configPath)

		if configPath != "" {
			os.Remove(configPath)
		}

		if c.expectedError {
			suite.Error(err)
		} else {
			suite.NoError(err)
			suite.Equal(c.expected, opt.Namespace)
		}
	}
}

func (suite OptionsSuite) TestSetWriter() {
	var b bytes.Buffer
	opt := NewOptions()
	err := opt.SetWriter(&b)
	suite.NoError(err)
	suite.Equal(&b, opt.Writer)

	opt = NewOptions()
	err = opt.SetWriter(nil)
	suite.Error(err)

	opt = NewOptions()
	opt.Filename = "./test.out"
	err = opt.SetWriter(nil)
	defer os.Remove(opt.Filename)
	suite.NoError(err)
}

func TestOptionsSuite(t *testing.T) {
	suite.Run(t, new(OptionsSuite))
}
