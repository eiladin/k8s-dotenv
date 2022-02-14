package client

import (
	"errors"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func Get() (*kubernetes.Clientset, error) {
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

	return clientset, nil
}

func CurrentNamespace() (string, error) {
	clientCfg, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return "", err
	}

	namespace := clientCfg.Contexts[clientCfg.CurrentContext].Namespace
	if namespace == "" {
		return "default", nil
	}
	return namespace, nil
}
