package secret

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mocks"
	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/kubernetes/fake"
)

type SecretSuite struct {
	suite.Suite
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
		s := mocks.Secret("test", "test", map[string][]byte{"n": []byte("v")})

		opt := &options.Options{
			Client:    fake.NewSimpleClientset(s),
			Namespace: c.namespace,
		}

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
