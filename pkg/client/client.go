package client

import (
	"fmt"
	"io"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ConfigureFunc is used for configuring `Client` settings.
type ConfigureFunc = func(client *Client)

// Client is used to interact with the kubernetes API.
type Client struct {
	kubernetes.Interface
	shouldExport bool
	namespace    string
	filename     string
	writer       io.Writer
	result       *Result
	Error        error
	appsv1       *AppsV1
	batchv1      *BatchV1
	batchv1beta1 *BatchV1Beta1
	corev1       *CoreV1
}

// NewClient creates `Client` from a kubernetes client.
func NewClient(kubeClient kubernetes.Interface, configures ...ConfigureFunc) *Client {
	client := Client{}
	client.Interface = kubeClient
	client.appsv1 = NewAppsV1(&client)
	client.batchv1 = NewBatchV1(&client)
	client.batchv1beta1 = NewBatchV1Beta1(&client)
	client.corev1 = NewCoreV1(&client)

	for _, configure := range configures {
		configure(&client)
	}

	return &client
}

// AppsV1 is used to interact with features provided by the apps group.
func (client *Client) AppsV1() *AppsV1 {
	return client.appsv1
}

// BatchV1 is used to interact with features provided by the batch group.
func (client *Client) BatchV1() *BatchV1 {
	return client.batchv1
}

// BatchV1Beta1 is used to interact with features provided by the batch group.
func (client *Client) BatchV1Beta1() *BatchV1Beta1 {
	return client.batchv1beta1
}

// CoreV1 is used to interact with features provided by the core group.
func (client *Client) CoreV1() *CoreV1 {
	return client.corev1
}

func (client *Client) setDefaultWriter() error {
	if client.writer != nil {
		return nil
	}

	if client.filename == "" {
		return ErrNoFilename
	}

	//nolint
	f, err := os.OpenFile(client.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}

	client.writer = f

	return nil
}

// CurrentNamespace returns the namespace from `~/.kube/config`.
func CurrentNamespace(namespace string, configPath string) (string, error) {
	if namespace != "" {
		return namespace, nil
	}

	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	if configPath != "" {
		rules.ExplicitPath = configPath
	}

	clientCfg, err := rules.Load()
	if err != nil {
		return "", ErrNamespaceResolution
	}

	ns := clientCfg.Contexts[clientCfg.CurrentContext].Namespace
	if ns == "" {
		return "default", nil
	}

	return ns, nil
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
