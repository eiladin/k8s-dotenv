package client

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type AppsV1 struct {
	v1.AppsV1Interface
	client *Client
}

func NewAppsV1(client *Client) *AppsV1 {
	return &AppsV1{
		client:          client,
		AppsV1Interface: client.Interface.AppsV1(),
	}
}

// DaemonSet returns a single resource with the given name.
func (appsv1 *AppsV1) DaemonSet(resource string) *Client {
	resp, err := appsv1.
		AppsV1Interface.
		DaemonSets(appsv1.client.namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		appsv1.client.Error = err

		return appsv1.client
	}

	appsv1.client.result = resultFromContainers(resp.Spec.Template.Spec.Containers)

	return appsv1.client
}

// DaemonSets returns a list of resources.
func (appsv1 *AppsV1) DaemonSets() ([]string, error) {
	resp, err := appsv1.
		AppsV1Interface.
		DaemonSets(appsv1.client.namespace).
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, NewResourceLoadError("DaemonSets", err)
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}

// Deployment returns a single resource with the given name.
func (appsv1 *AppsV1) Deployment(resource string) *Client {
	resp, err := appsv1.
		AppsV1Interface.
		Deployments(appsv1.client.namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		appsv1.client.Error = err

		return appsv1.client
	}

	appsv1.client.result = resultFromContainers(resp.Spec.Template.Spec.Containers)

	return appsv1.client
}

// Deployments returns a list of resources.
func (appsv1 *AppsV1) Deployments() ([]string, error) {
	resp, err := appsv1.
		AppsV1Interface.
		Deployments(appsv1.client.namespace).
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, NewResourceLoadError("Deployments", err)
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}

// ReplicaSet returns a single resource with the given name.
func (appsv1 *AppsV1) ReplicaSet(resource string) *Client {
	resp, err := appsv1.
		AppsV1Interface.
		ReplicaSets(appsv1.client.namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		appsv1.client.Error = err

		return appsv1.client
	}

	appsv1.client.result = resultFromContainers(resp.Spec.Template.Spec.Containers)

	return appsv1.client
}

// ReplicaSets returns a list of resources.
func (appsv1 *AppsV1) ReplicaSets() ([]string, error) {
	resp, err := appsv1.
		AppsV1Interface.
		ReplicaSets(appsv1.client.namespace).
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, NewResourceLoadError("ReplicaSets", err)
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
