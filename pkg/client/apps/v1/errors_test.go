package v1

import (
	"testing"

	"github.com/eiladin/k8s-dotenv/pkg/testing/mock"
)

func TestResourceLoadError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *ResourceLoadError
		want string
	}{
		{
			name: "return internal error",
			e: &ResourceLoadError{
				Err:      mock.AnError,
				Resource: "test",
			},
			want: "error loading test: mock.AnError general error for testing",
		},
		{
			name: "return message when there is no internal error",
			e:    &ResourceLoadError{Resource: "test"},
			want: "error loading test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("ResourceLoadError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceLoadError_Unwrap(t *testing.T) {
	tests := []struct {
		name    string
		e       *ResourceLoadError
		wantErr error
	}{
		{
			name: "return internal error",
			e: &ResourceLoadError{
				Err:      mock.AnError,
				Resource: "test",
			},
			wantErr: mock.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Unwrap(); err != tt.wantErr {
				t.Errorf("ResourceLoadError.Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewResourceLoadError(t *testing.T) {
	type args struct {
		resource string
		err      error
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "wrap errors",
			args: args{
				resource: "test",
				err:      mock.AnError,
			},
			wantErr: &ResourceLoadError{Resource: "test", Err: mock.AnError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewResourceLoadError(tt.args.resource, tt.args.err); err.Error() != tt.wantErr.Error() {
				t.Errorf("NewResourceLoadError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
