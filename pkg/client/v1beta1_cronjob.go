package client

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CronJob returns a single resource in a given namespace with the given name.
func (client *Client) CronJobV1Beta1(resource string) *Client {
	resp, err := client.BatchV1beta1().CronJobs(client.namespace).Get(context.TODO(), resource, metav1.GetOptions{})
	if err != nil {
		client.Error = NewResourceLoadError("CronJob", err)

		return client
	}

	client.result = resultFromContainers(resp.Spec.JobTemplate.Spec.Template.Spec.Containers)

	return client
}

// CronJobs returns a list of resources in a given namespace.
func (client *Client) CronJobsV1beta1() ([]string, error) {
	res := []string{}

	resp, err := client.BatchV1beta1().CronJobs(client.namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, NewResourceLoadError("CronJobs", err)
	}

	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
