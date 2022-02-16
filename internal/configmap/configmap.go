package configmap

import (
	"context"
	"fmt"

	"github.com/eiladin/k8s-dotenv/internal/options"
	"github.com/eiladin/k8s-dotenv/internal/parser"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Get(opt *options.Options, configmap string) (string, error) {
	resp, err := opt.Client.CoreV1().ConfigMaps(opt.Namespace).Get(context.TODO(), configmap, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf("##### CONFIGMAP - %s #####\n", configmap)
	for k, v := range resp.Data {
		res += parser.ParseStr(!opt.NoExport, k, v)
	}

	return res, nil
}
