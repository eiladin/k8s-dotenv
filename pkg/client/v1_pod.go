package client

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pod returns a single resource in a given namespace with the given name.
func (client *Client) PodV1(resource string) *Client {
	resp, err := client.CoreV1().Pods(client.namespace).Get(context.TODO(), resource, metav1.GetOptions{})
	if err != nil {
		client.Error = err

		return client
	}

	client.result = resultFromContainers(resp.Spec.Containers)

	return client
}

// PodsV1 returns a list of resources in a given namespace.
func (client *Client) PodsV1() ([]string, error) {
	resp, err := client.CoreV1().Pods(client.namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, NewResourceLoadError("Pods", err)
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
