package v1

import (
	"errors"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mocks"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

type NamespaceSuite struct {
	suite.Suite
}

func (suite NamespaceSuite) TestNamespaces() {
	cases := []struct {
		items         []string
		expectedCount int
		shouldErr     bool
	}{
		{
			items:         []string{"test"},
			expectedCount: 1,
		},
		{
			items:         []string{"test", "test2"},
			expectedCount: 2,
		},
		{
			shouldErr: true,
		},
	}

	for _, c := range cases {
		ms := []runtime.Object{}
		for _, item := range c.items {
			mock := mocks.Namespace(item)
			ms = append(ms, mock)
		}
		client := fake.NewSimpleClientset(ms...)
		if c.shouldErr {
			client.PrependReactor("list", "namespaces", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &corev1.NamespaceList{}, errors.New("error getting namespaces")
			})
		}

		opt := &options.Options{Client: client}

		got, err := Namespaces(opt)
		if c.shouldErr {
			suite.Error(err)
		} else {
			suite.NoError(err)
			suite.NotNil(got)
			suite.Len(got, c.expectedCount)
		}
	}
}

func TestNamespaceSuite(t *testing.T) {
	suite.Run(t, new(NamespaceSuite))
}
