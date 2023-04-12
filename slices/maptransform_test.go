package slices

import (
	"reflect"
	"testing"
)

func TestMapConversion_StringToInt(t *testing.T) {
	type args struct {
		input []string
		fun   MapFunc[string, int]
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "counts string lengths",
			args: args{
				input: []string{"a", "ab", "abc", "d"},
				fun: func(element string) int {
					return len(element)
				},
			},
			want: []int{1, 2, 3, 1},
		},
		{
			name: "nil input results in nil output",
			args: args{
				input: nil,
				fun: func(element string) int {
					return len(element)
				},
			},
			want: nil,
		},
		{
			name: "empty input results in nil output",
			args: args{
				input: []string{},
				fun: func(element string) int {
					return len(element)
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Map(tt.args.input, tt.args.fun)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
