package deployment

import (
	"context"

	"github.com/eiladin/k8s-dotenv/internal/client"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GetResult struct {
	Environment map[string]string
	Secrets     []string
	ConfigMaps  []string
}

func NewGetResult() *GetResult {
	return &GetResult{
		Environment: map[string]string{},
		Secrets:     []string{},
		ConfigMaps:  []string{},
	}
}

func Get(namespace string, name string) (*GetResult, error) {
	res := NewGetResult()

	clientset, err := client.Get()
	if err != nil {
		return nil, err
	}

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	for _, cont := range deployment.Spec.Template.Spec.Containers {
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

	return res, nil
}

func GetList(namespace string) ([]string, error) {
	clientset, err := client.Get()
	if err != nil {
		return nil, err
	}

	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, deploy := range deployments.Items {
		res = append(res, deploy.Name)
	}
	return res, nil
}
