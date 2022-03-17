package mock

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pod returns a mock struct.
func Pod(name, namespace string, env map[string]string, configmaps, secrets []string) *corev1.Pod {
	res := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: map[string]string{},
		},
	}

	res.Spec.Containers = []corev1.Container{Container(env, configmaps, secrets)}

	return res
}
