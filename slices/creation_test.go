package slices_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/constraints"
	"github.com/pickeringtech/go-collections/slices"
	"reflect"
	"strconv"
	"testing"
)

func ExampleGenerate() {
	type numberAndSquare struct {
		number int
		square int
	}

	sli := slices.Generate[numberAndSquare](10, func(index int) numberAndSquare {
		return numberAndSquare{
			number: index,
			square: index * index,
		}
	})

	fmt.Printf("%v", sli)
	// Output: [{0 0} {1 1} {2 4} {3 9} {4 16} {5 25} {6 36} {7 49} {8 64} {9 81}]
}

func TestGenerate(t *testing.T) {
	type args[T any] struct {
		amount int
		fn     slices.GeneratorFunc[T]
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
			got := slices.Generate(tt.args.amount, tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkGenerate(b *testing.B) {
	benchmarks := []struct {
		name   string
		amount int
	}{
		{
			name:   "generates 10 data points",
			amount: 10,
		},
		{
			name:   "generates 100 data points",
			amount: 100,
		},
		{
			name:   "generates 1000 data points",
			amount: 1000,
		},
		{
			name:   "generates 10000 data points",
			amount: 10000,
		},
		{
			name:   "generates 100000 data points",
			amount: 100000,
		},
		{
			name:   "generates 1000000 data points",
			amount: 1000000,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				slices.Generate(bm.amount, func(index int) int {
					return index * index
				})
			}
		})
	}
}

func TestNumericIdentityGenerator(t *testing.T) {
	type args struct {
		index int
	}
	type testCase[T constraints.Numeric] struct {
		name string
		args args
		want T
	}
	tests := []testCase[int]{
		{
			name: "outputs the index",
			args: args{
				index: 5,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.NumericIdentityGenerator[int](tt.args.index)
			if got != tt.want {
				t.Errorf("NumericIdentityGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}
