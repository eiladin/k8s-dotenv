package corev1

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestConfigMapV1(t *testing.T) {
	type testCase struct {
		Name           string
		Client         *CoreV1
		Configmap      string
		ExpectedString string
		ShouldExport   bool
		ExpectError    bool
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString, actualError := tc.Client.ConfigMapV1(tc.Configmap, true)

			assert.Equal(t, tc.ExpectedString, actualString)
			if tc.ExpectError {
				assert.Error(t, actualError)
			} else {
				assert.NoError(t, actualError)
			}
		})
	}

	cm := mock.ConfigMap("test", "test", map[string]string{"n": "v"})
	kubeClient := mock.NewFakeClient(cm)

	validate(t, &testCase{
		Name:           "Should find test.test",
		Client:         NewCoreV1(kubeClient.CoreV1(), "test"),
		Configmap:      "test",
		ExpectedString: "##### CONFIGMAP - test #####\nexport n=\"v\"\n",
	})

	validate(t, &testCase{
		Name:        "Should not find test.test1",
		Client:         NewCoreV1(kubeClient.CoreV1(), "test"),
		Configmap:   "test1",
		ExpectError: true,
	})

	validate(t, &testCase{
		Name:        "Should not find test2.test",
		Client:         NewCoreV1(kubeClient.CoreV1(), "test2"),
		Configmap:   "test",
		ExpectError: true,
	})
}
