package client

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReplicaSetV1 returns a single resource in a given namespace with the given name.
func (client *Client) ReplicaSetV1(resource string) *Client {
	resp, err := client.AppsV1().ReplicaSets(client.namespace).Get(context.TODO(), resource, metav1.GetOptions{})
	if err != nil {
		client.Error = err

		return client
	}

	client.result = resultFromContainers(resp.Spec.Template.Spec.Containers)

	return client
}

// ReplicaSetsV1 returns a list of resources in a given namespace.
func (client *Client) ReplicaSetsV1() ([]string, error) {
	resp, err := client.AppsV1().ReplicaSets(client.namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, NewResourceLoadError("Deployments", err)
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
