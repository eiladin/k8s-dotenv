package client

import (
	appsv1 "github.com/eiladin/k8s-dotenv/pkg/client/apps/v1"
	batchv1 "github.com/eiladin/k8s-dotenv/pkg/client/batch/v1"
	batchv1beta1 "github.com/eiladin/k8s-dotenv/pkg/client/batch/v1beta1"
	corev1 "github.com/eiladin/k8s-dotenv/pkg/client/core/v1"
	"k8s.io/client-go/kubernetes"
)

// WithKubeClient sets the underlying kubernetes API client.
func WithKubeClient(kubeClient kubernetes.Interface) ConfigureFunc {
	return func(client *Client) {
		client.Interface = kubeClient
		client.corev1 = corev1.NewCoreV1(kubeClient, client.options)
		client.appsv1 = appsv1.NewAppsV1(kubeClient, client.options)
		client.batchv1 = batchv1.NewBatchV1(kubeClient, client.options)
		client.batchv1beta1 = batchv1beta1.NewBatchV1Beta1(kubeClient, client.options)
	}
}

// WithExport flags the client to include `export` statements in the output.
func WithExport(shouldExport bool) ConfigureFunc {
	return func(client *Client) {
		client.options.ShouldExport = shouldExport
	}
}

// WithNamespace sets the namespace to use when interacting with the Kubernetes API.
func WithNamespace(namespace string) ConfigureFunc {
	return func(client *Client) {
		client.options.Namespace = namespace
	}
}
