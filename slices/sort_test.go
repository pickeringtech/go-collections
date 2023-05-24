package slices

import (
	"github.com/pickeringtech/go-collectionutil/constraints"
	"testing"
)

func TestAscendingSortFunc(t *testing.T) {
	type args[T constraints.Ordered] struct {
		a T
		b T
	}
	type testCase[T constraints.Ordered] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "a < b == true",
			args: args[int]{
				a: 0,
				b: 1,
			},
			want: true,
		},
		{
			name: "a > b == false",
			args: args[int]{
				a: 1,
				b: 0,
			},
			want: false,
		},
		{
			name: "(a < b) < 0",
			args: args[int]{
				a: -2,
				b: -1,
			},
			want: true,
		},
		{
			name: "(a > b) < 0",
			args: args[int]{
				a: -1,
				b: -2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AscendingSortFunc(tt.args.a, tt.args.b)
			if got != tt.want {
				t.Errorf("AscendingSortFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
