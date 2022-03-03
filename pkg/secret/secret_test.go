package secret

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	type testCase struct {
		Name           string
		Client         *client.Client
		Namespace      string
		Secret         string
		ExpectedString string
		ShouldExport   bool
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := Get(tc.Client, tc.Namespace, tc.Secret, tc.ShouldExport)

			assert.Equal(t, tc.ExpectedString, actualString)

			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	cm := mock.Secret("test", "test", map[string][]byte{"n": []byte("v")})
	cl := mock.NewFakeClient(cm)

	validate(t, &testCase{
		Name:   "Should find test.test",
		Secret: "test", Client: client.NewClient(cl),
		Namespace:      "test",
		ExpectedString: "##### SECRET - test #####\nn=\"v\"\n",
	})

	validate(t, &testCase{
		Name:        "Should not find test.test1",
		Secret:      "test1",
		Client:      client.NewClient(cl),
		Namespace:   "test",
		ExpectError: true,
	})

	validate(t, &testCase{
		Name:        "Should not find test2.test",
		Secret:      "test",
		Client:      client.NewClient(cl),
		Namespace:   "test2",
		ExpectError: true,
	})
}
