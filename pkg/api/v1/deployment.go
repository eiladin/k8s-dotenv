package v1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/environment"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Deployment returns a single resource in a given namespace with the given name.
func Deployment(client *client.Client, namespace string, resource string) (*environment.Result, error) {
	resp, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), resource, metav1.GetOptions{})
	if err != nil {
		return nil, NewResourceLoadError("Deployment", err)
	}

	return environment.FromContainers(resp.Spec.Template.Spec.Containers), nil
}

// Deployments returns a list of resources in a given namespace.
func Deployments(client *client.Client, namespace string) ([]string, error) {
	resp, err := client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, NewResourceLoadError("Deployments", err)
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
