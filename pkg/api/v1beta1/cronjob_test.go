package v1beta1

import (
	"errors"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mocks"
	"github.com/stretchr/testify/suite"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

type CronJobSuite struct {
	suite.Suite
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
		m := mocks.CronJobv1beta1(c.name, c.namespace, c.env, c.configmaps, c.secrets)
		client := fake.NewSimpleClientset(m)

		if c.shouldErr {
			client.PrependReactor("get", "cronjobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &batchv1beta1.CronJob{}, errors.New("error getting cronjob")
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
	type cronjob struct {
		name      string
		namespace string
	}

	cases := []struct {
		namespace     string
		cronjobs      []cronjob
		expectedCount int
		shouldErr     bool
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
		{
			namespace: "test",
			shouldErr: true,
		},
	}

	for _, c := range cases {
		ms := []runtime.Object{}
		for _, cj := range c.cronjobs {
			mock := mocks.CronJobv1beta1(cj.name, cj.namespace, nil, nil, nil)
			ms = append(ms, mock)
		}
		client := fake.NewSimpleClientset(ms...)
		if c.shouldErr {
			client.PrependReactor("list", "cronjobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &batchv1beta1.CronJobList{}, errors.New("error getting cronjob list")
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
