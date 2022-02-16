package client

import (
	"errors"
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	*kubernetes.Clientset
}

func Get() (*Client, error) {
	var home string
	if home = homedir.HomeDir(); home == "" {
		return nil, errors.New("err")
	}

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{clientset}, nil
}

func CurrentNamespace(namespace string) (string, error) {
	if namespace != "" {
		return namespace, nil
	}

	clientCfg, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return "", err
	}

	ns := clientCfg.Contexts[clientCfg.CurrentContext].Namespace
	if ns == "" {
		return "default", nil
	}
	return ns, nil
}

func (c *Client) GetApiGroup(resource string) (string, error) {
	serverResources, err := c.ServerResources()
	if err != nil {
		return "", err
	}

	for _, r := range serverResources {
		for _, ar := range r.APIResources {
			if ar.Kind == resource {
				return r.GroupVersion, nil
			}
		}
	}

	return "", fmt.Errorf("resource %s not found", resource)
}
