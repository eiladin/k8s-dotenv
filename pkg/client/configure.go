package client

import (
	"io"

	"k8s.io/client-go/kubernetes"
)

// WithKubeClient sets the underlying kubernetes API client.
func WithKubeClient(kubeClient kubernetes.Interface) ConfigureFunc {
	return func(client *Client) {
		client.Interface = kubeClient
		client.appsv1 = NewAppsV1(client)
		client.batchv1 = NewBatchV1(client)
		client.batchv1beta1 = NewBatchV1Beta1(client)
		client.corev1 = NewCoreV1(client)
	}
}

// WithExport flags the client to include `export` statements in the output.
func WithExport(shouldExport bool) ConfigureFunc {
	return func(client *Client) {
		client.shouldExport = shouldExport
	}
}

// WithWriter sets the `io.Writer` to use for output.
func WithWriter(writer io.Writer) ConfigureFunc {
	return func(client *Client) {
		client.writer = writer
	}
}

// WithNamespace sets the namespace to use when interacting with the Kubernetes API.
func WithNamespace(namespace string) ConfigureFunc {
	return func(client *Client) {
		client.namespace = namespace
	}
}

// WithFilename sets the name of the file to output into.
func WithFilename(filename string) ConfigureFunc {
	return func(client *Client) {
		client.filename = filename
	}
}
