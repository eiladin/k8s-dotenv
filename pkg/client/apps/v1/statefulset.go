package v1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/pkg/result"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StatefulSet returns a single resource with the given name.
func (appsv1 *AppsV1) StatefulSet(resource string) *result.Result {
	resp, err := appsv1.
		AppsV1Interface.
		StatefulSets(appsv1.options.Namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		return result.NewFromError(NewResourceLoadError("StatefulSet", err))
	}

	return result.NewFromContainers(

		appsv1.kubeClient,
		appsv1.options.Namespace,
		appsv1.options.ShouldExport,
		resp.Spec.Template.Spec.Containers,
	)
}

// StatefulSetList returns a list of daemonsets.
func (appsv1 *AppsV1) StatefulSetList() ([]string, error) {
	resp, err := appsv1.
		AppsV1Interface.
		StatefulSets(appsv1.options.Namespace).
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, NewResourceLoadError("StatefulSets", err)
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
