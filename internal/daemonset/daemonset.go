package daemonset

import (
	"context"

	"github.com/eiladin/k8s-dotenv/internal/client"
	"github.com/eiladin/k8s-dotenv/internal/environment"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Get(namespace string, name string) (*environment.Result, error) {
	res := environment.NewResult()

	clientset, err := client.Get()
	if err != nil {
		return nil, err
	}

	resp, err := clientset.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	for _, cont := range resp.Spec.Template.Spec.Containers {
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

	resp, err := clientset.AppsV1().DaemonSets(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}
	return res, nil
}
