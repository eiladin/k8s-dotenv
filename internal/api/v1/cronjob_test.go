package v1

import (
	"errors"
	"testing"

	"github.com/eiladin/k8s-dotenv/internal/options"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

type CronJobSuite struct {
	suite.Suite
}

func mockCronJob(name, namespace string, env map[string]string, configmaps, secrets []string) *v1.CronJob {
	cronjob := &v1.CronJob{
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
	cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers = containers
	return cronjob
}

func (suite CronJobSuite) TestCronJob() {
	cases := []struct {
		name       string
		namespace  string
		env        map[string]string
		configmaps []string
		secrets    []string
		shouldErr  bool
	}{
		{name: "my-cronjob", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-cronjob", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{}},
		{name: "my-cronjob", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-cronjob", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{}, secrets: []string{}},
		{name: "my-cronjob", namespace: "test", configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-cronjob", namespace: "test", configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{}},
		{name: "my-cronjob", namespace: "test", configmaps: []string{}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-cronjob", namespace: "test", configmaps: []string{}, secrets: []string{}},
		{name: "my-cronjob", namespace: "test", shouldErr: true},
	}

	for _, c := range cases {
		m := mockCronJob(c.name, c.namespace, c.env, c.configmaps, c.secrets)
		client := fake.NewSimpleClientset(m)
		if c.shouldErr {
			client.PrependReactor("get", "cronjobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &v1.CronJob{}, errors.New("error getting cronjob")
			})
		}

		opt := options.NewOptions()
		opt.Client = client
		opt.Namespace = c.namespace
		opt.Name = c.name

		got, err := CronJob(opt)
		if c.shouldErr {
			suite.Error(err)
		} else {
			suite.NoError(err)
			suite.NotNil(got)
			suite.Len(got.Environment, len(c.env))
			suite.Len(got.ConfigMaps, len(c.configmaps))
			suite.Len(got.Secrets, len(c.secrets))
		}
	}
}

func (suite CronJobSuite) TestCronJobs() {
	type item struct {
		name      string
		namespace string
	}

	cases := []struct {
		namespace     string
		items         []item
		expectedCount int
		shouldErr     bool
	}{
		{
			namespace:     "test",
			items:         []item{{name: "my-cronjob", namespace: "test"}},
			expectedCount: 1,
		},
		{
			namespace:     "test",
			items:         []item{{name: "my-cronjob", namespace: "test"}, {name: "my-cronjob-2", namespace: "test"}},
			expectedCount: 2,
		},
		{
			namespace:     "other",
			items:         []item{{name: "my-cronjob", namespace: "test"}, {name: "my-cronjob-2", namespace: "test"}},
			expectedCount: 0,
		},
		{
			namespace:     "test",
			items:         []item{{name: "my-cronjob", namespace: "test"}, {name: "my-cronjob-2", namespace: "other"}},
			expectedCount: 1,
		},
		{
			namespace: "test",
			shouldErr: true,
		},
	}

	for _, c := range cases {
		mocks := []runtime.Object{}
		for _, item := range c.items {
			mock := mockCronJob(item.name, item.namespace, nil, nil, nil)
			mocks = append(mocks, mock)
		}
		client := fake.NewSimpleClientset(mocks...)
		if c.shouldErr {
			client.PrependReactor("list", "cronjobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &v1.CronJobList{}, errors.New("error getting cronjob")
			})
		}

		opt := options.NewOptions()
		opt.Client = client
		opt.Namespace = c.namespace

		got, err := CronJobs(opt)
		if c.shouldErr {
			suite.Error(err)
		} else {
			suite.NoError(err)
			suite.NotNil(got)
			suite.Len(got, c.expectedCount)
		}
	}
}

func TestCronJobSuite(t *testing.T) {
	suite.Run(t, new(CronJobSuite))
}
