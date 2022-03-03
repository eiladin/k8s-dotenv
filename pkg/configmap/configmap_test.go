package configmap

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGet(t *testing.T) {
	type testCase struct {
		Name string

		Client       *client.Client
		Namespace    string
		Configmap    string
		ShouldExport bool

		ExpectedString string
		ExpectedError  error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := Get(tc.Client, tc.Namespace, tc.Configmap, true)

			assert.Equal(t, tc.ExpectedString, actualString)
			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	cm := mock.ConfigMap("test", "test", map[string]string{"n": "v"})
	cl := fake.NewSimpleClientset(cm)

	validate(t, &testCase{
		Name:           "Should find test.test",
		Client:         client.NewClient(cl),
		Namespace:      "test",
		Configmap:      "test",
		ExpectedString: "##### CONFIGMAP - test #####\nexport n=\"v\"\n",
	})

	validate(t, &testCase{
		Name:          "Should not find test.test1",
		Client:        client.NewClient(cl),
		Namespace:     "test",
		Configmap:     "test1",
		ExpectedError: ErrMissingResource,
	})

	validate(t, &testCase{
		Name:          "Should not find test2.test",
		Client:        client.NewClient(cl),
		Namespace:     "test2",
		Configmap:     "test",
		ExpectedError: ErrMissingResource,
	})
}
