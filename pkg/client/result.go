package client

import (
	"fmt"
	"sort"

	"github.com/eiladin/k8s-dotenv/pkg/parser"
	v1 "k8s.io/api/core/v1"
)

// Result contains the values of environment variables and names of configmaps and secrets related to a resource.
type Result struct {
	Environment envValues
	Secrets     map[string]envValues
	ConfigMaps  map[string]envValues
}

// NewResult constructor.
func NewResult() *Result {
	return &Result{
		Environment: envValues{},
		Secrets:     map[string]envValues{},
		ConfigMaps:  map[string]envValues{},
	}
}

type envValues map[string]string

func (env envValues) sortedKeys() []string {
	keys := make([]string, 0, len(env))

	for k := range env {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func (client *Client) resultFromContainers(containers []v1.Container) *Client {
	res := NewResult()

	for _, cont := range containers {
		for _, env := range cont.Env {
			res.Environment[env.Name] = env.Value
		}

		for _, envFrom := range cont.EnvFrom {
			if envFrom.ConfigMapRef != nil {
				n := envFrom.ConfigMapRef.Name
				res.ConfigMaps[n], client.Error = client.CoreV1().ConfigMapData(n, client.shouldExport)

				if client.Error != nil {
					return client
				}
			}

			if envFrom.SecretRef != nil {
				n := envFrom.SecretRef.Name
				res.Secrets[n], client.Error = client.CoreV1().SecretData(n, client.shouldExport)

				if client.Error != nil {
					return client
				}
			}
		}
	}

	client.result = res

	return client
}

func (r *Result) parse(client *Client) string {
	var res string

	envKeys := r.Environment.sortedKeys()
	for _, k := range envKeys {
		res += parser.ParseStr(client.shouldExport, k, r.Environment[k])
	}

	for k, v := range r.ConfigMaps {
		res += fmt.Sprintf("##### CONFIGMAP - %s #####\n", k)
		for _, key := range v.sortedKeys() {
			res += parser.ParseStr(client.shouldExport, key, v[key])
		}
	}

	for k, v := range r.Secrets {
		res += fmt.Sprintf("##### SECRET - %s #####\n", k)
		for _, key := range v.sortedKeys() {
			res += parser.ParseStr(client.shouldExport, key, v[key])
		}
	}

	return res
}

func (client *Client) Write() error {
	if client.Error != nil {
		return client.Error
	}

	output := client.result.parse(client)

	err := client.setDefaultWriter()
	if err != nil {
		return NewWriteError(err)
	}

	if _, err := client.writer.Write([]byte(output)); err != nil {
		return NewWriteError(err)
	}

	return nil
}
