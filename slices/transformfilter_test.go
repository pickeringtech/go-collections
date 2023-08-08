package slices_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/slices"
	"reflect"
	"testing"
)

func ExampleFilter() {
	input := []int{1, 2, 3, 4, 5}
	output := slices.Filter(input, func(element int) bool {
		return element > 2
	})
	fmt.Printf("Output: %v\n", output)

	// Output: Output: [3 4 5]
}

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

func BenchmarkFilter(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
		fn   func(int) bool
	}{
		{
			name: "3 elements",
			sli:  []int{1, 2, 3},
			fn: func(element int) bool {
				return element >= 2
			},
		},
		{
			name: "10 elements",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn: func(element int) bool {
				return element >= 5
			},
		},
		{
			name: "100 elements",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn: func(element int) bool {
				return element >= 50
			},
		},
		{
			name: "1_000 elements",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn: func(element int) bool {
				return element >= 500
			},
		},
		{
			name: "10_000 elements",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn: func(element int) bool {
				return element >= 5_000
			},
		},
		{
			name: "100_000 elements",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn: func(element int) bool {
				return element >= 50_000
			},
		},
		{
			name: "1_000_000 elements",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn: func(element int) bool {
				return element >= 500_000
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Filter(bm.sli, bm.fn)
			}
		})
	}
}
