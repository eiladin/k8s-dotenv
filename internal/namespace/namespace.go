package namespace

import (
	"context"

	"github.com/eiladin/k8s-dotenv/internal/client"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetList() ([]string, error) {
	clientset, err := client.Get()
	if err != nil {
		return nil, err
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, ns := range namespaces.Items {
		res = append(res, ns.Name)
	}
	return res, nil
}
