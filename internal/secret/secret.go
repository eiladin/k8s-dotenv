package secret

import (
	"context"
	"fmt"

	"github.com/eiladin/k8s-dotenv/internal/options"
	"github.com/eiladin/k8s-dotenv/internal/parser"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Get(opt *options.Options, secret string) (string, error) {
	resp, err := opt.Client.CoreV1().Secrets(opt.Namespace).Get(context.TODO(), secret, v1.GetOptions{})
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf("##### SECRET - %s #####\n", secret)
	for k, v := range resp.Data {
		res += parser.Parse(!opt.NoExport, k, v)
	}

	return res, nil
}
