package options

import (
	"fmt"
	"io"

	"github.com/eiladin/k8s-dotenv/pkg/kubeclient"
	"k8s.io/client-go/kubernetes"
)

// CLI stores configuration and arguments passed to the cli.
type CLI struct {
	KubeClient   kubernetes.Interface
	Namespace    string
	ResourceName string
	Filename     string
	NoExport     bool
	Writer       io.Writer
}

// ResolveNamespace sets the Namespace property of an Options struct.
func (cli *CLI) ResolveNamespace() error {
	ns, err := kubeclient.CurrentNamespace()
	if err != nil {
		return fmt.Errorf("resolve namespace: %w", err)
	}

	cli.Namespace = ns

	return nil
}
