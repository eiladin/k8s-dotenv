package v1

import (
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
)

func TestCoreV1_ConfigMapData(t *testing.T) {
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockConfigMap)

	type args struct {
		resource string
	}

	tests := []struct {
		name    string
		corev1  *CoreV1
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name:   "return configmap data",
			corev1: NewCoreV1(kubeClient, &options.Client{Namespace: "test"}),
			args:   args{resource: "test"}, want: map[string]string{"k": "v"},
		},
		{
			name:    "return API errors",
			corev1:  NewCoreV1(kubeClient, &options.Client{Namespace: "test"}),
			args:    args{resource: "test2"},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := testCase.corev1.ConfigMapData(testCase.args.resource)
			if (err != nil) != testCase.wantErr {
				t.Errorf("CoreV1.ConfigMapData() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("CoreV1.ConfigMapData() = %v, want %v", got, testCase.want)
			}
		})
	}
}
