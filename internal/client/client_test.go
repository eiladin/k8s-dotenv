package client

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type ClientSuite struct {
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

func (suite ClientSuite) TestCurrentNamespace() {
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

		got, err := CurrentNamespace(c.namespace, configPath)

		if configPath != "" {
			os.Remove(configPath)
		}

		if c.expectedError {
			suite.Error(err)
		} else {
			suite.NoError(err)
		}
		suite.Equal(c.expected, got)
	}
}

func (suite ClientSuite) TestGetApiGroup() {
	cases := []struct {
		resource  string
		shouldErr bool
	}{
		{resource: "Job"},
		{resource: "Job", shouldErr: true},
	}

	for _, c := range cases {
		client := fake.NewSimpleClientset(&v1.Job{})

		if !c.shouldErr {
			client.Fake.Resources = append(client.Fake.Resources, &metav1.APIResourceList{
				GroupVersion: "v1",
				APIResources: []metav1.APIResource{
					{Name: "Jobs", SingularName: "Job", Kind: "Job", Namespaced: false, Group: "v1"},
				},
			})
		}

		got, err := GetApiGroup(client, c.resource)
		if c.shouldErr {
			suite.Error(err)
		} else {
			suite.NoError(err)
			suite.NotNil(got)
		}
	}
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
