package v1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/pkg/environment"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Pod(opt *options.Options) (*environment.Result, error) {
	res := environment.NewResult()

	resp, err := opt.Client.CoreV1().Pods(opt.Namespace).Get(context.TODO(), opt.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	for _, cont := range resp.Spec.Containers {
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

func Pods(opt *options.Options) ([]string, error) {
	resp, err := opt.Client.CoreV1().Pods(opt.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}
	return res, nil
}
