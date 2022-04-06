package v1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (corev1 *CoreV1) SecretData(resource string) (map[string]string, error) {
	resp, err := corev1.
		Secrets(corev1.options.Namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		return nil, ErrMissingResource
	}

	res := make(map[string]string)

	for k, v := range resp.Data {
		res[k] = string(v)
	}

	return res, nil
}
