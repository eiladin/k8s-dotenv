package clioptions

import (
	"fmt"
	"io"

	"github.com/eiladin/k8s-dotenv/pkg/kubeclient"
	"k8s.io/client-go/kubernetes"
)

// CLIOptions contains configuration used to interact with the kubernetes API.
type CLIOptions struct {
	KubeClient   kubernetes.Interface
	Namespace    string
	ResourceName string
	Filename     string
	NoExport     bool
	Writer       io.Writer
}

// ResolveNamespace sets the Namespace property of an Options struct.
func (opt *CLIOptions) ResolveNamespace() error {
	ns, err := kubeclient.CurrentNamespace()
	if err != nil {
		return fmt.Errorf("resolve namespace: %w", err)
	}

	opt.Namespace = ns

	return nil
}