package kubeclient

import (
	"errors"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// ErrMissingHomeDir is returned when the home directory for the current user cannot be found.
var ErrMissingHomeDir = errors.New("unable to find home directory")

// ErrReadingKubeConfig is returned when ~/.kube/config cannot be read.
var ErrReadingKubeConfig = errors.New("unable to read ~/.kube/config")

// ErrCreatingKubeClient is returned when ~/.kube/config cannot be parsed.
var ErrCreatingKubeClient = errors.New("unable to parse ~/.kube/config")

// ErrNamespaceResolution is returned when the current namespace cannot be resolved.
var ErrNamespaceResolution = errors.New("current namespace could not be resolved")

// GetDefault returns a kubernetes clientset by reading the current users ~/.kube/config.
func GetDefault() (kubernetes.Interface, error) {
	var home string
	if home = homedir.HomeDir(); home == "" {
		return nil, ErrMissingHomeDir
	}

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		return nil, ErrReadingKubeConfig
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, ErrCreatingKubeClient
	}

	return clientset, nil
}

// CurrentNamespace returns the namespace from `~/.kube/config`.
func CurrentNamespace() (string, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()

	clientCfg, err := rules.Load()
	if err != nil {
		return "", ErrNamespaceResolution
	}

	if clientCfg.Contexts[clientCfg.CurrentContext].Namespace == "" {
		return "default", nil
	}

	return clientCfg.Contexts[clientCfg.CurrentContext].Namespace, nil
}
