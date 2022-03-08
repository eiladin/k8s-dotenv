package v1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/pkg/result"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReplicaSet returns a single resource with the given name.
func (appsv1 *AppsV1) ReplicaSet(resource string) *result.Result {
	resp, err := appsv1.
		AppsV1Interface.
		ReplicaSets(appsv1.options.Namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		return result.NewFromError(NewResourceLoadError("ReplicaSet", err))
	}

	return result.NewFromContainers(
		appsv1.kubeClient,
		appsv1.options.Namespace,
		appsv1.options.ShouldExport,
		resp.Spec.Template.Spec.Containers,
	)
}

// ReplicaSetList returns a list of replicasets.
func (appsv1 *AppsV1) ReplicaSetList() ([]string, error) {
	resp, err := appsv1.
		AppsV1Interface.
		ReplicaSets(appsv1.options.Namespace).
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
