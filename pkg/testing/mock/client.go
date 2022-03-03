package mock

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/testing"
)

type FakeClient struct {
	*fake.Clientset
}

func NewFakeClient(objects ...runtime.Object) *FakeClient {
	return &FakeClient{fake.NewSimpleClientset(objects...)}
}

func (c *FakeClient) WithResources(resourceList *metav1.APIResourceList) *FakeClient {
	c.Fake.Resources = append(c.Fake.Resources, resourceList)

	return c
}

func (c *FakeClient) PrependReactor(verb string, resource string, handled bool, ret runtime.Object, err error) *FakeClient {
	c.Fake.PrependReactor(verb, resource, func(action testing.Action) (bool, runtime.Object, error) {
		return handled, ret, err
	})

	return c
}
