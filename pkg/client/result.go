package client

import (
	"sort"

	"github.com/eiladin/k8s-dotenv/pkg/parser"
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

func resultFromContainers(containers []v1.Container) *Result {
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

func (r *Result) output(client *Client) (string, error) {
	res := ""
	keys := make([]string, 0, len(r.Environment))

	for k := range r.Environment {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	sort.Strings(r.Secrets)
	sort.Strings(r.ConfigMaps)

	for _, k := range keys {
		res += parser.ParseStr(client.shouldExport, k, r.Environment[k])
	}

	for _, s := range r.Secrets {
		secretVal, err := client.CoreV1().Secret(s, client.shouldExport)
		if err != nil {
			return "", NewSecretErr(err)
		}

		res += secretVal
	}

	for _, c := range r.ConfigMaps {
		configmapVal, err := client.CoreV1().ConfigMapV1(c, client.shouldExport)
		if err != nil {
			return "", NewConfigMapError(err)
		}

		res += configmapVal
	}

	return res, nil
}

func (client *Client) Write() error {
	if client.Error != nil {
		return client.Error
	}

	output, err := client.result.output(client)
	if err != nil {
		return err
	}

	err = client.setDefaultWriter()
	if err != nil {
		return NewWriteError(err)
	}

	if _, err := client.writer.Write([]byte(output)); err != nil {
		return NewWriteError(err)
	}

	return nil
}
