package client

import (
	"context"
	"fmt"
	"sort"

	"github.com/eiladin/k8s-dotenv/pkg/parser"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// CoreV1 is used to interact with features provided by the core group.
type CoreV1 struct {
	v1.CoreV1Interface
	client *Client
}

// NewCoreV1 creates `CoreV1`.
func NewCoreV1(client *Client) *CoreV1 {
	return &CoreV1{
		client:          client,
		CoreV1Interface: client.Interface.CoreV1(),
	}
}

// ConfigMapValues returns the export value(s) given a configmap name in a specific namespace.
func (corev1 *CoreV1) ConfigMapValues(resource string, shouldExport bool) (string, error) {
	resp, err := corev1.
		CoreV1Interface.
		ConfigMaps(corev1.client.namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		return "", ErrMissingResource
	}

	res := fmt.Sprintf("##### CONFIGMAP - %s #####\n", resource)

	keys := make([]string, 0, len(resp.Data))
	for k := range resp.Data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		res += parser.ParseStr(shouldExport, k, resp.Data[k])
	}

	return res, nil
}

// SecretValues returns the export value(s) given a secret name in a specific namespace.
func (corev1 *CoreV1) SecretValues(secret string, shouldExport bool) (string, error) {
	resp, err := corev1.
		CoreV1Interface.
		Secrets(corev1.client.namespace).
		Get(context.TODO(), secret, metav1.GetOptions{})

	if err != nil {
		return "", ErrMissingResource
	}

	res := fmt.Sprintf("##### SECRET - %s #####\n", secret)
	keys := make([]string, 0, len(resp.Data))

	for k := range resp.Data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		res += parser.Parse(shouldExport, k, resp.Data[k])
	}

	return res, nil
}

// Namespaces returns a single resource in a given namespace with the given name.
func (corev1 *CoreV1) Namespaces() ([]string, error) {
	resp, err := corev1.
		CoreV1Interface.
		Namespaces().
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, NewResourceLoadError("Namespaces", err)
	}

	res := []string{}
	for _, ns := range resp.Items {
		res = append(res, ns.Name)
	}

	return res, nil
}

// Pod returns a single resource in a given namespace with the given name.
func (corev1 *CoreV1) Pod(resource string) *Client {
	resp, err := corev1.
		CoreV1Interface.
		Pods(corev1.client.namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		corev1.client.Error = err

		return corev1.client
	}

	corev1.client.result = resultFromContainers(resp.Spec.Containers)

	return corev1.client
}

// Pods returns a list of resources in a given namespace.
func (corev1 *CoreV1) Pods() ([]string, error) {
	resp, err := corev1.
		CoreV1Interface.
		Pods(corev1.client.namespace).
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, NewResourceLoadError("Pods", err)
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
