package v1beta1

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/internal/options"
	"github.com/stretchr/testify/suite"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

type CronJobSuite struct {
	suite.Suite
}

func mockCronJob(name, namespace string, env map[string]string, configmaps, secrets []string) *v1beta1.CronJob {
	res := &v1beta1.CronJob{
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
	res.Spec.JobTemplate.Spec.Template.Spec.Containers = containers
	return res
}

func (suite CronJobSuite) TestCronJob() {
	cases := []struct {
		name       string
		namespace  string
		env        map[string]string
		configmaps []string
		secrets    []string
	}{
		{name: "my-cronjob", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-cronjob", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{}},
		{name: "my-cronjob", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-cronjob", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{}, secrets: []string{}},
		{name: "my-cronjob", namespace: "test", configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-cronjob", namespace: "test", configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{}},
		{name: "my-cronjob", namespace: "test", configmaps: []string{}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-cronjob", namespace: "test", configmaps: []string{}, secrets: []string{}},
	}

	for _, c := range cases {
		m := mockCronJob(c.name, c.namespace, c.env, c.configmaps, c.secrets)
		opt := options.NewOptions()
		opt.Client = fake.NewSimpleClientset(m)
		opt.Namespace = c.namespace
		opt.Name = c.name

		res, err := CronJob(opt)
		suite.NoError(err)
		suite.NotNil(res)
		suite.Len(res.Environment, len(c.env))
		suite.Len(res.ConfigMaps, len(c.configmaps))
		suite.Len(res.Secrets, len(c.secrets))
	}
}

func (suite CronJobSuite) TestCronJobs() {
	type cronjob struct {
		name      string
		namespace string
	}

	cases := []struct {
		namespace     string
		cronjobs      []cronjob
		expectedCount int
	}{
		{
			namespace:     "test",
			cronjobs:      []cronjob{{name: "my-cronjob", namespace: "test"}},
			expectedCount: 1,
		},
		{
			namespace:     "test",
			cronjobs:      []cronjob{{name: "my-cronjob", namespace: "test"}, {name: "my-cronjob-2", namespace: "test"}},
			expectedCount: 2,
		},
		{
			namespace:     "other",
			cronjobs:      []cronjob{{name: "my-cronjob", namespace: "test"}, {name: "my-cronjob-2", namespace: "test"}},
			expectedCount: 0,
		},
		{
			namespace:     "test",
			cronjobs:      []cronjob{{name: "my-cronjob", namespace: "test"}, {name: "my-cronjob-2", namespace: "other"}},
			expectedCount: 1,
		},
	}

	for _, c := range cases {
		mocks := []runtime.Object{}
		for _, cj := range c.cronjobs {
			mock := mockCronJob(cj.name, cj.namespace, nil, nil, nil)
			mocks = append(mocks, mock)
		}
		client := fake.NewSimpleClientset(mocks...)

		opt := options.NewOptions()
		opt.Client = client
		opt.Namespace = c.namespace

		res, err := CronJobs(opt)
		suite.NoError(err)
		suite.NotNil(res)
		suite.Len(res, c.expectedCount)
	}
}

func TestCronJobSuite(t *testing.T) {
	suite.Run(t, new(CronJobSuite))
}
