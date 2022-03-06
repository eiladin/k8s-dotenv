package corev1

import v1 "k8s.io/client-go/kubernetes/typed/core/v1"

type CoreV1 struct {
	namespace string
	v1.CoreV1Interface
}

func NewCoreV1(coreV1Interface v1.CoreV1Interface, namespace string) *CoreV1 {
	return &CoreV1{
		CoreV1Interface: coreV1Interface,
		namespace:       namespace,
	}
}

func (v1 *CoreV1) SetNamespace(namespace string) {
	v1.namespace = namespace
}
