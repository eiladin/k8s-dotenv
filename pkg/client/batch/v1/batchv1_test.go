package v1

import (
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"k8s.io/client-go/kubernetes"
)

func TestNewBatchV1(t *testing.T) {
	client := mock.NewFakeClient()

	type args struct {
		client  kubernetes.Interface
		options *options.Client
	}

	tests := []struct {
		name string
		args args
		want *BatchV1
	}{
		{
			name: "create batchV1 client",
			args: args{
				client:  client,
				options: &options.Client{},
			},
			want: NewBatchV1(client, &options.Client{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBatchV1(tt.args.client, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBatchV1() = %v, want %v", got, tt.want)
			}
		})
	}
}
