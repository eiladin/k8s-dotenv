package secret

import (
	"context"
	"fmt"

	"github.com/eiladin/k8s-dotenv/internal/client"
	"github.com/eiladin/k8s-dotenv/internal/parser"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Get(namespace string, name string) (string, error) {
	clientset, err := client.Get()
	if err != nil {
		return "", err
	}

	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf("##### SECRET - %s #####\n", name)
	for k, v := range secret.Data {
		res += parser.Parse(k, v)
	}

	return res, nil
}
