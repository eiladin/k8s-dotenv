package v1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/environment"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Daemonset returns a single resource in a given namespace with the given name.
func DaemonSet(client *client.Client, namespace string, resource string) (*environment.Result, error) {
	resp, err := client.AppsV1().DaemonSets(namespace).Get(context.TODO(), resource, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return environment.FromContainers(resp.Spec.Template.Spec.Containers), nil
}

// DaemonSets returns a list of resources in a given namespace.
func DaemonSets(client *client.Client, namespace string) ([]string, error) {
	resp, err := client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}
	return res, nil
}
