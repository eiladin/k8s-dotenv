package v1

import (
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"github.com/eiladin/k8s-dotenv/pkg/result"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestAppsV1_Deployment(t *testing.T) {
	mockv1 := mock.Deployment("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)
	errorClient := mock.NewFakeClient().PrependReactor("get", "deployments", true, nil, mock.AnError)

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
			want:   result.NewFromError(mock.AnError),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			opts := []cmp.Option{
				cmp.AllowUnexported(result.Result{}),
				cmpopts.EquateErrors(),
			}

			if got := testCase.appsv1.Deployment(testCase.args.resource); !cmp.Equal(got, testCase.want, opts...) {
				t.Errorf("AppsV1.Deployment() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestAppsV1_DeploymentList(t *testing.T) {
	mockv1 := mock.Deployment("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	kubeClient := mock.NewFakeClient(mockv1)
	errorClient := mock.NewFakeClient().PrependReactor("list", "deployments", true, nil, mock.AnError)

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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := testCase.appsv1.DeploymentList()
			if (err != nil) != testCase.wantErr {
				t.Errorf("AppsV1.DeploymentList() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("AppsV1.DeploymentList() = %v, want %v", got, testCase.want)
			}
		})
	}
}
