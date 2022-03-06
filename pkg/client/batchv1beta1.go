package client

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1beta1"
)

type BatchV1Beta1 struct {
	v1.BatchV1beta1Interface
	client *Client
}

func NewBatchV1Beta1(client *Client) *BatchV1Beta1 {
	return &BatchV1Beta1{
		client:                client,
		BatchV1beta1Interface: client.Interface.BatchV1beta1(),
	}
}

// CronJob returns a single resource.
func (batchv1beta1 *BatchV1Beta1) CronJob(resource string) *Client {
	resp, err := batchv1beta1.
		BatchV1beta1Interface.
		CronJobs(batchv1beta1.client.namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		batchv1beta1.client.Error = NewResourceLoadError("CronJob", err)

		return batchv1beta1.client
	}

	batchv1beta1.client.result = resultFromContainers(resp.Spec.JobTemplate.Spec.Template.Spec.Containers)

	return batchv1beta1.client
}

// CronJobs returns a list of resources.
func (batchv1beta1 *BatchV1Beta1) CronJobs() ([]string, error) {
	res := []string{}

	resp, err := batchv1beta1.
		BatchV1beta1Interface.
		CronJobs(batchv1beta1.client.namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, NewResourceLoadError("CronJobs", err)
	}

	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
