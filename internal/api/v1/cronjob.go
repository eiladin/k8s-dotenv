package v1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/internal/environment"
	"github.com/eiladin/k8s-dotenv/internal/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CronJob(opt *options.Options) (*environment.Result, error) {
	res := environment.NewResult()
	resp, err := opt.Client.BatchV1().CronJobs(opt.Namespace).Get(context.TODO(), opt.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	containers := resp.Spec.JobTemplate.Spec.Template.Spec.Containers

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

	return res, nil
}

func CronJobs(opt *options.Options) ([]string, error) {
	res := []string{}

	resp, err := opt.Client.BatchV1().CronJobs(opt.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}