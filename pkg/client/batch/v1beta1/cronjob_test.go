package v1beta1

import (
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"github.com/eiladin/k8s-dotenv/pkg/result"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestBatchV1Beta1_CronJob(t *testing.T) {
	mockv1 := mock.CronJobv1beta1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)
	errorClient := mock.NewFakeClient().PrependReactor("get", "cronjobs", true, nil, assert.AnError)

	type args struct {
		resource string
	}
	tests := []struct {
		name         string
		batchv1beta1 *BatchV1Beta1
		args         args
		want         *result.Result
	}{
		{
			name:         "return cronjobs",
			batchv1beta1: NewBatchV1Beta1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
			args:         args{resource: "test"},
			want: &result.Result{
				Environment: result.EnvValues{"k": "v"},
				Secrets:     map[string]result.EnvValues{"test": {"k": "v"}},
				ConfigMaps:  map[string]result.EnvValues{"test": {"k": "v"}},
			},
		},
		{
			name:         "return API errors",
			batchv1beta1: NewBatchV1Beta1(errorClient, &clientoptions.Clientoptions{Namespace: "test"}),
			args:         args{resource: "test"},
			want:         result.NewFromError(NewResourceLoadError("CronJob", assert.AnError)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.batchv1beta1.CronJob(tt.args.resource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BatchV1Beta1.CronJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBatchV1Beta1_CronJobList(t *testing.T) {
	mockv1 := mock.CronJobv1beta1("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	kubeClient := mock.NewFakeClient(mockv1)
	errorClient := mock.NewFakeClient().PrependReactor("list", "cronjobs", true, nil, assert.AnError)

	tests := []struct {
		name         string
		batchv1beta1 *BatchV1Beta1
		want         []string
		wantErr      bool
	}{
		{
			name:         "return cronjobs",
			batchv1beta1: NewBatchV1Beta1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
			want:         []string{"test"},
		},
		{
			name:         "return API errors",
			batchv1beta1: NewBatchV1Beta1(errorClient, &clientoptions.Clientoptions{Namespace: "test"}),
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.batchv1beta1.CronJobList()
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchV1Beta1.CronJobList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BatchV1Beta1.CronJobList() = %v, want %v", got, tt.want)
			}
		})
	}
}
