package options

import (
	"errors"
	"io"
	"os"

	"github.com/eiladin/k8s-dotenv/internal/client"
	"k8s.io/client-go/kubernetes"
)

type Options struct {
	Client    kubernetes.Interface
	Namespace string
	Name      string
	Filename  string
	NoExport  bool
	Writer    io.Writer
}

func NewOptions() *Options {
	return &Options{}
}

func (opt *Options) ResolveNamespace(configPath string) error {
	ns, err := client.CurrentNamespace(opt.Namespace, configPath)
	if err != nil {
		return err
	}
	opt.Namespace = ns
	return nil
}

func (opt *Options) SetWriter(writer io.Writer) error {
	if writer == nil {
		if opt.Filename == "" {
			return errors.New("no filename provided")
		}
		f, err := os.OpenFile(opt.Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		writer = f
	}
	opt.Writer = writer
	return nil
}
