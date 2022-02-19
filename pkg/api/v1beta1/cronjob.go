package v1beta1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/pkg/environment"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CronJob(opt *options.Options) (*environment.Result, error) {
	resp, err := opt.Client.BatchV1beta1().CronJobs(opt.Namespace).Get(context.TODO(), opt.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return environment.FromContainers(resp.Spec.JobTemplate.Spec.Template.Spec.Containers), nil
}

func CronJobs(opt *options.Options) ([]string, error) {
	res := []string{}

	resp, err := opt.Client.BatchV1beta1().CronJobs(opt.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
