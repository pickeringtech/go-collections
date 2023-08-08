package slices_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/slices"
	"reflect"
	"testing"
)

func ExampleMap() {
	a := []string{"a", "ab", "abc", "d"}
	b := slices.Map(a, func(element string) int {
		return len(element)
	})
	fmt.Printf("%v\n", b)

	// Output:
	// [1 2 3 1]
}

func TestMap_StringToInt(t *testing.T) {
	type args struct {
		input []string
		fun   slices.MapFunc[string, int]
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
			got := slices.Map(tt.args.input, tt.args.fun)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkMap(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []string
		fn   slices.MapFunc[string, int]
	}{
		{
			name: "3 elements",
			sli:  []string{"a", "ab", "abc", "d"},
			fn: func(element string) int {
				return len(element)
			},
		},
		{
			name: "10 elements",
			sli: slices.Generate(10, func(i int) string {
				return "a"
			}),
			fn: func(element string) int {
				return len(element)
			},
		},
		{
			name: "100 elements",
			sli: slices.Generate(100, func(i int) string {
				return "a"
			}),
			fn: func(element string) int {
				return len(element)
			},
		},
		{
			name: "1_000 elements",
			sli: slices.Generate(1_000, func(i int) string {
				return "a"
			}),
			fn: func(element string) int {
				return len(element)
			},
		},
		{
			name: "10_000 elements",
			sli: slices.Generate(10_000, func(i int) string {
				return "a"
			}),
			fn: func(element string) int {
				return len(element)
			},
		},
		{
			name: "100_000 elements",
			sli: slices.Generate(100_000, func(i int) string {
				return "a"
			}),
			fn: func(element string) int {
				return len(element)
			},
		},
		{
			name: "1_000_000 elements",
			sli: slices.Generate(1_000_000, func(i int) string {
				return "a"
			}),
			fn: func(element string) int {
				return len(element)
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Map(bm.sli, bm.fn)
			}
		})
	}
}
