package secret

import (
	"context"
	"fmt"
	"sort"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/parser"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Get(client *client.Client, namespace string, secret string, shouldExport bool) (string, error) {
	resp, err := client.CoreV1().Secrets(namespace).Get(context.TODO(), secret, metav1.GetOptions{})
	if err != nil {
		return "", err
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
