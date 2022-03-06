package client

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CronJob returns a single resource in a given namespace with the given name.
func (client *Client) Namespaces() ([]string, error) {
	resp, err := client.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, NewResourceLoadError("Namespaces", err)
	}

	res := []string{}
	for _, ns := range resp.Items {
		res = append(res, ns.Name)
	}

	return res, nil
}
