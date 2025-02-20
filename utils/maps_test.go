package utils

import (
	"reflect"
	"testing"
)

func TestMapKeys(t *testing.T) {
	type args struct {
		m map[string]string
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "One key", args: args{m: map[string]string{"k": "v"}}, want: []string{"k"}},
		{name: "Two keys", args: args{m: map[string]string{"k": "v", "k1": "v1"}}, want: []string{"k", "k1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapKeys(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
