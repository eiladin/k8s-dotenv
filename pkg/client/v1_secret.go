package client

import (
	"context"
	"fmt"
	"sort"

	"github.com/eiladin/k8s-dotenv/pkg/parser"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Get returns the export value(s) given a secret name in a specific namespace.
func (client *Client) SecretV1(secret string, shouldExport bool) (string, error) {
	resp, err := client.CoreV1().Secrets(client.namespace).Get(context.TODO(), secret, metav1.GetOptions{})
	if err != nil {
		return "", ErrMissingResource
	}

	res := fmt.Sprintf("##### SECRET - %s #####\n", secret)
	keys := make([]string, 0, len(resp.Data))

	for k := range resp.Data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		res += parser.Parse(shouldExport, k, resp.Data[k])
	}

	return res, nil
}
