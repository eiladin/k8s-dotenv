package client

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

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

func TestCurrentNamespace(t *testing.T) {
	type testCase struct {
		Name string

		Namespace  string
		ConfigPath string

		ExpectedString string
		ExpectedError  error
		ErrorChecker   func(err error) bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := CurrentNamespace(tc.Namespace, tc.ConfigPath)

			assert.Equal(t, tc.ExpectedString, actualString)
			if tc.ErrorChecker != nil {
				assert.True(t, tc.ErrorChecker(actualError))
			} else {
				assert.Equal(t, tc.ExpectedError, actualError)
			}
		})
	}

	validate(t, &testCase{
		Name:           "Should return passed in namespace",
		Namespace:      "test",
		ExpectedString: "test",
	})

	err := ioutil.WriteFile("./default.config", []byte(defaultNamespaceConfig), 0644)
	assert.NoError(t, err)
	defer os.Remove("./default.config")
	validate(t, &testCase{
		Name:           "Should resolve default",
		ConfigPath:     "default.config",
		ExpectedString: "default",
	})

	err = ioutil.WriteFile("./dev.config", []byte(devNamespaceConfig), 0644)
	assert.NoError(t, err)
	defer os.Remove("./dev.config")
	validate(t, &testCase{
		Name:           "Should resolve dev",
		ConfigPath:     "dev.config",
		ExpectedString: "dev",
	})

	err = ioutil.WriteFile("./error.config", []byte(errorConfig), 0644)
	assert.NoError(t, err)
	defer os.Remove("./error.config")
	validate(t, &testCase{
		Name:       "Should throw an error on invalid config",
		ConfigPath: "error.config",
		ErrorChecker: func(err error) bool {
			return err != nil
		},
	})
}

func TestGetAPIGroup(t *testing.T) {
	type testCase struct {
		Name string

		Client   Client
		Resource string

		ExpectedString string
		ExpectedError  error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := tc.Client.GetAPIGroup(tc.Resource)

			assert.Equal(t, tc.ExpectedString, actualString)
			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	cl := fake.NewSimpleClientset(&v1.Job{})
	cl.Fake.Resources = []*metav1.APIResourceList{
		{
			GroupVersion: "v1",
			APIResources: []metav1.APIResource{
				{Name: "Jobs", SingularName: "Job", Kind: "Job", Namespaced: false, Group: "v1"},
			},
		},
	}
	validate(t, &testCase{
		Name:           "Should detect resource group",
		Client:         Client{NewClient(cl)},
		Resource:       "Job",
		ExpectedString: "v1",
	})

	cl = fake.NewSimpleClientset(&v1.Job{})
	validate(t, &testCase{
		Name:          "Should error if the resource is not found",
		Client:        Client{NewClient(cl)},
		Resource:      "Job",
		ExpectedError: errors.New("resource Job not found"),
	})

	cl = fake.NewSimpleClientset(&v1.Job{})
	cl.Fake.Resources = []*metav1.APIResourceList{
		{
			GroupVersion: "a/b/c",
			APIResources: []metav1.APIResource{
				{Name: "Jobs", SingularName: "Job", Kind: "Job", Namespaced: false, Group: "a/b/c"},
			},
		},
	}
	validate(t, &testCase{
		Name:          "Should return API errors",
		Client:        Client{NewClient(cl)},
		Resource:      "Job",
		ExpectedError: errors.New("unexpected GroupVersion string: a/b/c"),
	})
}
