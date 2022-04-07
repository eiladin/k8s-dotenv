package v1

import (
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// CoreV1 is used to interact with features provided by the core group.
type CoreV1 struct {
	v1.CoreV1Interface
	kubeClient kubernetes.Interface
	options    *options.Client
}

// NewCoreV1 creates `CoreV1`.
func NewCoreV1(kubeClient kubernetes.Interface, options *options.Client) *CoreV1 {
	return &CoreV1{
		options:         options,
		kubeClient:      kubeClient,
		CoreV1Interface: kubeClient.CoreV1(),
	}
}
