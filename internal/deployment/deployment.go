package deployment

import (
	"context"
	"fmt"
	"strings"

	"github.com/eiladin/k8s-dotenv/internal/client"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Get(namespace string, name string) ([]string, []string, []string, error) {
	secrets := []string{}
	configmaps := []string{}
	environment := []string{}

	clientset, err := client.Get()
	if err != nil {
		return nil, nil, nil, err
	}

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		return nil, nil, nil, err
	}

	for _, cont := range deployment.Spec.Template.Spec.Containers {
		for _, env := range cont.Env {
			environment = append(environment, fmt.Sprintf("export %s=\"%s\"", strings.ReplaceAll(env.Name, ".", ""), strings.ReplaceAll(string(env.Value), "\n", "\\n")))
		}

		for _, envFrom := range cont.EnvFrom {
			if envFrom.SecretRef != nil {
				secrets = append(secrets, envFrom.SecretRef.Name)
			}
			if envFrom.ConfigMapRef != nil {
				configmaps = append(configmaps, envFrom.ConfigMapRef.Name)
			}
		}
	}

	return environment, secrets, configmaps, nil
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
