package v1

import (
	"reflect"
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/clientoptions"
	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
	"github.com/stretchr/testify/assert"
)

func TestCoreV1_NamespaceList(t *testing.T) {
	oneNamespaceClient := mock.NewFakeClient(mock.Namespace("one"))
	twoNamespaceClient := mock.NewFakeClient(mock.Namespace("one"), mock.Namespace("two"))
	errorClient := mock.NewFakeClient().PrependReactor("list", "namespaces", true, nil, assert.AnError)

	tests := []struct {
		name    string
		corev1  *CoreV1
		want    []string
		wantErr bool
	}{
		{name: "return a single namespace", corev1: NewCoreV1(oneNamespaceClient, &clientoptions.Clientoptions{}), want: []string{"one"}},
		{name: "return multiple namespaces", corev1: NewCoreV1(twoNamespaceClient, &clientoptions.Clientoptions{}), want: []string{"one", "two"}},
		{name: "return multiple namespaces", corev1: NewCoreV1(errorClient, &clientoptions.Clientoptions{}), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.corev1.NamespaceList()
			if (err != nil) != tt.wantErr {
				t.Errorf("CoreV1.NamespaceList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CoreV1.NamespaceList() = %v, want %v", got, tt.want)
			}
		})
	}
}
