package v1

import (
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"github.com/eiladin/k8s-dotenv/pkg/result"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestAppsV1_Deployment(t *testing.T) {
	mockv1 := mock.Deployment("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)
	errorClient := mock.NewFakeClient().PrependReactor("get", "deployments", true, nil, assert.AnError)

	type args struct {
		resource string
	}
	tests := []struct {
		name   string
		appsv1 *AppsV1
		args   args
		want   *result.Result
	}{
		{
			name:   "return deployment",
			appsv1: NewAppsV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
			args:   args{resource: "test"},
			want: &result.Result{
				Environment: result.EnvValues{"k": "v"},
				Secrets:     map[string]result.EnvValues{"test": {"k": "v"}},
				ConfigMaps:  map[string]result.EnvValues{"test": {"k": "v"}},
			},
		},
		{
			name:   "return API errors",
			appsv1: NewAppsV1(errorClient, &clientoptions.Clientoptions{Namespace: "test"}),
			args:   args{resource: "test"},
			want:   result.NewFromError(NewResourceLoadError("Deployment", assert.AnError)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.appsv1.Deployment(tt.args.resource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppsV1.Deployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppsV1_DeploymentList(t *testing.T) {
	mockv1 := mock.Deployment("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	kubeClient := mock.NewFakeClient(mockv1)
	errorClient := mock.NewFakeClient().PrependReactor("list", "deployments", true, nil, assert.AnError)

	tests := []struct {
		name    string
		appsv1  *AppsV1
		want    []string
		wantErr bool
	}{
		{
			name:   "return deployments",
			appsv1: NewAppsV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
			want:   []string{"test"},
		},
		{
			name:    "return API errors",
			appsv1:  NewAppsV1(errorClient, &clientoptions.Clientoptions{Namespace: "test"}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.appsv1.DeploymentList()
			if (err != nil) != tt.wantErr {
				t.Errorf("AppsV1.DeploymentList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppsV1.DeploymentList() = %v, want %v", got, tt.want)
			}
		})
	}
}
