package v1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/environment"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pod returns a single resource in a given namespace with the given name.
func Pod(client *client.Client, namespace string, resource string) (*environment.Result, error) {
	resp, err := client.CoreV1().Pods(namespace).Get(context.TODO(), resource, metav1.GetOptions{})
	if err != nil {
		return nil, NewResourceLoadError("Pod", err)
	}

	return environment.FromContainers(resp.Spec.Containers), nil
}

// Pods returns a list of resources in a given namespace.
func Pods(client *client.Client, namespace string) ([]string, error) {
	resp, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, NewResourceLoadError("Pods", err)
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
