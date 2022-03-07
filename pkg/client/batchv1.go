package client

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

// BatchV1 is used to interact with features provided by the batch group.
type BatchV1 struct {
	v1.BatchV1Interface
	client *Client
}

// NewBatchV1 creates `BatchV1`.
func NewBatchV1(client *Client) *BatchV1 {
	return &BatchV1{
		client:           client,
		BatchV1Interface: client.Interface.BatchV1(),
	}
}

// CronJob returns a single resource.
func (batchv1 *BatchV1) CronJob(resource string) *Client {
	resp, err := batchv1.
		BatchV1Interface.
		CronJobs(batchv1.client.namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		batchv1.client.Error = NewResourceLoadError("CronJob", err)

		return batchv1.client
	}

	batchv1.client.resultFromContainers(resp.Spec.JobTemplate.Spec.Template.Spec.Containers)

	return batchv1.client
}

// CronJobList returns a list of cronjobs.
func (batchv1 *BatchV1) CronJobList() ([]string, error) {
	res := []string{}

	resp, err := batchv1.
		BatchV1Interface.
		CronJobs(batchv1.client.namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, NewResourceLoadError("CronJobs", err)
	}

	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}

// Job returns a single resource with the given name.
func (batchv1 *BatchV1) Job(resource string) *Client {
	resp, err := batchv1.
		BatchV1Interface.
		Jobs(batchv1.client.namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		batchv1.client.Error = err

		return batchv1.client
	}

	batchv1.client.resultFromContainers(resp.Spec.Template.Spec.Containers)

	return batchv1.client
}

// JobList returns a list of jobs.
func (batchv1 *BatchV1) JobList() ([]string, error) {
	resp, err := batchv1.
		BatchV1Interface.
		Jobs(batchv1.client.namespace).
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, NewResourceLoadError("Jobs", err)
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
