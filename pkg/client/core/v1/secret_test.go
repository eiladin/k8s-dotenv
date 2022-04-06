package v1

import (
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
)

func TestCoreV1_SecretData(t *testing.T) {
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	kubeClient := mock.NewFakeClient(mockSecret)
	errorClient := mock.NewFakeClient(mockSecret).PrependReactor("get", "secrets", true, nil, mock.AnError)

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
			name:   "return secret data",
			corev1: NewCoreV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
			args:   args{resource: "test"},
			want:   map[string]string{"k": "v"},
		},
		{
			name:    "return API errors",
			corev1:  NewCoreV1(errorClient, &clientoptions.Clientoptions{Namespace: "test"}),
			args:    args{resource: "test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.corev1.SecretData(tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("CoreV1.SecretData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CoreV1.SecretData() = %v, want %v", got, tt.want)
			}
		})
	}
}
