package slices

import (
	"reflect"
	"strconv"
	"testing"
)

func TestGenerate(t *testing.T) {
	type args[T any] struct {
		amount int
		fn     GeneratorFunc[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[string]{
		{
			name: "generates correctly",
			args: args[string]{
				amount: 3,
				fn: func(index int) string {
					return strconv.Itoa(index * 3)
				},
			},
			want: []string{"0", "3", "6"},
		},
		{
			name: "amount 0 provides nil output",
			args: args[string]{
				amount: 0,
				fn: func(index int) string {
					return strconv.Itoa(index * 3)
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Generate(tt.args.amount, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
