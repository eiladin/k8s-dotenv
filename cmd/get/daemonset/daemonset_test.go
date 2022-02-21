package daemonset

import (
	"bytes"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mocks"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

type DaemonsetCmdSuite struct {
	suite.Suite
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
		{args: []string{"my-daemonset"}, name: "my-daemonset", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{args: []string{"my-daemonset"}, shouldErr: true},
		{shouldErr: true},
	}

	for _, c := range cases {
		ms := []runtime.Object{}
		ms = append(ms, mocks.DaemonSet(c.name, c.namespace, c.env, c.configmaps, c.secrets))
		for _, cm := range c.configmaps {
			ms = append(ms, mocks.ConfigMap(cm, c.namespace, map[string]string{"config": "value"}))
		}
		for _, s := range c.secrets {
			ms = append(ms, mocks.Secret(s, c.namespace, map[string][]byte{"secret": []byte("value")}))
		}

		client := fake.NewSimpleClientset(ms...)

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

func TestDaemonsetCmdSuite(t *testing.T) {
	suite.Run(t, new(DaemonsetCmdSuite))
}
