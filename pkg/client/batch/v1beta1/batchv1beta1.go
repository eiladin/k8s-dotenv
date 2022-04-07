package v1beta1

import (
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/batch/v1beta1"
)

// BatchV1Beta1 is used to interact with features provided by the batch group.
type BatchV1Beta1 struct {
	v1.BatchV1beta1Interface
	kubeClient kubernetes.Interface
	options    *options.Client
}

// NewBatchV1Beta1 creates `BatchV1Beta1`.
func NewBatchV1Beta1(client kubernetes.Interface, options *options.Client) *BatchV1Beta1 {
	return &BatchV1Beta1{
		options:               options,
		kubeClient:            client,
		BatchV1beta1Interface: client.BatchV1beta1(),
	}
}
