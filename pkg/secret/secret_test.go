package secret

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGet(t *testing.T) {
	type testCase struct {
		Name string

		Opt    *options.Options
		Secret string

		ExpectedString string
		ErrorChecker   func(err error) bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := Get(tc.Opt, tc.Secret)

			assert.Equal(t, tc.ExpectedString, actualString)
			if tc.ErrorChecker != nil {
				assert.Equal(t, true, tc.ErrorChecker(actualError))
			}
		})
	}

	cm := mock.Secret("test", "test", map[string][]byte{"n": []byte("v")})
	client := fake.NewSimpleClientset(cm)
	validate(t, &testCase{Name: "Should find test.test", Secret: "test", Opt: &options.Options{Client: client, Namespace: "test"}, ExpectedString: "##### SECRET - test #####\nexport n=\"v\"\n"})
	validate(t, &testCase{Name: "Should not find test.test1", Secret: "test1", Opt: &options.Options{Client: client, Namespace: "test"}, ErrorChecker: errors.IsNotFound})
	validate(t, &testCase{Name: "Should not find test2.test", Secret: "test", Opt: &options.Options{Client: client, Namespace: "test2"}, ErrorChecker: errors.IsNotFound})
}
