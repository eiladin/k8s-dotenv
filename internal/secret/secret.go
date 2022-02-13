package secret

import (
	"context"
	"fmt"
	"strings"

	"github.com/eiladin/k8s-dotenv/internal/client"
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
		res += fmt.Sprintf("export %s=\"%s\"\n", strings.ReplaceAll(k, ".", ""), strings.ReplaceAll(string(v), "\n", "\\n"))
	}

	return res, nil
}
