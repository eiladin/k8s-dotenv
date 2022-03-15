package client

import (
	appsv1 "github.com/eiladin/k8s-dotenv/pkg/client/apps/v1"
	batchv1 "github.com/eiladin/k8s-dotenv/pkg/client/batch/v1"
	batchv1beta1 "github.com/eiladin/k8s-dotenv/pkg/client/batch/v1beta1"
	corev1 "github.com/eiladin/k8s-dotenv/pkg/client/core/v1"
	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"k8s.io/client-go/kubernetes"
)

// ConfigureFunc is used for configuring `Client` settings.
type ConfigureFunc = func(client *Client)

// Client is used to interact with the kubernetes API.
type Client struct {
	kubernetes.Interface
	options      *clientoptions.Clientoptions
	Error        error
	appsv1       *appsv1.AppsV1
	batchv1      *batchv1.BatchV1
	batchv1beta1 *batchv1beta1.BatchV1Beta1
	corev1       *corev1.CoreV1
}

// NewClient creates `Client` from a kubernetes client.
func NewClient(configures ...ConfigureFunc) *Client {
	client := Client{options: clientoptions.New()}

	for _, configure := range configures {
		configure(&client)
	}

	return &client
}

// AppsV1 is used to interact with features provided by the apps group.
func (client *Client) AppsV1() *appsv1.AppsV1 {
	if client.Interface == nil {
		panic(newMissingKubeClientError("AppsV1"))
	}

	return client.appsv1
}

// BatchV1 is used to interact with features provided by the batch group.
func (client *Client) BatchV1() *batchv1.BatchV1 {
	if client.Interface == nil {
		panic(newMissingKubeClientError("BatchV1"))
	}

	return client.batchv1
}

// BatchV1Beta1 is used to interact with features provided by the batch group.
func (client *Client) BatchV1Beta1() *batchv1beta1.BatchV1Beta1 {
	if client.Interface == nil {
		panic(newMissingKubeClientError("BatchV1Beta1"))
	}

	return client.batchv1beta1
}

// CoreV1 is used to interact with features provided by the core group.
func (client *Client) CoreV1() *corev1.CoreV1 {
	if client.Interface == nil {
		panic(newMissingKubeClientError("CoreV1"))
	}

	return client.corev1
}

// GetAPIGroup returns the GroupVersion (batch/v1, batch/v1beta1, etc) for the given resource.
func (client *Client) GetAPIGroup(resource string) (string, error) {
	serverResources, err := client.Discovery().ServerResources()
	if err != nil {
		return "", ErrAPIGroup
	}

	for _, r := range serverResources {
		for _, ar := range r.APIResources {
			if ar.Kind == resource {
				return r.GroupVersion, nil
			}
		}
	}

	return "", ErrMissingResource
}
