package environment

import (
	"sort"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/configmap"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/parser"
	"github.com/eiladin/k8s-dotenv/pkg/secret"
	v1 "k8s.io/api/core/v1"
)

// Result contains the values of environment variables and names of configmaps and secrets related to a resource.
type Result struct {
	Environment map[string]string
	Secrets     []string
	ConfigMaps  []string
}

// NewResult constructor.
func NewResult() *Result {
	return &Result{
		Environment: map[string]string{},
		Secrets:     []string{},
		ConfigMaps:  []string{},
	}
}

// FromContainers creates a result object from a list of containers.
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

func (r *Result) output(client *client.Client, namespace string, shouldExport bool) (string, error) {
	res := ""
	keys := make([]string, 0, len(r.Environment))
	for k := range r.Environment {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		res += parser.ParseStr(shouldExport, k, r.Environment[k])
	}

	sort.Strings(r.Secrets)
	for _, s := range r.Secrets {
		secretVal, err := secret.Get(client, namespace, s, shouldExport)
		if err != nil {
			return "", err
		}
		res += secretVal
	}

	sort.Strings(r.ConfigMaps)
	for _, c := range r.ConfigMaps {
		configmapVal, err := configmap.Get(client, namespace, c, shouldExport)
		if err != nil {
			return "", err
		}
		res += configmapVal
	}
	return res, nil
}

func (r *Result) Write(opt *options.Options) error {
	output, err := r.output(opt.Client, opt.Namespace, !opt.NoExport)
	if err != nil {
		return err
	}

	err = opt.SetDefaultWriter()
	if err != nil {
		return err
	}

	_, err = opt.Writer.Write([]byte(output))
	return err
}
