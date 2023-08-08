package slices_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/slices"
	"reflect"
	"testing"
)

func ExampleReduce() {
	a := []int{1, 2, 3, 4, 5}
	b := slices.Reduce(a, slices.TotalReducer[int])
	fmt.Printf("total: %v\n", b)

	// Output:
	// total: 15
}

func TestReduce(t *testing.T) {
	type args[I any, O any] struct {
		input []I
		fn    slices.ReductionFunc[I, O]
	}
	type testCase[I any, O any] struct {
		name string
		args args[I, O]
		want O
	}
	tests := []testCase[int, int]{
		{
			name: "reduce can add up all inputs",
			args: args[int, int]{
				input: []int{1, 2, 3, 4, 5},
				fn:    slices.TotalReducer[int],
			},
			want: 15,
		},
		{
			name: "empty input provides zero value output",
			args: args[int, int]{
				input: []int{},
				fn:    slices.TotalReducer[int],
			},
			want: 0,
		},
		{
			name: "nil input provides zero value output",
			args: args[int, int]{
				input: []int{},
				fn:    slices.TotalReducer[int],
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Reduce(tt.args.input, tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkReduce(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
		fn   slices.ReductionFunc[int, int]
	}{
		{
			name: "3 elements",
			sli:  []int{1, 2, 3},
			fn:   slices.TotalReducer[int],
		},
		{
			name: "10 elements",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn:   slices.TotalReducer[int],
		},
		{
			name: "100 elements",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn:   slices.TotalReducer[int],
		},
		{
			name: "1_000 elements",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn:   slices.TotalReducer[int],
		},
		{
			name: "10_000 elements",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn:   slices.TotalReducer[int],
		},
		{
			name: "100_000 elements",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn:   slices.TotalReducer[int],
		},
		{
			name: "1_000_000 elements",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn:   slices.TotalReducer[int],
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Reduce(bm.sli, bm.fn)
			}
		})
	}
}

func ExampleNewCountOccurrencesReducer() {
	a := []int{1, 2, 3, 4, 5, 3, 2, 2, 5, 4, 1}
	b := slices.Reduce(a, slices.NewCountOccurrencesReducer[int, int]([]int{1, 2, 3}))
	fmt.Printf("occurrences: %v\n", b)

	// Output:
	// occurrences: 7
}

func TestReduce_CountOccurrences(t *testing.T) {
	type args[I any, O any] struct {
		input []I
	}
	type testCase[I any, O any] struct {
		name string
		args args[I, O]
		want O
	}
	tests := []testCase[int, int]{
		{
			name: "reduce can add up all inputs",
			args: args[int, int]{
				input: []int{1, 2, 3, 4, 5, 3, 2, 2, 5, 4, 1},
			},
			want: 7,
		},
		{
			name: "empty input provides zero value output",
			args: args[int, int]{
				input: []int{},
			},
			want: 0,
		},
		{
			name: "nil input provides zero value output",
			args: args[int, int]{
				input: []int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Reduce[int, int](tt.args.input, slices.NewCountOccurrencesReducer[int, int]([]int{1, 2, 3}))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkReduce_CountOccurrences(b *testing.B) {
	benchmarks := []struct {
		name   string
		sli    []int
		values []int
	}{
		{
			name:   "3 elements",
			sli:    []int{1, 2, 3},
			values: []int{2, 3},
		},
		{
			name:   "10 elements",
			sli:    slices.Generate(10, slices.NumericIdentityGenerator[int]),
			values: []int{2, 3},
		},
		{
			name:   "100 elements",
			sli:    slices.Generate(100, slices.NumericIdentityGenerator[int]),
			values: []int{2, 3},
		},
		{
			name:   "1_000 elements",
			sli:    slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			values: []int{2, 3},
		},
		{
			name:   "10_000 elements",
			sli:    slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			values: []int{2, 3},
		},
		{
			name:   "100_000 elements",
			sli:    slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			values: []int{2, 3},
		},
		{
			name:   "1_000_000 elements",
			sli:    slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			values: []int{2, 3},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Reduce(bm.sli, slices.NewCountOccurrencesReducer[int, int](bm.values))
			}
		})
	}
}
