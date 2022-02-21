package cronjob

import (
	"bytes"
	"errors"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mocks"
	"github.com/stretchr/testify/suite"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

type CronjobCmdSuite struct {
	suite.Suite
}

func (suite CronjobCmdSuite) TestNewCmd() {
	got := NewCmd(options.NewOptions())
	suite.NotNil(got)
}

func (suite CronjobCmdSuite) TestValidArgs() {
	cases := []struct {
		group         string
		shouldBeEmpty bool
	}{
		{group: "batch/v1beta1"},
		{group: "batch/v1"},
		{group: "batch/not-a-version", shouldBeEmpty: true},
	}

	for _, c := range cases {
		opt := options.NewOptions()

		var mock runtime.Object
		if c.group == "batch/v1" {
			mock = mocks.CronJobv1("my-cronjob", "test", nil, nil, nil)
		} else if c.group == "batch/v1beta1" {
			mock = mocks.CronJobv1beta1("my-cronjob", "test", nil, nil, nil)
		}

		var client *fake.Clientset
		if mock != nil {
			client = fake.NewSimpleClientset(mock)
		} else {
			client = fake.NewSimpleClientset()
		}

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
		if c.shouldBeEmpty {
			suite.Empty(got)
		} else {
			suite.NotEmpty(got)
		}
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
		{group: "batch/v1", args: []string{"my-cronjob"}, name: "my-cronjob", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{group: "batch/v1", args: []string{"my-cronjob"}, shouldErr: true},
		{group: "batch/v1", shouldErr: true},
		{group: "batch/v1beta1", args: []string{"my-cronjob"}, name: "my-cronjob", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{group: "batch/v1beta1", args: []string{"my-cronjob"}, shouldErr: true},
		{group: "batch/v1beta1", shouldErr: true},
		{group: "batch/v1beta1", args: []string{"my-cronjob"}, shouldErr: true},
		{group: "batch/v1", args: []string{"my-cronjob"}, shouldErr: true, testApiErr: true},
		{group: "batch/not-a-version", args: []string{"my-cronjob"}, shouldErr: true},
	}

	for _, c := range cases {
		ms := []runtime.Object{}
		if c.group == "batch/v1" {
			ms = append(ms, mocks.CronJobv1(c.name, c.namespace, c.env, c.configmaps, c.secrets))
		} else if c.group == "batch/v1beta1" {
			ms = append(ms, mocks.CronJobv1beta1(c.name, c.namespace, c.env, c.configmaps, c.secrets))
		}
		for _, cm := range c.configmaps {
			ms = append(ms, mocks.ConfigMap(cm, c.namespace, map[string]string{"config": "value"}))
		}
		for _, s := range c.secrets {
			ms = append(ms, mocks.Secret(s, c.namespace, map[string][]byte{"secret": []byte("value")}))
		}

		client := fake.NewSimpleClientset(ms...)

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
					return true, &batchv1.CronJob{}, errors.New("error getting cronjob")
				} else {
					return true, &batchv1beta1.CronJob{}, errors.New("error getting cronjob")
				}
			})
		}

		var b bytes.Buffer
		opt := options.NewOptions()
		opt.Client = client
		opt.Namespace = c.namespace
		opt.Name = c.name
		opt.FileWriter = &b

		cmd := NewCmd(opt)
		err := cmd.RunE(cmd, c.args)

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
