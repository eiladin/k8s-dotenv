package v1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceList returns all namespaces in a cluster.
func (corev1 *CoreV1) NamespaceList() ([]string, error) {
	resp, err := corev1.
		CoreV1Interface.
		Namespaces().
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, NewResourceLoadError("Namespaces", err)
	}

	res := []string{}
	for _, ns := range resp.Items {
		res = append(res, ns.Name)
	}

	return res, nil
}
