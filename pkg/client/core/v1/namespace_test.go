package v1

import (
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
)

func TestCoreV1_NamespaceList(t *testing.T) {
	oneNamespaceClient := mock.NewFakeClient(mock.Namespace("one"))
	twoNamespaceClient := mock.NewFakeClient(mock.Namespace("one"), mock.Namespace("two"))
	errorClient := mock.NewFakeClient().PrependReactor("list", "namespaces", true, nil, mock.AnError)

	tests := []struct {
		name    string
		corev1  *CoreV1
		want    []string
		wantErr bool
	}{
		{
			name: "return a single namespace",
			corev1: NewCoreV1(oneNamespaceClient,
				&clientoptions.Clientoptions{}),
			want: []string{"one"},
		},
		{
			name: "return multiple namespaces",
			corev1: NewCoreV1(twoNamespaceClient,
				&clientoptions.Clientoptions{}),
			want: []string{"one", "two"},
		},
		{
			name: "return multiple namespaces",
			corev1: NewCoreV1(errorClient,
				&clientoptions.Clientoptions{}),
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := testCase.corev1.NamespaceList()
			if (err != nil) != testCase.wantErr {
				t.Errorf("CoreV1.NamespaceList() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("CoreV1.NamespaceList() = %v, want %v", got, testCase.want)
			}
		})
	}
}
