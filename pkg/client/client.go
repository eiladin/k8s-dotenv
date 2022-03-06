package client

import (
	"fmt"
	"io"
	"os"

	corev1 "github.com/eiladin/k8s-dotenv/pkg/client/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ConfigureFunc = func(client *Client)

type Client struct {
	kubernetes.Interface
	shouldExport bool
	namespace    string
	filename     string
	writer       io.Writer
	result       *Result
	Error        error
	corev1       *corev1.CoreV1
}

func NewClient(kubeClient kubernetes.Interface, configures ...ConfigureFunc) *Client {
	client := Client{}
	client.Interface = kubeClient
	client.corev1 = corev1.NewCoreV1(client.Interface.CoreV1(), "")

	for _, configure := range configures {
		configure(&client)
	}

	return &client
}

func (client *Client) CoreV1() *corev1.CoreV1 {
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
