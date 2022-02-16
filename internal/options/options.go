package options

import (
	"github.com/eiladin/k8s-dotenv/internal/client"
)

type Options struct {
	Client    *client.Client
	Namespace string
	Name      string
	Filename  string
	NoExport  bool
}

func NewOptions() *Options {
	return &Options{}
}

func (opt *Options) ResolveNamespace() error {
	ns, err := client.CurrentNamespace(opt.Namespace)
	if err != nil {
		return err
	}
	opt.Namespace = ns
	return nil
}
