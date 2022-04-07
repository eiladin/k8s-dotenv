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

func TestBatchV1_CronJob(t *testing.T) {
	mockv1 := mock.CronJobv1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)
	errorClient := mock.NewFakeClient().PrependReactor("get", "cronjobs", true, nil, mock.AnError)

	type args struct {
		resource string
	}

	tests := []struct {
		name    string
		batchv1 *BatchV1
		args    args
		want    *result.Result
	}{
		{
			name:    "return cronjobs",
			batchv1: NewBatchV1(kubeClient, &options.Client{Namespace: "test"}),
			args:    args{resource: "test"},
			want: &result.Result{
				Environment: result.EnvValues{"k": "v"},
				Secrets:     map[string]result.EnvValues{"test": {"k": "v"}},
				ConfigMaps:  map[string]result.EnvValues{"test": {"k": "v"}},
			},
		},
		{
			name:    "return API errors",
			batchv1: NewBatchV1(errorClient, &options.Client{}),
			args:    args{resource: "test"},
			want:    result.NewFromError(mock.AnError),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			opts := []cmp.Option{
				cmp.AllowUnexported(result.Result{}),
				cmpopts.EquateErrors(),
			}

			if got := testCase.batchv1.CronJob(testCase.args.resource); !cmp.Equal(got, testCase.want, opts...) {
				t.Errorf("BatchV1.CronJob() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestBatchV1_CronJobList(t *testing.T) {
	mockv1 := mock.CronJobv1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	kubeClient := mock.NewFakeClient(mockv1)
	errorClient := mock.NewFakeClient().PrependReactor("list", "cronjobs", true, nil, mock.AnError)

	tests := []struct {
		name    string
		batchv1 *BatchV1
		want    []string
		wantErr bool
	}{
		{
			name:    "return cronjobs",
			batchv1: NewBatchV1(kubeClient, &options.Client{Namespace: "test"}),
			want:    []string{"test"},
		},
		{
			name:    "return API errors",
			batchv1: NewBatchV1(errorClient, &options.Client{Namespace: "test"}),
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := testCase.batchv1.CronJobList()
			if (err != nil) != testCase.wantErr {
				t.Errorf("BatchV1.CronJobList() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("BatchV1.CronJobList() = %v, want %v", got, testCase.want)
			}
		})
	}
}
