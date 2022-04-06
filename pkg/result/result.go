package result

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/eiladin/k8s-dotenv/pkg/parser"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ErrMissingResource is returned when the resource is not found.
var ErrMissingResource = errors.New("resource not found")

// ErrMissingWriter is returned when no writer has been set.
var ErrMissingWriter = errors.New("missing writer")

func newWriteError(err error) error {
	return fmt.Errorf("write error: %w", err)
}

// Result contains the values of environment variables and names of configmaps and secrets related to a resource.
type Result struct {
	Error        error
	shouldExport bool
	Environment  EnvValues
	Secrets      map[string]EnvValues
	ConfigMaps   map[string]EnvValues
}

func newResult() *Result {
	return &Result{
		Environment: map[string]string{},
		Secrets:     map[string]EnvValues{},
		ConfigMaps:  map[string]EnvValues{},
	}
}

type EnvValues map[string]string

func (env EnvValues) sortedKeys() []string {
	keys := make([]string, 0, len(env))

	for k := range env {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func configMapData(client kubernetes.Interface, namespace, resource string) (map[string]string, error) {
	resp, err := client.
		CoreV1().
		ConfigMaps(namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		return nil, ErrMissingResource
	}

	return resp.Data, nil
}

func secretData(client kubernetes.Interface, namespace, resource string) (map[string]string, error) {
	resp, err := client.
		CoreV1().
		Secrets(namespace).
		Get(context.TODO(), resource, metav1.GetOptions{})

	if err != nil {
		return nil, ErrMissingResource
	}

	res := make(map[string]string)

	for k, v := range resp.Data {
		res[k] = string(v)
	}

	return res, nil
}

func NewFromError(err error) *Result {
	res := newResult()
	res.Error = err

	return res
}

func NewFromContainers(
	client kubernetes.Interface,
	namespace string,
	shouldExport bool,
	containers []corev1.Container,
) *Result {
	res := newResult()
	res.shouldExport = shouldExport

	for _, cont := range containers {
		for _, env := range cont.Env {
			res.Environment[env.Name] = env.Value
		}

		for _, envFrom := range cont.EnvFrom {
			if envFrom.ConfigMapRef != nil {
				name := envFrom.ConfigMapRef.Name
				configMap, err := configMapData(client, namespace, name)

				if err != nil {
					return NewFromError(err)
				}

				res.ConfigMaps[name] = configMap
			}

			if envFrom.SecretRef != nil {
				name := envFrom.SecretRef.Name
				sec, err := secretData(client, namespace, name)

				if err != nil {
					return NewFromError(err)
				}

				res.Secrets[name] = sec
			}
		}
	}

	return res
}

func (r *Result) parse() string {
	var res string

	envKeys := r.Environment.sortedKeys()
	for _, k := range envKeys {
		res += parser.ParseStr(r.shouldExport, k, r.Environment[k])
	}

	for k, v := range r.ConfigMaps {
		res += fmt.Sprintf("##### CONFIGMAP - %s #####\n", k)
		for _, key := range v.sortedKeys() {
			res += parser.ParseStr(r.shouldExport, key, v[key])
		}
	}

	for k, v := range r.Secrets {
		res += fmt.Sprintf("##### SECRET - %s #####\n", k)
		for _, key := range v.sortedKeys() {
			res += parser.ParseStr(r.shouldExport, key, v[key])
		}
	}

	return res
}

func (r *Result) Write(writer io.Writer) error {
	if r.Error != nil {
		return r.Error
	}

	if writer == nil {
		return ErrMissingWriter
	}

	output := r.parse()

	if _, err := writer.Write([]byte(output)); err != nil {
		return newWriteError(err)
	}

	return nil
}
