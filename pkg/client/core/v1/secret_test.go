package corev1

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestSecretV1(t *testing.T) {
	type testCase struct {
		Name           string
		Client         *CoreV1
		Secret         string
		ExpectedString string
		ShouldExport   bool
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := tc.Client.Secret(tc.Secret, tc.ShouldExport)

			assert.Equal(t, tc.ExpectedString, actualString)

			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	cm := mock.Secret("test", "test", map[string][]byte{"n": []byte("v")})
	kubeClient := mock.NewFakeClient(cm)

	validate(t, &testCase{
		Name:           "Should find test.test",
		Secret:         "test",
		Client:         NewCoreV1(kubeClient.CoreV1(), "test"),
		ExpectedString: "##### SECRET - test #####\nn=\"v\"\n",
	})

	validate(t, &testCase{
		Name:        "Should not find test.test1",
		Secret:      "test1",
		Client:      NewCoreV1(kubeClient.CoreV1(), "test"),
		ExpectError: true,
	})

	validate(t, &testCase{
		Name:        "Should not find test2.test",
		Secret:      "test",
		Client:      NewCoreV1(kubeClient.CoreV1(), "test2"),
		ExpectError: true,
	})
}
