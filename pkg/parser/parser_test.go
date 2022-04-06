package parser

import "testing"

func TestParse(t *testing.T) {
	type args struct {
		shouldExport bool
		key          string
		value        []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "without export", args: args{key: "key", value: []byte("value")}, want: "key=\"value\"\n"},
		{name: "with export", args: args{shouldExport: true, key: "key", value: []byte("value")}, want: "export key=\"value\"\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args.shouldExport, tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseStr(t *testing.T) {
	type args struct {
		shouldExport bool
		key          string
		value        string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "without export", args: args{key: "key", value: "value"}, want: "key=\"value\"\n"},
		{name: "with export", args: args{shouldExport: true, key: "key", value: "value"}, want: "export key=\"value\"\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseStr(tt.args.shouldExport, tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("ParseStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
