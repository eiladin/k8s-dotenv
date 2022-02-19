package environment

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mocks"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

type EnvironmentSuite struct {
	suite.Suite
}

func (suite EnvironmentSuite) TestOutput() {
	cases := []struct {
		env           map[string]string
		configmapName string
		configmap     map[string]string
		secretName    string
		secrets       map[string][]byte
		shouldErr     bool
	}{
		{
			env:       map[string]string{"env1": "val", "env2": "val2"},
			configmap: map[string]string{"config": "val", "config2": "val2"},
			secrets:   map[string][]byte{"secret": []byte("val"), "secret2": []byte("val2")},
		},
		{
			configmap: map[string]string{"config": "val", "config2": "val2"},
			secrets:   map[string][]byte{"secret": []byte("val"), "secret2": []byte("val2")},
		},
		{
			env:     map[string]string{"env1": "val", "env2": "val2"},
			secrets: map[string][]byte{"secret": []byte("val"), "secret2": []byte("val2")},
		},
		{
			env:       map[string]string{"env1": "val", "env2": "val2"},
			configmap: map[string]string{"config": "val", "config2": "val2"},
		},
		{
			configmap:     map[string]string{"config": "val"},
			configmapName: "test1",
			shouldErr:     true,
		},
		{
			secrets:    map[string][]byte{"secret": []byte("val")},
			secretName: "test1",
			shouldErr:  true,
		},
	}

	buildList := func(name string) []string {
		if name != "" {
			return []string{name}
		}
		return []string{"test"}
	}

	for i, c := range cases {
		r := NewResult()
		if c.configmap != nil {
			r.ConfigMaps = buildList(c.configmapName)
		}
		if c.secrets != nil {
			r.Secrets = buildList(c.secretName)
		}
		if c.env != nil {
			r.Environment = c.env
		}

		objs := []runtime.Object{}
		if c.configmap != nil {
			objs = append(objs, mocks.ConfigMap("test", "test", c.configmap))
		}
		if c.secrets != nil {
			objs = append(objs, mocks.Secret("test", "test", c.secrets))
		}

		opt := options.NewOptions()
		opt.Client = fake.NewSimpleClientset(objs...)
		opt.Namespace = "test"

		got, err := r.Output(opt)
		caseDesc := fmt.Sprintf("Test case %d", i)
		if c.shouldErr {
			suite.Error(err, caseDesc)
		} else {
			suite.NoError(err, caseDesc)
			suite.NotNil(got, caseDesc)
			for k, v := range c.configmap {
				suite.Contains(got, fmt.Sprintf("export %s=\"%s\"", k, v), caseDesc)
			}
			for k, v := range c.secrets {
				suite.Contains(got, fmt.Sprintf("export %s=\"%s\"", k, v), caseDesc)
			}
			for k, v := range c.env {
				suite.Contains(got, fmt.Sprintf("export %s=\"%s\"", k, v), caseDesc)
			}
		}
	}
}

func (suite EnvironmentSuite) TestWrite() {
	cases := []struct {
		env           map[string]string
		configmapName string
		configmap     map[string]string
		secretName    string
		secrets       map[string][]byte
		shouldErr     bool
		filename      string
		useFileWriter bool
	}{
		{
			env:       map[string]string{"env1": "val", "env2": "val2"},
			configmap: map[string]string{"config": "val", "config2": "val2"},
			secrets:   map[string][]byte{"secret": []byte("val"), "secret2": []byte("val2")},
		},
		{
			configmap: map[string]string{"config": "val", "config2": "val2"},
			secrets:   map[string][]byte{"secret": []byte("val"), "secret2": []byte("val2")},
		},
		{
			env:     map[string]string{"env1": "val", "env2": "val2"},
			secrets: map[string][]byte{"secret": []byte("val"), "secret2": []byte("val2")},
		},
		{
			env:       map[string]string{"env1": "val", "env2": "val2"},
			configmap: map[string]string{"config": "val", "config2": "val2"},
		},
		{
			configmapName: "test2",
			configmap:     map[string]string{"config": "val", "config2": "val2"},
			shouldErr:     true,
		},
		{
			env:           map[string]string{"env1": "val", "env2": "val2"},
			configmap:     map[string]string{"config": "val", "config2": "val2"},
			useFileWriter: true,
			filename:      "test.out",
		},
		{
			env:           map[string]string{"env1": "val", "env2": "val2"},
			configmap:     map[string]string{"config": "val", "config2": "val2"},
			useFileWriter: true,
			shouldErr:     true,
		},
	}

	buildList := func(name string) []string {
		if name != "" {
			return []string{name}
		}
		return []string{"test"}
	}

	for i, c := range cases {
		r := NewResult()
		if c.configmap != nil {
			r.ConfigMaps = buildList(c.configmapName)
		}
		if c.secrets != nil {
			r.Secrets = buildList(c.secretName)
		}
		if c.env != nil {
			r.Environment = c.env
		}

		ms := []runtime.Object{}
		if c.configmap != nil {
			ms = append(ms, mocks.ConfigMap("test", "test", c.configmap))
		}
		if c.secrets != nil {
			ms = append(ms, mocks.Secret("test", "test", c.secrets))
		}

		var b bytes.Buffer
		var err error
		var got string
		opt := options.NewOptions()
		opt.Client = fake.NewSimpleClientset(ms...)
		opt.Namespace = "test"
		if c.useFileWriter {
			opt.Filename = c.filename
			defer os.Remove(c.filename)
			err = r.Write(opt)
			var fileBytes []byte
			fileBytes, _ = os.ReadFile(c.filename)
			got = string(fileBytes)
		} else {
			opt.Writer = &b
			err = r.Write(opt)
			got = b.String()
		}

		caseDesc := fmt.Sprintf("Test case %d", i)
		if c.shouldErr {
			suite.Error(err)
		} else {
			suite.NoError(err, caseDesc)
			suite.NotNil(got, caseDesc)
			for k, v := range c.configmap {
				suite.Contains(got, fmt.Sprintf("export %s=\"%s\"", k, v), caseDesc)
			}
			for k, v := range c.secrets {
				suite.Contains(got, fmt.Sprintf("export %s=\"%s\"", k, v), caseDesc)
			}
			for k, v := range c.env {
				suite.Contains(got, fmt.Sprintf("export %s=\"%s\"", k, v), caseDesc)
			}
		}
	}
}

func TestEnvironmentSuite(t *testing.T) {
	suite.Run(t, new(EnvironmentSuite))
}
