package options

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/eiladin/k8s-dotenv/pkg/client"
)

// ErrNoFilename is returned when no filename is provided.
var ErrNoFilename = errors.New("no filename provided")

// ErrFileCreate is returned when the output file cannot be created.
var ErrFileCreate = errors.New("error creating output file")

// NewClientError wraps client errors.
func NewClientError(err error) error {
	return fmt.Errorf("client error: %w", err)
}

// Options contains configuration used to interact with the kubernetes API.
type Options struct {
	Client       *client.Client
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
		return NewClientError(err)
	}

	opt.Namespace = ns

	return nil
}

// SetDefaultWriter sets the Writer property of an Options struct.
func (opt *Options) SetDefaultWriter() error {
	if opt.Writer != nil {
		return nil
	}

	if opt.Filename == "" {
		return ErrNoFilename
	}

	f, err := os.OpenFile(opt.Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		return ErrFileCreate
	}

	opt.Writer = f

	return nil
}
