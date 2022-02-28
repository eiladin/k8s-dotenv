package job

import (
	"bytes"
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

type JobCmdSuite struct {
	suite.Suite
}

func (suite JobCmdSuite) TestNewCmd() {
	got := NewCmd(nil)
	suite.NotNil(got)
}

func (suite JobCmdSuite) TestValidArgs() {
	opt := &options.Options{
		Client:    fake.NewSimpleClientset(),
		Namespace: "test",
		Name:      "test",
	}
	cmd := NewCmd(opt)
	got, _ := cmd.ValidArgsFunction(cmd, []string{}, "")
	suite.NotNil(got)
}

func (suite JobCmdSuite) TestRun() {
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
		ms := []runtime.Object{}
		ms = append(ms, mocks.Job(c.name, c.namespace, c.env, c.configmaps, c.secrets))
		for _, cm := range c.configmaps {
			ms = append(ms, mocks.ConfigMap(cm, c.namespace, map[string]string{"config": "value"}))
		}
		for _, s := range c.secrets {
			ms = append(ms, mocks.Secret(s, c.namespace, map[string][]byte{"secret": []byte("value")}))
		}

		client := fake.NewSimpleClientset(ms...)
		if c.shouldErr {
			client.PrependReactor("get", "jobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &batchv1.Job{}, errors.New("error getting job")
			})
		}

		var b bytes.Buffer
		opt := &options.Options{
			Client:    client,
			Namespace: c.namespace,
			Name:      c.name,
			Writer:    &b,
		}

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

func TestJobCmdSuite(t *testing.T) {
	suite.Run(t, new(JobCmdSuite))
}
