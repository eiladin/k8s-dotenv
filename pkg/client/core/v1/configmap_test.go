package v1

import (
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
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
		{name: "return configmap data", corev1: NewCoreV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}), args: args{resource: "test"}, want: map[string]string{"k": "v"}},
		{name: "return API errors", corev1: NewCoreV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}), args: args{resource: "test2"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.corev1.ConfigMapData(tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("CoreV1.ConfigMapData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CoreV1.ConfigMapData() = %v, want %v", got, tt.want)
			}
		})
	}
}
