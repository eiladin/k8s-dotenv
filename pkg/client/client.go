package client

import (
	"errors"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// ErrResourceNameRequired is returned when no resource name is provided.
var ErrResourceNameRequired = errors.New("resource name required")

// ErrMissingResource is returned when the resource is not found.
var ErrMissingResource = errors.New("resource not found")

// ErrCreatingClient is returned when the client cannot be created.
var ErrCreatingClient = errors.New("client create error")

// ErrNamespaceResolution is returned when the current namespace cannot be resolved.
var ErrNamespaceResolution = errors.New("current resolution error")

// ErrAPIGroup is returned when a kubernetes api call fails.
var ErrAPIGroup = errors.New("api group error")

// Client is a wrapper around kubernetes.Interface.
type Client struct {
	kubernetes.Interface
}

// NewClient returns a Client wrapping the given kubernetes.Interface.
func NewClient(cl kubernetes.Interface) *Client {
	return &Client{cl}
}

// Get returns a configured Client loaded from ~/.kube/config.
func Get() (*Client, error) {
	var home string
	if home = homedir.HomeDir(); home == "" {
		return nil, ErrCreatingClient
	}

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		return nil, ErrCreatingClient
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, ErrCreatingClient
	}

	return &Client{clientset}, nil
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
