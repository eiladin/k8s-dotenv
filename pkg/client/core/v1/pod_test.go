package v1

import (
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"github.com/eiladin/k8s-dotenv/pkg/result"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
)

func TestCoreV1_Pod(t *testing.T) {
	mockv1 := mock.Pod("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})

	podClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)
	errorClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret).
		PrependReactor("get", "pods", true, nil, mock.AnError)

	type args struct {
		resource string
	}

	tests := []struct {
		name   string
		corev1 *CoreV1
		args   args
		want   *result.Result
	}{
		{
			name:   "return pods",
			corev1: NewCoreV1(podClient, &clientoptions.Clientoptions{Namespace: "test"}),
			args:   args{resource: "test"},
			want: &result.Result{
				Environment: result.EnvValues{"k": "v"},
				Secrets:     map[string]result.EnvValues{"test": {"k": "v"}},
				ConfigMaps:  map[string]result.EnvValues{"test": {"k": "v"}},
			},
		},
		{
			name:   "return API errors",
			corev1: NewCoreV1(errorClient, &clientoptions.Clientoptions{Namespace: "test"}),
			args:   args{resource: "test"},
			want:   result.NewFromError(NewResourceLoadError("Pod", mock.AnError)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.corev1.Pod(tt.args.resource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CoreV1.Pod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoreV1_PodList(t *testing.T) {
	mockv1 := mock.Pod("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	kubeClient := mock.NewFakeClient(mockv1)
	errorClient := mock.NewFakeClient(mockv1).PrependReactor("list", "pods", true, nil, mock.AnError)

	tests := []struct {
		name    string
		corev1  *CoreV1
		want    []string
		wantErr bool
	}{
		{
			name:   "return pods",
			corev1: NewCoreV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
			want:   []string{"test"},
		},
		{
			name:    "return API errors",
			corev1:  NewCoreV1(errorClient, &clientoptions.Clientoptions{Namespace: "test"}),
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := testCase.corev1.PodList()
			if (err != nil) != testCase.wantErr {
				t.Errorf("CoreV1.PodList() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("CoreV1.PodList() = %v, want %v", got, testCase.want)
			}
		})
	}
}
