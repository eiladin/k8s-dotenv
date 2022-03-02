package mock

import (
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CronJobv1beta1(name, namespace string, env map[string]string, configmaps, secrets []string) *batchv1beta1.CronJob {
	res := &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: map[string]string{},
		},
	}

	res.Spec.JobTemplate.Spec.Template.Spec.Containers = []corev1.Container{Container(env, configmaps, secrets)}
	return res
}

func CronJobv1(name, namespace string, env map[string]string, configmaps, secrets []string) *batchv1.CronJob {
	res := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: map[string]string{},
		},
	}

	res.Spec.JobTemplate.Spec.Template.Spec.Containers = []corev1.Container{Container(env, configmaps, secrets)}
	return res
}
