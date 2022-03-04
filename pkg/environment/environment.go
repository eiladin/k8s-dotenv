package environment

import (
	"fmt"
	"sort"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/configmap"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/parser"
	"github.com/eiladin/k8s-dotenv/pkg/secret"
	v1 "k8s.io/api/core/v1"
)

// NewSecretErr wraps secret errors.
func NewSecretErr(err error) error {
	return fmt.Errorf("secret error: %w", err)
}

// NewConfigMapError wraps config map errors.
func NewConfigMapError(err error) error {
	return fmt.Errorf("configmap error: %w", err)
}

// NewOptionsError wraps options errors.
func NewOptionsError(err error) error {
	return fmt.Errorf("options error: %w", err)
}

// NewWriteError wraps writer errors.
func NewWriteError(err error) error {
	return fmt.Errorf("write error: %w", err)
}

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
	sort.Strings(r.Secrets)
	sort.Strings(r.ConfigMaps)

	for _, k := range keys {
		res += parser.ParseStr(shouldExport, k, r.Environment[k])
	}

	for _, s := range r.Secrets {
		secretVal, err := secret.Get(client, namespace, s, shouldExport)
		if err != nil {
			return "", NewSecretErr(err)
		}

		res += secretVal
	}

	for _, c := range r.ConfigMaps {
		configmapVal, err := configmap.Get(client, namespace, c, shouldExport)
		if err != nil {
			return "", NewConfigMapError(err)
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
		return NewOptionsError(err)
	}

	if _, err := opt.Writer.Write([]byte(output)); err != nil {
		return NewWriteError(err)
	}

	return nil
}
