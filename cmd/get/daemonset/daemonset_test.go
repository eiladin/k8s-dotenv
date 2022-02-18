package daemonset

import (
	"bytes"
	"testing"

	"github.com/eiladin/k8s-dotenv/internal/options"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

type DaemonsetCmdSuite struct {
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

func mockSecret(name string, namespace string, data map[string][]byte) *corev1.Secret {
	res := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
	return res
}

func mockConfigMap(name string, namespace string, data map[string]string) *corev1.ConfigMap {
	res := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
	return res
}

func (suite DaemonsetCmdSuite) TestNewCmd() {
	got := NewCmd(options.NewOptions())
	suite.NotNil(got)
}

func (suite DaemonsetCmdSuite) TestValidArgs() {
	opt := options.NewOptions()
	client := fake.NewSimpleClientset()
	opt.Name = "test"
	opt.Namespace = "test"
	opt.Client = client
	cmd := NewCmd(opt)
	got, _ := cmd.ValidArgsFunction(cmd, []string{}, "")
	suite.NotNil(got)
}

func (suite DaemonsetCmdSuite) TestRun() {
	cases := []struct {
		name       string
		namespace  string
		env        map[string]string
		configmaps []string
		secrets    []string
		args       []string
		shouldErr  bool
	}{
		{args: []string{"my-job"}, name: "my-job", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{args: []string{"my-job"}, shouldErr: true},
		{shouldErr: true},
	}

	for _, c := range cases {
		mocks := []runtime.Object{}
		mocks = append(mocks, mockDaemonSet(c.name, c.namespace, c.env, c.configmaps, c.secrets))
		for _, cm := range c.configmaps {
			mocks = append(mocks, mockConfigMap(cm, c.namespace, map[string]string{"config": "value"}))
		}
		for _, s := range c.secrets {
			mocks = append(mocks, mockSecret(s, c.namespace, map[string][]byte{"secret": []byte("value")}))
		}

		client := fake.NewSimpleClientset(mocks...)

		opt := options.NewOptions()
		opt.Client = client
		opt.Namespace = c.namespace
		opt.Name = c.name

		var b bytes.Buffer
		err := opt.SetWriter(&b)
		suite.NoError(err)
		cmd := NewCmd(opt)
		err = cmd.RunE(cmd, c.args)

		if c.shouldErr {
			suite.Error(err)
		} else {
			suite.NoError(err)
			got := b.String()
			for k, v := range c.env {
				suite.Contains(got, k)
				suite.Contains(got, v)
			}
			for _, cm := range c.configmaps {
				suite.Contains(got, cm)
			}
			for _, s := range c.secrets {
				suite.Contains(got, s)
			}
		}
	}
}

func TestDaemonsetCmdSuite(t *testing.T) {
	suite.Run(t, new(DaemonsetCmdSuite))
}
