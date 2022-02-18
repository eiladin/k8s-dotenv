package cronjob

import (
	"bytes"
	"errors"
	"testing"

	"github.com/eiladin/k8s-dotenv/internal/options"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/batch/v1"
	v1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

type CronjobCmdSuite struct {
	suite.Suite
}

func mockContainer(env map[string]string, configmaps, secrets []string) corev1.Container {
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

	return c
}

func mockCronjobv1beta1(name, namespace string, env map[string]string, configmaps, secrets []string) *v1beta1.CronJob {
	res := &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: map[string]string{},
		},
	}

	res.Spec.JobTemplate.Spec.Template.Spec.Containers = []corev1.Container{mockContainer(env, configmaps, secrets)}
	return res
}

func mockCronjobv1(name, namespace string, env map[string]string, configmaps, secrets []string) *v1.CronJob {
	res := &v1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: map[string]string{},
		},
	}

	res.Spec.JobTemplate.Spec.Template.Spec.Containers = []corev1.Container{mockContainer(env, configmaps, secrets)}
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

func (suite CronjobCmdSuite) TestNewCmd() {
	got := NewCmd(options.NewOptions())
	suite.NotNil(got)
}

func (suite CronjobCmdSuite) TestValidArgs() {
	cases := []struct {
		group string
	}{
		{group: "batch/v1beta1"},
		{group: "batch/v1"},
	}

	for _, c := range cases {
		opt := options.NewOptions()
		client := fake.NewSimpleClientset()

		client.Fake.Resources = append(client.Fake.Resources, &metav1.APIResourceList{
			GroupVersion: c.group,
			APIResources: []metav1.APIResource{
				{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: c.group},
			},
		})

		opt.Name = "test"
		opt.Namespace = "test"
		opt.Client = client
		cmd := NewCmd(opt)
		got, _ := cmd.ValidArgsFunction(cmd, []string{}, "")
		suite.NotNil(got)
	}
}

func (suite CronjobCmdSuite) TestRun() {
	cases := []struct {
		group      string
		name       string
		namespace  string
		env        map[string]string
		configmaps []string
		secrets    []string
		args       []string
		shouldErr  bool
		testApiErr bool
	}{
		{group: "batch/v1", args: []string{"my-job"}, name: "my-job", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{group: "batch/v1", args: []string{"my-job"}, shouldErr: true},
		{group: "batch/v1", shouldErr: true},
		{group: "batch/v1beta1", args: []string{"my-job"}, name: "my-job", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{group: "batch/v1beta1", args: []string{"my-job"}, shouldErr: true},
		{group: "batch/v1beta1", shouldErr: true},
		{group: "batch/v1beta1", args: []string{"my-job"}, shouldErr: true},
		{group: "batch/v1", args: []string{"my-job"}, shouldErr: true, testApiErr: true},
	}

	for _, c := range cases {
		mocks := []runtime.Object{}
		if c.group == "batch/v1" {
			mocks = append(mocks, mockCronjobv1(c.name, c.namespace, c.env, c.configmaps, c.secrets))
		} else {
			mocks = append(mocks, mockCronjobv1beta1(c.name, c.namespace, c.env, c.configmaps, c.secrets))
		}
		for _, cm := range c.configmaps {
			mocks = append(mocks, mockConfigMap(cm, c.namespace, map[string]string{"config": "value"}))
		}
		for _, s := range c.secrets {
			mocks = append(mocks, mockSecret(s, c.namespace, map[string][]byte{"secret": []byte("value")}))
		}

		client := fake.NewSimpleClientset(mocks...)

		if c.testApiErr {
			client.Fake.Resources = append(client.Fake.Resources, &metav1.APIResourceList{
				GroupVersion: "a/b/c",
				APIResources: []metav1.APIResource{
					{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: c.group},
				},
			})
		} else {
			client.Fake.Resources = append(client.Fake.Resources, &metav1.APIResourceList{
				GroupVersion: c.group,
				APIResources: []metav1.APIResource{
					{Name: "CronJob", SingularName: "CronJob", Kind: "CronJob", Namespaced: true, Group: c.group},
				},
			})
		}

		if c.shouldErr {
			client.PrependReactor("get", "cronjobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				if c.group == "batch/v1" {
					return true, &v1.CronJob{}, errors.New("error getting cronjob")
				} else {
					return true, &v1beta1.CronJob{}, errors.New("error getting cronjob")
				}
			})
		}

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

func TestCronjobCmdSuite(t *testing.T) {
	suite.Run(t, new(CronjobCmdSuite))
}
