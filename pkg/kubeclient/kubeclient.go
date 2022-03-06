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

// Get returns a kubernetes clientset by reading the current users ~/.kube/config.
func Get() (kubernetes.Interface, error) {
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
