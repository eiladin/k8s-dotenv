package v1

import (
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

// AppsV1 is used to interact with features provided by the apps group.
type AppsV1 struct {
	v1.AppsV1Interface
	kubeClient kubernetes.Interface
	options    *options.Client
}

// NewAppsV1 creates `AppsV1`.
func NewAppsV1(client kubernetes.Interface, options *options.Client) *AppsV1 {
	return &AppsV1{
		options:         options,
		kubeClient:      client,
		AppsV1Interface: client.AppsV1(),
	}
}
