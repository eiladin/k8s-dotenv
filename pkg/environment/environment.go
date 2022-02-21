package environment

import (
	"github.com/eiladin/k8s-dotenv/pkg/configmap"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/parser"
	"github.com/eiladin/k8s-dotenv/pkg/secret"
	v1 "k8s.io/api/core/v1"
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

func FromContainers(containers []v1.Container) *Result {
	res := NewResult()
	for _, cont := range containers {
		for _, env := range cont.Env {
			res.Environment[env.Name] = env.Value
		}

		for _, envFrom := range cont.EnvFrom {
			if envFrom.SecretRef != nil {
				res.Secrets = append(res.Secrets, envFrom.SecretRef.Name)
			}
			if envFrom.ConfigMapRef != nil {
				res.ConfigMaps = append(res.ConfigMaps, envFrom.ConfigMapRef.Name)
			}
		}
	}

	return res
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

	if opt.FileWriter == nil {
		err = opt.SetDefaultFileWriter()
		if err != nil {
			return err
		}
	}

	_, err = opt.FileWriter.Write([]byte(output))
	return err
}
