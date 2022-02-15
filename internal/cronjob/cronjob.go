package cronjob

import (
	"github.com/eiladin/k8s-dotenv/internal/client"
	"github.com/eiladin/k8s-dotenv/internal/environment"
)

func Get(namespace string, name string) (*environment.Result, error) {
	clientset, err := client.Get()
	if err != nil {
		return nil, err
	}

	apiVersion, err := client.GetApiGroup("CronJob")
	if err != nil {
		return nil, err
	}
	beta1 := apiVersion == "batch/v1beta1"

	if beta1 {
		return GetV1beta1(clientset, namespace, name)
	} else {
		return GetV1(clientset, namespace, name)
	}
}

func GetList(namespace string) ([]string, error) {
	clientset, err := client.Get()
	if err != nil {
		return nil, err
	}

	apiVersion, err := client.GetApiGroup("CronJob")
	if err != nil {
		return nil, err
	}
	beta1 := apiVersion == "batch/v1beta1"

	if beta1 {
		return GetListV1beta1(clientset, namespace)
	} else {
		return GetListV1(clientset, namespace)
	}
}
