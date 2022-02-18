package secret

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/internal/options"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type SecretSuite struct {
	suite.Suite
}

func mockSecret(name string, namespace string, data map[string][]byte) *v1.Secret {
	res := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
	return res
}

func (suite SecretSuite) TestGet() {
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
		s := mockSecret("test", "test", map[string][]byte{"n": []byte("v")})
		opt := options.NewOptions()
		opt.Client = fake.NewSimpleClientset(s)
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

func TestSecretSuite(t *testing.T) {
	suite.Run(t, new(SecretSuite))
}
