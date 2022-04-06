package v1

import (
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"github.com/eiladin/k8s-dotenv/pkg/result"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
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
			batchv1: NewBatchV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
			args:    args{resource: "test"},
			want: &result.Result{
				Environment: result.EnvValues{"k": "v"},
				Secrets:     map[string]result.EnvValues{"test": {"k": "v"}},
				ConfigMaps:  map[string]result.EnvValues{"test": {"k": "v"}},
			},
		},
		{
			name:    "return API errors",
			batchv1: NewBatchV1(errorClient, &clientoptions.Clientoptions{}),
			args:    args{resource: "test"},
			want:    result.NewFromError(NewResourceLoadError("CronJob", mock.AnError)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.batchv1.CronJob(tt.args.resource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BatchV1.CronJob() = %v, want %v", got, tt.want)
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
			batchv1: NewBatchV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
			want:    []string{"test"},
		},
		{
			name:    "return API errors",
			batchv1: NewBatchV1(errorClient, &clientoptions.Clientoptions{Namespace: "test"}),
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
