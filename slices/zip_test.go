package slices_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/slices"
)

func ExampleZip() {
	names := []string{"alice", "bob"}
	ages := []int{30, 25}
	output := slices.Zip(names, ages)
	fmt.Printf("Output: %v\n", output)

	// Output: Output: [{alice 30} {bob 25}]
}

func ExampleZipWith() {
	a := []int{1, 2, 3}
	b := []int{10, 20, 30}
	output := slices.ZipWith(a, b, func(x, y int) int {
		return x + y
	})
	fmt.Printf("Output: %v\n", output)

	// Output: Output: [11 22 33]
}

func TestZip(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []string
		want []slices.Pair[int, string]
	}{
		{
			name: "equal lengths pair every element",
			a:    []int{1, 2, 3},
			b:    []string{"a", "b", "c"},
			want: []slices.Pair[int, string]{{1, "a"}, {2, "b"}, {3, "c"}},
		},
		{
			name: "truncates to shorter first input",
			a:    []int{1, 2},
			b:    []string{"a", "b", "c"},
			want: []slices.Pair[int, string]{{1, "a"}, {2, "b"}},
		},
		{
			name: "truncates to shorter second input",
			a:    []int{1, 2, 3},
			b:    []string{"a"},
			want: []slices.Pair[int, string]{{1, "a"}},
		},
		{
			name: "nil first input yields non-nil empty output",
			a:    nil,
			b:    []string{"a"},
			want: []slices.Pair[int, string]{},
		},
		{
			name: "nil second input yields non-nil empty output",
			a:    []int{1},
			b:    nil,
			want: []slices.Pair[int, string]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Zip(tt.a, tt.b)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Zip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZipWith(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		fun  slices.ZipFunc[int, int, int]
		want []int
	}{
		{
			name: "sums equal-length inputs",
			a:    []int{1, 2, 3},
			b:    []int{10, 20, 30},
			fun:  func(x, y int) int { return x + y },
			want: []int{11, 22, 33},
		},
		{
			name: "truncates to shorter input",
			a:    []int{1, 2, 3, 4},
			b:    []int{10, 20},
			fun:  func(x, y int) int { return x * y },
			want: []int{10, 40},
		},
		{
			name: "nil input yields non-nil empty output",
			a:    nil,
			b:    []int{1, 2},
			fun:  func(x, y int) int { return x + y },
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.ZipWith(tt.a, tt.b, tt.fun)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZipWith() = %v, want %v", got, tt.want)
			}
		})
	}
}
