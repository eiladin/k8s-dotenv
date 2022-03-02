package v1beta1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/environment"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CronJob returns a single resource in a given namespace with the given name.
func CronJob(client *client.Client, namespace string, resource string) (*environment.Result, error) {
	resp, err := client.BatchV1beta1().CronJobs(namespace).Get(context.TODO(), resource, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return environment.FromContainers(resp.Spec.JobTemplate.Spec.Template.Spec.Containers), nil
}

// CronJobs returns a list of resources in a given namespace.
func CronJobs(client *client.Client, namespace string) ([]string, error) {
	res := []string{}

	resp, err := client.BatchV1beta1().CronJobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
