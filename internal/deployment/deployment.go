package deployment

import (
	"context"

	"github.com/eiladin/k8s-dotenv/internal/client"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Get(namespace string, name string) ([]string, []string, error) {
	secrets := []string{}
	configmaps := []string{}

	clientset, err := client.Get()
	if err != nil {
		return nil, nil, err
	}

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	for _, cont := range deployment.Spec.Template.Spec.Containers {
		for _, envFrom := range cont.EnvFrom {
			if envFrom.SecretRef != nil {
				secrets = append(secrets, envFrom.SecretRef.Name)
			}
			if envFrom.ConfigMapRef != nil {
				configmaps = append(configmaps, envFrom.ConfigMapRef.Name)
			}
		}
	}

	return secrets, configmaps, nil
}
