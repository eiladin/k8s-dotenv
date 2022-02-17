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
	cases := []struct {
		name      string
		namespace string
		shouldErr bool
	}{
		{name: "test", namespace: "test"},
		{name: "test1", namespace: "test", shouldErr: true},
		{name: "test", namespace: "test2", shouldErr: true},
	}

	for _, c := range cases {
		cm := mockConfigMap("test", "test", map[string]string{"n": "v"})
		opt := options.NewOptions()
		opt.Client = fake.NewSimpleClientset(cm)
		opt.Namespace = c.namespace

		got, err := Get(opt, c.name)
		if c.shouldErr {
			suite.Error(err)
		} else {
			suite.NoError(err)
			suite.Greater(len(got), 0, "result should have a length greater than 0")
		}
	}
}

func TestConfigMapSuite(t *testing.T) {
	suite.Run(t, new(ConfigMapSuite))
}
