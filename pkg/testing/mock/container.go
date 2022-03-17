package mock

import (
	corev1 "k8s.io/api/core/v1"
)

// Container returns a mock struct.
func Container(env map[string]string, configmaps, secrets []string) corev1.Container {
	container := corev1.Container{}

	for k, v := range env {
		container.Env = append(container.Env, corev1.EnvVar{Name: k, Value: v})
	}

	for _, cm := range configmaps {
		container.EnvFrom = append(container.EnvFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cm,
				},
			},
		})
	}

	for _, s := range secrets {
		container.EnvFrom = append(container.EnvFrom, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: s,
				},
			},
		})
	}

	return container
}
