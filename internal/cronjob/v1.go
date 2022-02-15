package cronjob

import (
	"context"

	"github.com/eiladin/k8s-dotenv/internal/environment"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetV1(clientset *kubernetes.Clientset, namespace string, name string) (*environment.Result, error) {
	res := environment.NewResult()
	resp, err := clientset.BatchV1().CronJobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	containers := resp.Spec.JobTemplate.Spec.Template.Spec.Containers

	for _, cont := range containers {
		for _, env := range cont.Env {
			res.Environment[env.Name] = env.Value
		}

		for _, envFrom := range cont.EnvFrom {
			if envFrom.SecretRef != nil {
				res.Secrets = append(res.Secrets, envFrom.SecretRef.Name)
			}
			if envFrom.ConfigMapRef != nil {
				res.ConfigMaps = append(res.ConfigMaps, envFrom.ConfigMapRef.Name)
			}
		}
	}

	return res, nil
}

func GetListV1(clientset *kubernetes.Clientset, namespace string) ([]string, error) {
	res := []string{}

	resp, err := clientset.BatchV1().CronJobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, item := range resp.Items {
		res = append(res, item.Name)
	}

	return res, nil
}
