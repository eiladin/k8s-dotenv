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

type JobSuite struct {
	suite.Suite
}

func mockJob(name, namespace string, env map[string]string, configmaps, secrets []string) *v1.Job {
	res := &v1.Job{
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

func (suite JobSuite) TestJob() {
	cases := []struct {
		name       string
		namespace  string
		env        map[string]string
		configmaps []string
		secrets    []string
		shouldErr  bool
	}{
		{name: "my-job", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-job", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{}},
		{name: "my-job", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-job", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{}, secrets: []string{}},
		{name: "my-job", namespace: "test", configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-job", namespace: "test", configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{}},
		{name: "my-job", namespace: "test", configmaps: []string{}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-job", namespace: "test", configmaps: []string{}, secrets: []string{}},
		{shouldErr: true},
	}

	for _, c := range cases {
		m := mockJob(c.name, c.namespace, c.env, c.configmaps, c.secrets)
		client := fake.NewSimpleClientset(m)
		if c.shouldErr {
			client.PrependReactor("get", "jobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &v1.Job{}, errors.New("error getting job")
			})
		}

		opt := options.NewOptions()
		opt.Client = client
		opt.Namespace = c.namespace
		opt.Name = c.name

		got, err := Job(opt)
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

func (suite JobSuite) TestJobs() {
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
			items:         []item{{name: "my-job", namespace: "test"}},
			expectedCount: 1,
		},
		{
			namespace:     "test",
			items:         []item{{name: "my-job", namespace: "test"}, {name: "my-job-2", namespace: "test"}},
			expectedCount: 2,
		},
		{
			namespace:     "other",
			items:         []item{{name: "my-job", namespace: "test"}, {name: "my-job-2", namespace: "test"}},
			expectedCount: 0,
		},
		{
			namespace:     "test",
			items:         []item{{name: "my-job", namespace: "test"}, {name: "my-job-2", namespace: "other"}},
			expectedCount: 1,
		},
		{
			shouldErr: true,
		},
	}

	for _, c := range cases {
		mocks := []runtime.Object{}
		for _, item := range c.items {
			mock := mockJob(item.name, item.namespace, nil, nil, nil)
			mocks = append(mocks, mock)
		}
		client := fake.NewSimpleClientset(mocks...)
		if c.shouldErr {
			client.PrependReactor("list", "jobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &v1.JobList{}, errors.New("error getting jobs list")
			})
		}

		opt := options.NewOptions()
		opt.Client = client
		opt.Namespace = c.namespace

		got, err := Jobs(opt)
		if c.shouldErr {
			suite.Error(err)
		} else {
			suite.NoError(err)
			suite.NotNil(got)
			suite.Len(got, c.expectedCount)
		}
	}
}

func TestJobSuite(t *testing.T) {
	suite.Run(t, new(JobSuite))
}
