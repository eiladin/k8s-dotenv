package corev1

import (
	"context"
	"fmt"
	"sort"

	"github.com/eiladin/k8s-dotenv/pkg/parser"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Get returns the export value(s) given a configmap name in a specific namespace.
func (client *CoreV1) ConfigMapV1(resource string, shouldExport bool) (string, error) {
	resp, err := client.ConfigMaps(client.namespace).Get(context.TODO(), resource, metav1.GetOptions{})
	if err != nil {
		return "", ErrMissingResource
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
