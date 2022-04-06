package client

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	v1 "k8s.io/api/batch/v1"
)

func TestClient_AppsV1(t *testing.T) {
	tests := []struct {
		name       string
		client     *Client
		wantNotNil bool
		wantPanic  bool
	}{
		{name: "error", wantPanic: true},
		{name: "create", client: NewClient(WithKubeClient(mock.NewFakeClient())), wantNotNil: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err interface{} = nil
			defer func() {
				if err == nil && tt.wantPanic {
					t.Errorf("Client.AppV1() did not panic")
				} else if err != nil && !tt.wantPanic {
					t.Errorf("Client.AppV1() panicked")
				}
			}()
			defer func() { err = recover() }()

			if got := tt.client.AppsV1(); (got != nil) != tt.wantNotNil {
				t.Errorf("Client.AppsV1() = %v, want %v", got != nil, tt.wantNotNil)
			}
		})
	}
}

func TestClient_BatchV1(t *testing.T) {
	tests := []struct {
		name       string
		client     *Client
		wantNotNil bool
		wantPanic  bool
	}{
		{name: "error", wantPanic: true},
		{name: "create", client: NewClient(WithKubeClient(mock.NewFakeClient())), wantNotNil: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err interface{} = nil
			defer func() {
				if err == nil && tt.wantPanic {
					t.Errorf("Client.BatchV1() did not panic")
				} else if err != nil && !tt.wantPanic {
					t.Errorf("Client.BatchV1() panicked")
				}
			}()
			defer func() {
				t.Log(2)
				err = recover()
			}()

			if got := tt.client.BatchV1(); (got != nil) != tt.wantNotNil {
				t.Errorf("Client.BatchV1() = %v, want %v", got != nil, tt.wantNotNil)
			}
		})
	}
}

func TestClient_BatchV1Beta1(t *testing.T) {
	tests := []struct {
		name       string
		client     *Client
		wantNotNil bool
		wantPanic  bool
	}{
		{name: "error", wantPanic: true},
		{name: "create", client: NewClient(WithKubeClient(mock.NewFakeClient())), wantNotNil: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err interface{} = nil
			defer func() {
				if err == nil && tt.wantPanic {
					t.Errorf("Client.BatchV1Beta1() did not panic")
				} else if err != nil && !tt.wantPanic {
					t.Errorf("Client.BatchV1Beta1() panicked")
				}
			}()
			defer func() { err = recover() }()

			if got := tt.client.BatchV1Beta1(); (got != nil) != tt.wantNotNil {
				t.Errorf("Client.BatchV1Beta1() = %v, want %v", got != nil, tt.wantNotNil)
			}
		})
	}
}

func TestClient_CoreV1(t *testing.T) {
	tests := []struct {
		name       string
		client     *Client
		wantNotNil bool
		wantPanic  bool
	}{
		{name: "error", wantPanic: true},
		{name: "create", client: NewClient(WithKubeClient(mock.NewFakeClient())), wantNotNil: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err interface{} = nil
			defer func() {
				if err == nil && tt.wantPanic {
					t.Errorf("Client.CoreV1() did not panic")
				} else if err != nil && !tt.wantPanic {
					t.Errorf("Client.CoreV1() panicked")
				}
			}()
			defer func() { err = recover() }()

			if got := tt.client.CoreV1(); (got != nil) != tt.wantNotNil {
				t.Errorf("Client.CoreV1() = %v, want %v", got != nil, tt.wantNotNil)
			}
		})
	}
}

func TestClient_GetAPIGroup(t *testing.T) {
	kubeClient := mock.NewFakeClient(&v1.Job{}).WithResources(mock.Jobv1Resource())
	missingResourceClient := mock.NewFakeClient(&v1.Job{})
	errorClient := mock.NewFakeClient(&v1.Job{}).WithResources(mock.InvalidGroupResource())

	type args struct {
		resource string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    string
		wantErr bool
	}{
		{name: "detect resource group", client: NewClient(WithKubeClient(kubeClient)), args: args{resource: "Job"}, want: "v1"},
		{name: "error if resource not found", client: NewClient(WithKubeClient(missingResourceClient)), args: args{resource: "Job"}, wantErr: true},
		{name: "return API errors", client: NewClient(WithKubeClient(errorClient)), args: args{resource: "Job"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.GetAPIGroup(tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetAPIGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.GetAPIGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		configures []ConfigureFunc
	}
	tests := []struct {
		name       string
		args       args
		wantNotNil bool
	}{
		{name: "run configures", args: args{configures: []ConfigureFunc{WithKubeClient(mock.NewFakeClient())}}, wantNotNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewClient(tt.args.configures...)
			if (got.appsv1 != nil) != tt.wantNotNil {
				t.Errorf("NewClient().appsv1 is nil, want not nil")
			}
			if (got.batchv1 != nil) != tt.wantNotNil {
				t.Errorf("NewClient().batchv1 is nil, want not nil")
			}
			if (got.batchv1beta1 != nil) != tt.wantNotNil {
				t.Errorf("NewClient().batchv1beta1 is nil, want not nil")
			}
			if (got.corev1 != nil) != tt.wantNotNil {
				t.Errorf("NewClient().corev1 is nil, want not nil")
			}
			if (got.options != nil) != tt.wantNotNil {
				t.Errorf("NewClient().options is nil, want not nil")
			}
			if (got.Interface != nil) != tt.wantNotNil {
				t.Errorf("NewClient().Interface is nil, want not nil")
			}
		})
	}
}
