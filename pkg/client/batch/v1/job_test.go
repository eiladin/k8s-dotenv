package v1

import (
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"github.com/eiladin/k8s-dotenv/pkg/result"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
)

func TestBatchV1_Job(t *testing.T) {
	mockv1 := mock.Job("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	mockSecret := mock.Secret("test", "test", map[string][]byte{"k": []byte("v")})
	mockConfigMap := mock.ConfigMap("test", "test", map[string]string{"k": "v"})
	kubeClient := mock.NewFakeClient(mockv1, mockConfigMap, mockSecret)
	errorClient := mock.NewFakeClient().PrependReactor("get", "jobs", true, nil, mock.AnError)

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
			name:    "return jobs",
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
			batchv1: NewBatchV1(errorClient, &clientoptions.Clientoptions{Namespace: "test"}),
			args:    args{resource: "test"},
			want:    result.NewFromError(NewResourceLoadError("Job", mock.AnError)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.batchv1.Job(tt.args.resource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BatchV1.Job() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBatchV1_JobList(t *testing.T) {
	mockv1 := mock.Job("test", "test", map[string]string{"k": "v"}, []string{"test"}, []string{"test"})
	kubeClient := mock.NewFakeClient(mockv1)
	errorClient := mock.NewFakeClient().PrependReactor("list", "jobs", true, nil, mock.AnError)

	tests := []struct {
		name    string
		batchv1 *BatchV1
		want    []string
		wantErr bool
	}{
		{
			name:    "return jobs",
			batchv1: NewBatchV1(kubeClient, &clientoptions.Clientoptions{Namespace: "test"}),
			want:    []string{"test"},
		},
		{
			name:    "return API errors",
			batchv1: NewBatchV1(errorClient, &clientoptions.Clientoptions{Namespace: "test"}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.batchv1.JobList()
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchV1.JobList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BatchV1.JobList() = %v, want %v", got, tt.want)
			}
		})
	}
}
