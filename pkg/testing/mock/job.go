package mock

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Job returns a mock struct.
func Job(name, namespace string, env map[string]string, configmaps, secrets []string) *batchv1.Job {
	res := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: map[string]string{},
		},
	}

	res.Spec.Template.Spec.Containers = []corev1.Container{Container(env, configmaps, secrets)}

	return res
}
