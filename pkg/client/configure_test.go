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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			fn := WithKubeClient(testCase.args.kubeClient)
			cl := NewClient()
			fn(cl)
			if !reflect.DeepEqual(cl.Interface, testCase.want.Interface) {
				t.Errorf("WithKubeClient() = %v, want %v", cl, testCase.want)
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			fn := WithExport(testCase.args.shouldExport)
			cl := NewClient()
			fn(cl)
			if !reflect.DeepEqual(cl, testCase.want) {
				t.Errorf("WithExport() = %v, want %v", cl, testCase.want)
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			fn := WithNamespace(testCase.args.namespace)
			cl := NewClient()
			fn(cl)
			if !reflect.DeepEqual(cl, testCase.want) {
				t.Errorf("WithNamespace() = %v, want %v", cl, testCase.want)
			}
		})
	}
}
