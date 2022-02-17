package v1

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/internal/options"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

type DaemonSetSuite struct {
	suite.Suite
}

func mockDaemonSet(name, namespace string, env map[string]string, configmaps, secrets []string) *v1.DaemonSet {
	res := &v1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: map[string]string{},
		},
	}

	containers := []corev1.Container{}
	c := corev1.Container{}

	for k, v := range env {
		c.Env = append(c.Env, corev1.EnvVar{Name: k, Value: v})
	}

	for _, cm := range configmaps {
		c.EnvFrom = append(c.EnvFrom, corev1.EnvFromSource{ConfigMapRef: &corev1.ConfigMapEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: cm}}})
	}
	for _, s := range secrets {
		c.EnvFrom = append(c.EnvFrom, corev1.EnvFromSource{SecretRef: &corev1.SecretEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: s}}})
	}

	containers = append(containers, c)
	res.Spec.Template.Spec.Containers = containers
	return res
}

func (suite DaemonSetSuite) TestDaemonSet() {
	cases := []struct {
		name       string
		namespace  string
		env        map[string]string
		configmaps []string
		secrets    []string
	}{
		{name: "my-daemonset", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-daemonset", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{}},
		{name: "my-daemonset", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-daemonset", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{}, secrets: []string{}},
		{name: "my-daemonset", namespace: "test", configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-daemonset", namespace: "test", configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{}},
		{name: "my-daemonset", namespace: "test", configmaps: []string{}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-daemonset", namespace: "test", configmaps: []string{}, secrets: []string{}},
	}

	for _, c := range cases {
		m := mockDaemonSet(c.name, c.namespace, c.env, c.configmaps, c.secrets)
		opt := options.NewOptions()
		opt.Client = fake.NewSimpleClientset(m)
		opt.Namespace = c.namespace
		opt.Name = c.name

		got, err := DaemonSet(opt)
		suite.NoError(err)
		suite.NotNil(got)
		suite.Len(got.Environment, len(c.env))
		suite.Len(got.ConfigMaps, len(c.configmaps))
		suite.Len(got.Secrets, len(c.secrets))
	}
}

func (suite DaemonSetSuite) TestDaemonSets() {
	type item struct {
		name      string
		namespace string
	}

	cases := []struct {
		namespace     string
		items         []item
		expectedCount int
	}{
		{
			namespace:     "test",
			items:         []item{{name: "my-daemonset", namespace: "test"}},
			expectedCount: 1,
		},
		{
			namespace:     "test",
			items:         []item{{name: "my-daemonset", namespace: "test"}, {name: "my-daemonset-2", namespace: "test"}},
			expectedCount: 2,
		},
		{
			namespace:     "other",
			items:         []item{{name: "my-daemonset", namespace: "test"}, {name: "my-daemonset-2", namespace: "test"}},
			expectedCount: 0,
		},
		{
			namespace:     "test",
			items:         []item{{name: "my-daemonset", namespace: "test"}, {name: "my-daemonset-2", namespace: "other"}},
			expectedCount: 1,
		},
	}

	for _, c := range cases {
		mocks := []runtime.Object{}
		for _, item := range c.items {
			mock := mockDaemonSet(item.name, item.namespace, nil, nil, nil)
			mocks = append(mocks, mock)
		}
		opt := options.NewOptions()
		opt.Client = fake.NewSimpleClientset(mocks...)
		opt.Namespace = c.namespace

		got, err := DaemonSets(opt)
		suite.NoError(err)
		suite.NotNil(got)
		suite.Len(got, c.expectedCount)
	}
}

func TestDaemonSetSuite(t *testing.T) {
	suite.Run(t, new(DaemonSetSuite))
}
