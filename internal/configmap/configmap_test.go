package configmap

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/internal/options"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type ConfigMapSuite struct {
	suite.Suite
}

func mockConfigMap(name string, namespace string, data map[string]string) *v1.ConfigMap {
	res := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
	return res
}

func (suite ConfigMapSuite) TestGet() {
	cm := mockConfigMap("test", "test", map[string]string{"n": "v"})
	opt := options.NewOptions()
	opt.Client = fake.NewSimpleClientset(cm)
	opt.Namespace = "test"

	res, err := Get(opt, "test")
	suite.NoError(err)
	suite.Greater(len(res), 0, "result should have a length greater than 0")
}

func TestConfigMapSuite(t *testing.T) {
	suite.Run(t, new(ConfigMapSuite))
}
