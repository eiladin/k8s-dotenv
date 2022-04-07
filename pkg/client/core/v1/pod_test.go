package v1

import (
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/result"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
			corev1: NewCoreV1(podClient, &options.Client{Namespace: "test"}),
			args:   args{resource: "test"},
			want: &result.Result{
				Environment: result.EnvValues{"k": "v"},
				Secrets:     map[string]result.EnvValues{"test": {"k": "v"}},
				ConfigMaps:  map[string]result.EnvValues{"test": {"k": "v"}},
			},
		},
		{
			name:   "return API errors",
			corev1: NewCoreV1(errorClient, &options.Client{Namespace: "test"}),
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

			if got := testCase.corev1.Pod(testCase.args.resource); !cmp.Equal(got, testCase.want, opts...) {
				t.Errorf("CoreV1.Pod() = %v, want %v", got, testCase.want)
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
			corev1: NewCoreV1(kubeClient, &options.Client{Namespace: "test"}),
			want:   []string{"test"},
		},
		{
			name:    "return API errors",
			corev1:  NewCoreV1(errorClient, &options.Client{Namespace: "test"}),
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
