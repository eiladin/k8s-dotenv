package v1

import (
	"context"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Namespaces returns a list of namespaces.
func Namespaces(client *client.Client) ([]string, error) {
	namespaces, err := client.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, NewResourceLoadError("Namespaces", err)
	}

	res := []string{}
	for _, ns := range namespaces.Items {
		res = append(res, ns.Name)
	}

	return res, nil
}
