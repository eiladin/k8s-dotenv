package environment

import (
	"io/ioutil"

	"github.com/eiladin/k8s-dotenv/internal/configmap"
	"github.com/eiladin/k8s-dotenv/internal/parser"
	"github.com/eiladin/k8s-dotenv/internal/secret"
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

func (r *Result) GetOutput(namespace string, shouldExport bool) (string, error) {
	res := ""
	for k, v := range r.Environment {
		res += parser.ParseStr(shouldExport, k, v)
	}

	for _, s := range r.Secrets {
		secretVal, err := secret.Get(namespace, s, shouldExport)
		if err != nil {
			return "", err
		}
		res += secretVal
	}

	for _, c := range r.ConfigMaps {
		configmapVal, err := configmap.Get(namespace, c, shouldExport)
		if err != nil {
			return "", err
		}
		res += configmapVal
	}
	return res, nil
}

func (r *Result) Write(namespace string, shouldExport bool, filename string) error {
	output, err := r.GetOutput(namespace, shouldExport)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		return err
	}
	return nil
}
