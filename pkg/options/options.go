package options

import (
	"errors"
	"fmt"
	"io"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"k8s.io/client-go/kubernetes"
)

// ErrNoFilename is returned when no filename is provided.
var ErrNoFilename = errors.New("no filename provided")

// Options contains configuration used to interact with the kubernetes API.
type Options struct {
	Client       kubernetes.Interface
	Namespace    string
	ResourceName string
	Filename     string
	NoExport     bool
	Writer       io.Writer
}

// ResolveNamespace sets the Namespace property of an Options struct.
func (opt *Options) ResolveNamespace(configPath string) error {
	ns, err := client.CurrentNamespace(opt.Namespace, configPath)
	if err != nil {
		return fmt.Errorf("resolve namespace: %w", err)
	}

	opt.Namespace = ns

	return nil
}
