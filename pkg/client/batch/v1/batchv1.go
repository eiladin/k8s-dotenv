package v1

import (
	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

// BatchV1 is used to interact with features provided by the batch group.
type BatchV1 struct {
	v1.BatchV1Interface
	kubeClient kubernetes.Interface
	options    *clientoptions.Clientoptions
}

// NewBatchV1 creates `BatchV1`.
func NewBatchV1(client kubernetes.Interface, options *clientoptions.Clientoptions) *BatchV1 {
	return &BatchV1{
		options:          options,
		kubeClient:       client,
		BatchV1Interface: client.BatchV1(),
	}
}
