package environment

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mocks"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

type EnvironmentSuite struct {
	suite.Suite
}

func (suite EnvironmentSuite) TestNewResult() {
	got := NewResult()
	suite.NotNil(got)
}

func (suite EnvironmentSuite) TestFromContainers() {
	containers := []v1.Container{
		mocks.Container(map[string]string{"a": "b"}, []string{"a"}, []string{"b"}),
	}

	got := FromContainers(containers)
	suite.NotNil(got)
	suite.NotEmpty(got.Environment)
	suite.NotEmpty(got.ConfigMaps)
	suite.NotEmpty(got.Secrets)
}

func (suite EnvironmentSuite) TestOutput() {
	envCase := map[string]string{"env1": "val", "env2": "val2"}
	configmapCase := map[string]string{"config": "val", "config2": "val2"}
	secretsCase := map[string][]byte{"secret": []byte("val"), "secret2": []byte("val2")}

	cases := []struct {
		env           map[string]string
		configmapName string
		configmap     map[string]string
		secretName    string
		secrets       map[string][]byte
		shouldErr     bool
	}{
		{env: envCase, configmap: configmapCase, secrets: secretsCase, configmapName: "test", secretName: "test"},
		{configmap: configmapCase, secrets: secretsCase, configmapName: "test", secretName: "test"},
		{env: envCase, secrets: secretsCase, secretName: "test"},
		{env: envCase, configmap: configmapCase, configmapName: "test"},
		{configmap: configmapCase, configmapName: "test1", shouldErr: true},
		{secrets: secretsCase, secretName: "test1", shouldErr: true},
	}

	for i, c := range cases {
		var cm []string
		var sec []string

		objs := []runtime.Object{}
		if c.configmap != nil {
			cm = []string{c.configmapName}
			objs = append(objs, mocks.ConfigMap("test", "test", c.configmap))
		}
		if c.secrets != nil {
			sec = []string{c.secretName}
			objs = append(objs, mocks.Secret("test", "test", c.secrets))
		}

		containers := []v1.Container{mocks.Container(c.env, cm, sec)}
		r := FromContainers(containers)

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
	envCase := map[string]string{"env1": "val", "env2": "val2"}
	configmapCase := map[string]string{"config": "val", "config2": "val2"}
	secretsCase := map[string][]byte{"secret": []byte("val"), "secret2": []byte("val2")}
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
		{env: envCase, configmap: configmapCase, secrets: secretsCase, configmapName: "test", secretName: "test"},
		{configmap: configmapCase, secrets: secretsCase, configmapName: "test", secretName: "test"},
		{env: envCase, secrets: secretsCase, secretName: "test"},
		{env: envCase, configmap: configmapCase, configmapName: "test"},
		{configmap: configmapCase, configmapName: "test1", shouldErr: true},
		{secrets: secretsCase, secretName: "test1", shouldErr: true},
		{env: envCase, configmap: configmapCase, useFileWriter: true, filename: "test.out", configmapName: "test"},
		{env: envCase, configmap: configmapCase, useFileWriter: true, configmapName: "test", shouldErr: true},
	}

	for i, c := range cases {
		var cm []string
		var sec []string

		objs := []runtime.Object{}
		if c.configmap != nil {
			cm = []string{c.configmapName}
			objs = append(objs, mocks.ConfigMap("test", "test", c.configmap))
		}
		if c.secrets != nil {
			sec = []string{c.secretName}
			objs = append(objs, mocks.Secret("test", "test", c.secrets))
		}

		containers := []v1.Container{mocks.Container(c.env, cm, sec)}
		r := FromContainers(containers)

		var b bytes.Buffer
		var err error
		var got string
		opt := options.NewOptions()
		opt.Client = fake.NewSimpleClientset(objs...)
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
