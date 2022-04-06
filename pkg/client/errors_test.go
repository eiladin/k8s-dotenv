package client

import "testing"

func Test_newMissingKubeClientError(t *testing.T) {
	type args struct {
		client string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "wraps error", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := newMissingKubeClientError(tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("newMissingKubeClientError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
