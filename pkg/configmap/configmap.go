package configmap

import (
	"context"
	"fmt"
	"sort"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/parser"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Get(opt *options.Options, configmap string) (string, error) {
	resp, err := opt.Client.CoreV1().ConfigMaps(opt.Namespace).Get(context.TODO(), configmap, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf("##### CONFIGMAP - %s #####\n", configmap)

	keys := make([]string, 0, len(resp.Data))
	for k := range resp.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		res += parser.ParseStr(!opt.NoExport, k, resp.Data[k])
	}

	return res, nil
}
