package v1

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/internal/options"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

type NamespaceSuite struct {
	suite.Suite
}

func mockNamespace(name string) *corev1.Namespace {
	res := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: map[string]string{},
		},
	}
	return res
}

func (suite NamespaceSuite) TestNamespaces() {
	cases := []struct {
		items         []string
		expectedCount int
	}{
		{
			items:         []string{"test"},
			expectedCount: 1,
		},
		{
			items:         []string{"test", "test2"},
			expectedCount: 2,
		},
	}

	for _, c := range cases {
		mocks := []runtime.Object{}
		for _, item := range c.items {
			mock := mockNamespace(item)
			mocks = append(mocks, mock)
		}
		opt := options.NewOptions()
		opt.Client = fake.NewSimpleClientset(mocks...)

		got, err := Namespaces(opt)
		suite.NoError(err)
		suite.NotNil(got)
		suite.Len(got, c.expectedCount)
	}
}

func TestNamespaceSuite(t *testing.T) {
	suite.Run(t, new(NamespaceSuite))
}
