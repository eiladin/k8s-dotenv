package client

import (
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"k8s.io/client-go/kubernetes"
)

func TestWithKubeClient(t *testing.T) {
	kubeclient := mock.NewFakeClient()
	type args struct {
		kubeClient kubernetes.Interface
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "update kube client",
			args: args{kubeClient: kubeclient},
			want: &Client{Interface: kubeclient},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := WithKubeClient(tt.args.kubeClient)
			cl := NewClient()
			fn(cl)
			if !reflect.DeepEqual(cl.Interface, tt.want.Interface) {
				t.Errorf("WithKubeClient() = %v, want %v", cl, tt.want)
			}
		})
	}
}

func TestWithExport(t *testing.T) {
	type args struct {
		shouldExport bool
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "update Client ShouldExport",
			args: args{shouldExport: true},
			want: &Client{options: &clientoptions.Clientoptions{ShouldExport: true}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := WithExport(tt.args.shouldExport)
			cl := NewClient()
			fn(cl)
			if !reflect.DeepEqual(cl, tt.want) {
				t.Errorf("WithExport() = %v, want %v", cl, tt.want)
			}
		})
	}
}

func TestWithNamespace(t *testing.T) {
	type args struct {
		namespace string
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "update Client Namespace",
			args: args{namespace: "test"},
			want: &Client{options: &clientoptions.Clientoptions{Namespace: "test"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := WithNamespace(tt.args.namespace)
			cl := NewClient()
			fn(cl)
			if !reflect.DeepEqual(cl, tt.want) {
				t.Errorf("WithNamespace() = %v, want %v", cl, tt.want)
			}
		})
	}
}
