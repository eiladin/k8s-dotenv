package secret

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGet(t *testing.T) {
	type testCase struct {
		Name string

		Client       *client.Client
		Namespace    string
		Secret       string
		ShouldExport bool

		ExpectedString string
		ErrorChecker   func(err error) bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := Get(tc.Client, tc.Namespace, tc.Secret, tc.ShouldExport)

			assert.Equal(t, tc.ExpectedString, actualString)
			if tc.ErrorChecker != nil {
				assert.Equal(t, true, tc.ErrorChecker(actualError))
			}
		})
	}

	cm := mock.Secret("test", "test", map[string][]byte{"n": []byte("v")})
	cl := fake.NewSimpleClientset(cm)
	validate(t, &testCase{Name: "Should find test.test", Secret: "test", Client: client.NewClient(cl), Namespace: "test", ExpectedString: "##### SECRET - test #####\nn=\"v\"\n"})
	validate(t, &testCase{Name: "Should not find test.test1", Secret: "test1", Client: client.NewClient(cl), Namespace: "test", ErrorChecker: errors.IsNotFound})
	validate(t, &testCase{Name: "Should not find test2.test", Secret: "test", Client: client.NewClient(cl), Namespace: "test2", ErrorChecker: errors.IsNotFound})
}
