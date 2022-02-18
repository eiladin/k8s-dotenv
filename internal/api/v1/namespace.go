package v1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/internal/options"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Namespaces(opt *options.Options) ([]string, error) {
	namespaces, err := opt.Client.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, ns := range namespaces.Items {
		res = append(res, ns.Name)
	}
	return res, nil
}
