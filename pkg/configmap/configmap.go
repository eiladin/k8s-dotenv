package configmap

import (
	"context"
	"fmt"
	"sort"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/parser"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Get returns the export value(s) given a configmap name in a specific namespace
func Get(client *client.Client, namespace string, resource string, shouldExport bool) (string, error) {
	resp, err := client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), resource, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf("##### CONFIGMAP - %s #####\n", resource)

	keys := make([]string, 0, len(resp.Data))
	for k := range resp.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		res += parser.ParseStr(shouldExport, k, resp.Data[k])
	}

	return res, nil
}
