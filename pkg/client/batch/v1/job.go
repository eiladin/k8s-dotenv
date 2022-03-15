package v1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/pkg/result"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Job returns a single resource with the given name.
func (batchv1 *BatchV1) Job(resource string) *result.Result {
	resp, err := batchv1.
		BatchV1Interface.
		Jobs(batchv1.options.Namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		return result.NewFromError(NewResourceLoadError("Job", err))
	}

	return result.NewFromContainers(
		batchv1.kubeClient,
		batchv1.options.Namespace,
		batchv1.options.ShouldExport,
		resp.Spec.Template.Spec.Containers,
	)
}

// JobList returns a list of jobs.
func (batchv1 *BatchV1) JobList() ([]string, error) {
	resp, err := batchv1.
		BatchV1Interface.
		Jobs(batchv1.options.Namespace).
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
