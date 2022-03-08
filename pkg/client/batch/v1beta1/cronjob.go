package v1beta1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/pkg/result"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CronJob returns a single resource.
func (batchv1beta1 *BatchV1Beta1) CronJob(resource string) *result.Result {
	resp, err := batchv1beta1.
		BatchV1beta1Interface.
		CronJobs(batchv1beta1.options.Namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		return result.NewFromError(NewResourceLoadError("CronJob", err))
	}

	return result.NewFromContainers(
		batchv1beta1.kubeClient,
		batchv1beta1.options.Namespace,
		batchv1beta1.options.ShouldExport,
		resp.Spec.JobTemplate.Spec.Template.Spec.Containers,
	)
}

// CronJobList returns a list of cronjobs.
func (batchv1beta1 *BatchV1Beta1) CronJobList() ([]string, error) {
	res := []string{}

	resp, err := batchv1beta1.
		BatchV1beta1Interface.
		CronJobs(batchv1beta1.options.Namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, NewResourceLoadError("CronJobs", err)
	}

	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
