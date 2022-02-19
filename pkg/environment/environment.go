package environment

import (
	"github.com/eiladin/k8s-dotenv/pkg/configmap"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/parser"
	"github.com/eiladin/k8s-dotenv/pkg/secret"
)

type Result struct {
	Environment map[string]string
	Secrets     []string
	ConfigMaps  []string
}

func NewResult() *Result {
	return &Result{
		Environment: map[string]string{},
		Secrets:     []string{},
		ConfigMaps:  []string{},
	}
}

func (r *Result) Output(opt *options.Options) (string, error) {
	res := ""
	for k, v := range r.Environment {
		res += parser.ParseStr(!opt.NoExport, k, v)
	}

	for _, s := range r.Secrets {
		secretVal, err := secret.Get(opt, s)
		if err != nil {
			return "", err
		}
		res += secretVal
	}

	for _, c := range r.ConfigMaps {
		configmapVal, err := configmap.Get(opt, c)
		if err != nil {
			return "", err
		}
		res += configmapVal
	}
	return res, nil
}

func (r *Result) Write(opt *options.Options) error {
	output, err := r.Output(opt)
	if err != nil {
		return err
	}

	if opt.Writer == nil {
		err = opt.SetDefaultWriter()
		if err != nil {
			return err
		}
	}

	_, err = opt.Writer.Write([]byte(output))
	return err
}
