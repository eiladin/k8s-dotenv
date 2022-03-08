package v1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/pkg/result"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pod returns a single resource in a given namespace with the given name.
func (corev1 *CoreV1) Pod(resource string) *result.Result {
	resp, err := corev1.
		CoreV1Interface.
		Pods(corev1.options.Namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		return result.NewFromError(NewResourceLoadError("Pod", err))
	}

	return result.NewFromContainers(
		corev1.kubeClient,
		corev1.options.Namespace,
		corev1.options.ShouldExport,
		resp.Spec.Containers,
	)
}

// PodList returns a list of pods.
func (corev1 *CoreV1) PodList() ([]string, error) {
	resp, err := corev1.
		CoreV1Interface.
		Pods(corev1.options.Namespace).
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, NewResourceLoadError("Pods", err)
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
