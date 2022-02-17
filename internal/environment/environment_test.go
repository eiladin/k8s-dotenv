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
		name      string
		namespace string
		env       map[string]string
		configmap map[string]string
		secrets   map[string][]byte
	}{
		{
			name:      "test",
			namespace: "test",
			env:       map[string]string{"env1": "val", "env2": "val2"},
			configmap: map[string]string{"config": "val", "config2": "val2"},
			secrets:   map[string][]byte{"secret": []byte("val"), "secret2": []byte("val2")},
		},
		{
			name:      "test",
			namespace: "test",
			configmap: map[string]string{"config": "val", "config2": "val2"},
			secrets:   map[string][]byte{"secret": []byte("val"), "secret2": []byte("val2")},
		},
		{
			name:      "test",
			namespace: "test",
			env:       map[string]string{"env1": "val", "env2": "val2"},
			secrets:   map[string][]byte{"secret": []byte("val"), "secret2": []byte("val2")},
		},
		{
			name:      "test",
			namespace: "test",
			env:       map[string]string{"env1": "val", "env2": "val2"},
			configmap: map[string]string{"config": "val", "config2": "val2"},
		},
	}

	for i, c := range cases {
		r := NewResult()
		if c.configmap != nil {
			r.ConfigMaps = []string{c.name}
		}
		if c.secrets != nil {
			r.Secrets = []string{c.name}
		}
		if c.env != nil {
			r.Environment = c.env
		}

		objs := []runtime.Object{}
		if c.configmap != nil {
			objs = append(objs, mockConfigMap(c.name, c.namespace, c.configmap))
		}
		if c.secrets != nil {
			objs = append(objs, mockSecret(c.name, c.namespace, c.secrets))
		}

		opt := options.NewOptions()
		opt.Client = fake.NewSimpleClientset(objs...)
		opt.Namespace = c.namespace

		got, err := r.Output(opt)
		suite.NoError(err, fmt.Sprintf("Test case %d", i))
		suite.NotNil(got, fmt.Sprintf("Test case %d", i))
		for k, v := range c.configmap {
			suite.Contains(got, fmt.Sprintf("export %s=\"%s\"", k, v), fmt.Sprintf("Test case %d", i))
		}
		for k, v := range c.secrets {
			suite.Contains(got, fmt.Sprintf("export %s=\"%s\"", k, v), fmt.Sprintf("Test case %d", i))
		}
		for k, v := range c.env {
			suite.Contains(got, fmt.Sprintf("export %s=\"%s\"", k, v), fmt.Sprintf("Test case %d", i))
		}
	}
}

func (suite EnvironmentSuite) TestWrite() {
	cases := []struct {
		name      string
		namespace string
		env       map[string]string
		configmap map[string]string
		secrets   map[string][]byte
	}{
		{
			name:      "test",
			namespace: "test",
			env:       map[string]string{"env1": "val", "env2": "val2"},
			configmap: map[string]string{"config": "val", "config2": "val2"},
			secrets:   map[string][]byte{"secret": []byte("val"), "secret2": []byte("val2")},
		},
		{
			name:      "test",
			namespace: "test",
			configmap: map[string]string{"config": "val", "config2": "val2"},
			secrets:   map[string][]byte{"secret": []byte("val"), "secret2": []byte("val2")},
		},
		{
			name:      "test",
			namespace: "test",
			env:       map[string]string{"env1": "val", "env2": "val2"},
			secrets:   map[string][]byte{"secret": []byte("val"), "secret2": []byte("val2")},
		},
		{
			name:      "test",
			namespace: "test",
			env:       map[string]string{"env1": "val", "env2": "val2"},
			configmap: map[string]string{"config": "val", "config2": "val2"},
		},
	}

	for i, c := range cases {
		r := NewResult()
		if c.configmap != nil {
			r.ConfigMaps = []string{c.name}
		}
		if c.secrets != nil {
			r.Secrets = []string{c.name}
		}
		if c.env != nil {
			r.Environment = c.env
		}

		objs := []runtime.Object{}
		if c.configmap != nil {
			objs = append(objs, mockConfigMap(c.name, c.namespace, c.configmap))
		}
		if c.secrets != nil {
			objs = append(objs, mockSecret(c.name, c.namespace, c.secrets))
		}

		opt := options.NewOptions()
		opt.Client = fake.NewSimpleClientset(objs...)
		opt.Namespace = c.namespace
		var b bytes.Buffer

		err := r.Write(&b, opt)
		got := b.String()
		suite.NoError(err, fmt.Sprintf("Test case %d", i))
		suite.NotNil(got, fmt.Sprintf("Test case %d", i))
		for k, v := range c.configmap {
			suite.Contains(got, fmt.Sprintf("export %s=\"%s\"", k, v), fmt.Sprintf("Test case %d", i))
		}
		for k, v := range c.secrets {
			suite.Contains(got, fmt.Sprintf("export %s=\"%s\"", k, v), fmt.Sprintf("Test case %d", i))
		}
		for k, v := range c.env {
			suite.Contains(got, fmt.Sprintf("export %s=\"%s\"", k, v), fmt.Sprintf("Test case %d", i))
		}
	}
}

func TestEnvironmentSuite(t *testing.T) {
	suite.Run(t, new(EnvironmentSuite))
}
