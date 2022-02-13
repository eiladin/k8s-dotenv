package configmap

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

	configmap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf("##### CONFIGMAP - %s #####\n", name)
	for k, v := range configmap.Data {
		res += fmt.Sprintf("export %s=\"%s\"\n", strings.ReplaceAll(k, ".", ""), strings.ReplaceAll(string(v), "\n", "\\n"))
	}

	return res, nil
}
