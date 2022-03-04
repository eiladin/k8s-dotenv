package mock

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/testing"
)

// FakeClient used in tests.
type FakeClient struct {
	*fake.Clientset
}

// NewFakeClient returns a `FakeClient` with a `fake.ClientSet` internally.
func NewFakeClient(objects ...runtime.Object) *FakeClient {
	return &FakeClient{fake.NewSimpleClientset(objects...)}
}

// WithResources adds an `APIResourceList` to a `FakeClient`.
func (c *FakeClient) WithResources(resourceList *metav1.APIResourceList) *FakeClient {
	c.Fake.Resources = append(c.Fake.Resources, resourceList)

	return c
}

// PrependReactor adds a reactor to the beginning of the chain.
func (c *FakeClient) PrependReactor(verb string, resource string, handled bool, ret runtime.Object, err error) *FakeClient {
	c.Fake.PrependReactor(verb, resource, func(action testing.Action) (bool, runtime.Object, error) {
		return handled, ret, err
	})

	return c
}
