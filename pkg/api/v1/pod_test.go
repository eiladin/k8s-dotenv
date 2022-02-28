package v1

import (
	"errors"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mocks"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

type PodSuite struct {
	suite.Suite
}

func (suite PodSuite) TestPod() {
	cases := []struct {
		name       string
		namespace  string
		env        map[string]string
		configmaps []string
		secrets    []string
		shouldErr  bool
	}{
		{name: "my-pod", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-pod", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{}},
		{name: "my-pod", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-pod", namespace: "test", env: map[string]string{"k1": "v1", "k2": "v2"}, configmaps: []string{}, secrets: []string{}},
		{name: "my-pod", namespace: "test", configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-pod", namespace: "test", configmaps: []string{"ConfigMap0", "ConfigMap1"}, secrets: []string{}},
		{name: "my-pod", namespace: "test", configmaps: []string{}, secrets: []string{"Secret0", "Secret1"}},
		{name: "my-pod", namespace: "test", configmaps: []string{}, secrets: []string{}},
		{shouldErr: true},
	}

	for _, c := range cases {
		m := mocks.Pod(c.name, c.namespace, c.env, c.configmaps, c.secrets)
		client := fake.NewSimpleClientset(m)
		if c.shouldErr {
			client.PrependReactor("get", "pods", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &corev1.Pod{}, errors.New("error getting deployment")
			})
		}

		opt := &options.Options{
			Client:    client,
			Namespace: c.namespace,
			Name:      c.name,
		}

		got, err := Pod(opt)
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

func (suite PodSuite) TestPods() {
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
			items:         []item{{name: "my-pod", namespace: "test"}},
			expectedCount: 1,
		},
		{
			namespace:     "test",
			items:         []item{{name: "my-pod", namespace: "test"}, {name: "my-pod-2", namespace: "test"}},
			expectedCount: 2,
		},
		{
			namespace:     "other",
			items:         []item{{name: "my-pod", namespace: "test"}, {name: "my-pod-2", namespace: "test"}},
			expectedCount: 0,
		},
		{
			namespace:     "test",
			items:         []item{{name: "my-pod", namespace: "test"}, {name: "my-pod-2", namespace: "other"}},
			expectedCount: 1,
		},
		{
			shouldErr: true,
		},
	}

	for _, c := range cases {
		ms := []runtime.Object{}
		for _, item := range c.items {
			mock := mocks.Pod(item.name, item.namespace, nil, nil, nil)
			ms = append(ms, mock)
		}
		client := fake.NewSimpleClientset(ms...)
		if c.shouldErr {
			client.PrependReactor("list", "pods", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, &corev1.PodList{}, errors.New("error getting pod list")
			})
		}

		opt := &options.Options{
			Client:    client,
			Namespace: c.namespace,
		}

		got, err := Pods(opt)
		if c.shouldErr {
			suite.Error(err)
		} else {
			suite.NoError(err)
			suite.NotNil(got)
			suite.Len(got, c.expectedCount)
		}
	}
}

func TestPodSuite(t *testing.T) {
	suite.Run(t, new(PodSuite))
}
