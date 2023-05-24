package slices_test

import (
	"github.com/pickeringtech/go-collectionutil/slices"
	"reflect"
	"testing"
)

func TestFilter_Strings(t *testing.T) {
	type args struct {
		input []string
		fun   slices.FilterFunc[string]
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "filters input when length is below certain level",
			args: args{
				input: []string{"a", "ab", "abc", "abcd"},
				fun: func(element string) bool {
					return len(element) > 2
				},
			},
			want: []string{"abc", "abcd"},
		},
		{
			name: "nil input results in nil output",
			args: args{
				input: nil,
				fun: func(element string) bool {
					return len(element) > 2
				},
			},
			want: nil,
		},
		{
			name: "empty input results in nil output",
			args: args{
				input: []string{},
				fun: func(element string) bool {
					return len(element) > 2
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Filter(tt.args.input, tt.args.fun)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
