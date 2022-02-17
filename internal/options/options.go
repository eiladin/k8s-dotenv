package options

import (
	"github.com/eiladin/k8s-dotenv/internal/client"
	"k8s.io/client-go/kubernetes"
)

type Options struct {
	Client    kubernetes.Interface
	Namespace string
	Name      string
	Filename  string
	NoExport  bool
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
