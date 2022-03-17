package mock

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReplicaSet returns a mock struct.
func ReplicaSet(name, namespace string, env map[string]string, configmaps, secrets []string) *appsv1.ReplicaSet {
	res := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: map[string]string{},
		},
	}

	res.Spec.Template.Spec.Containers = []corev1.Container{Container(env, configmaps, secrets)}

	return res
}
