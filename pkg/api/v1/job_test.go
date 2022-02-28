package v1

import (
	"errors"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mocks"
	"github.com/stretchr/testify/suite"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

type JobSuite struct {
	suite.Suite
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
		m := mocks.Job(c.name, c.namespace, c.env, c.configmaps, c.secrets)
		client := fake.NewSimpleClientset(m)
		if c.shouldErr {
			client.PrependReactor("get", "jobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &batchv1.Job{}, errors.New("error getting job")
			})
		}

		opt := &options.Options{
			Client:    client,
			Namespace: c.namespace,
			Name:      c.name,
		}

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
		ms := []runtime.Object{}
		for _, item := range c.items {
			mock := mocks.Job(item.name, item.namespace, nil, nil, nil)
			ms = append(ms, mock)
		}
		client := fake.NewSimpleClientset(ms...)
		if c.shouldErr {
			client.PrependReactor("list", "jobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &batchv1.JobList{}, errors.New("error getting jobs list")
			})
		}

		opt := &options.Options{
			Client:    client,
			Namespace: c.namespace,
		}

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
