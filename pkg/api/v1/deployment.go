package v1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/pkg/environment"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Deployment(opt *options.Options) (*environment.Result, error) {
	resp, err := opt.Client.AppsV1().Deployments(opt.Namespace).Get(context.TODO(), opt.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return environment.FromContainers(resp.Spec.Template.Spec.Containers), nil
}

func Deployments(opt *options.Options) ([]string, error) {
	resp, err := opt.Client.AppsV1().Deployments(opt.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}
	return res, nil
}
