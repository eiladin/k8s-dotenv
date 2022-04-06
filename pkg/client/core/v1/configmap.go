package v1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMapData returns a map of key/value pairs given a config map.
func (corev1 *CoreV1) ConfigMapData(resource string) (map[string]string, error) {
	resp, err := corev1.
		ConfigMaps(corev1.options.Namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		return nil, ErrMissingResource
	}

	return resp.Data, nil
}
