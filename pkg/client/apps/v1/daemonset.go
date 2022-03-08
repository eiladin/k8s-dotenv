package v1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/pkg/result"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DaemonSet returns a single resource with the given name.
func (appsv1 *AppsV1) DaemonSet(resource string) *result.Result {
	resp, err := appsv1.
		AppsV1Interface.
		DaemonSets(appsv1.options.Namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		return result.NewFromError(NewResourceLoadError("DaemonSet", err))
	}

	return result.NewFromContainers(

		appsv1.kubeClient,
		appsv1.options.Namespace,
		appsv1.options.ShouldExport,
		resp.Spec.Template.Spec.Containers,
	)
}

// DaemonSetList returns a list of daemonsets.
func (appsv1 *AppsV1) DaemonSetList() ([]string, error) {
	resp, err := appsv1.
		AppsV1Interface.
		DaemonSets(appsv1.options.Namespace).
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
