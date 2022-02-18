package environment

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/eiladin/k8s-dotenv/internal/options"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

type EnvironmentSuite struct {
	suite.Suite
}

func mockConfigMap(name string, namespace string, data map[string]string) *v1.ConfigMap {
	res := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
	return res
}

func mockSecret(name string, namespace string, data map[string][]byte) *v1.Secret {
	res := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
	return res
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
			objs = append(objs, mockConfigMap("test", "test", c.configmap))
		}
		if c.secrets != nil {
			objs = append(objs, mockSecret("test", "test", c.secrets))
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
			objs = append(objs, mockConfigMap("test", "test", c.configmap))
		}
		if c.secrets != nil {
			objs = append(objs, mockSecret("test", "test", c.secrets))
		}

		opt := options.NewOptions()
		opt.Client = fake.NewSimpleClientset(objs...)
		opt.Namespace = "test"
		var b bytes.Buffer
		opt.SetWriter(&b)
		err := r.Write(opt)
		got := b.String()
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
