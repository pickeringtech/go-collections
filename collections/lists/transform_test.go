package lists_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pickeringtech/go-collections/collections/lists"
)

func ExampleMap() {
	l := lists.NewArray("a", "ab", "abc", "d")
	lengths := lists.Map(l, func(s string) int {
		return len(s)
	})
	fmt.Printf("%v\n", lengths.AsSlice())

	// Output:
	// [1 2 3 1]
}

func ExampleFlatMap() {
	l := lists.NewArray(1, 2, 3)
	repeated := lists.FlatMap(l, func(n int) lists.List[int] {
		out := []int{}
		for i := 0; i < n; i++ {
			out = append(out, n)
		}
		return lists.NewArray(out...)
	})
	fmt.Printf("%v\n", repeated.AsSlice())

	// Output:
	// [1 2 2 3 3 3]
}

func ExampleReduce() {
	l := lists.NewArray(1, 2, 3, 4)
	sum := lists.Reduce(l, 0, func(acc, n int) int {
		return acc + n
	})
	fmt.Printf("%d\n", sum)

	// Output:
	// 10
}

func TestMap(t *testing.T) {
	tests := []struct {
		name  string
		input lists.List[string]
		fn    func(string) int
		want  []int
	}{
		{
			name:  "maps string lengths to a new element type",
			input: lists.NewArray("a", "ab", "abc", "d"),
			fn:    func(s string) int { return len(s) },
			want:  []int{1, 2, 3, 1},
		},
		{
			name:  "empty input yields non-nil empty output",
			input: lists.NewArray[string](),
			fn:    func(s string) int { return len(s) },
			want:  []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lists.Map(tt.input, tt.fn)
			if !reflect.DeepEqual(got.AsSlice(), tt.want) {
				t.Errorf("Map() = %v, want %v", got.AsSlice(), tt.want)
			}
		})
	}
}

func TestMap_PreservesOrderAcrossImplementations(t *testing.T) {
	// The function takes the List interface, so it works over any
	// implementation without per-impl duplication.
	impls := map[string]lists.List[int]{
		"Array":          lists.NewArray(1, 2, 3),
		"Linked":         lists.NewLinked(1, 2, 3),
		"DoublyLinked":   lists.NewDoublyLinked(1, 2, 3),
		"ConcurrentList": lists.NewConcurrentArray(1, 2, 3),
	}
	want := []int{2, 4, 6}
	for name, impl := range impls {
		t.Run(name, func(t *testing.T) {
			got := lists.Map(impl, func(n int) int { return n * 2 })
			if !reflect.DeepEqual(got.AsSlice(), want) {
				t.Errorf("Map() = %v, want %v", got.AsSlice(), want)
			}
		})
	}
}

func TestFlatMap(t *testing.T) {
	tests := []struct {
		name  string
		input lists.List[int]
		fn    func(int) lists.List[int]
		want  []int
	}{
		{
			name:  "expands each element into n copies",
			input: lists.NewArray(1, 2, 3),
			fn: func(n int) lists.List[int] {
				out := []int{}
				for i := 0; i < n; i++ {
					out = append(out, n)
				}
				return lists.NewArray(out...)
			},
			want: []int{1, 2, 2, 3, 3, 3},
		},
		{
			name:  "empty inner lists drop elements",
			input: lists.NewArray(1, 2, 3),
			fn: func(n int) lists.List[int] {
				if n%2 == 0 {
					return lists.NewArray(n)
				}
				return lists.NewArray[int]()
			},
			want: []int{2},
		},
		{
			name:  "empty input yields non-nil empty output",
			input: lists.NewArray[int](),
			fn:    func(n int) lists.List[int] { return lists.NewArray(n) },
			want:  []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lists.FlatMap(tt.input, tt.fn)
			if !reflect.DeepEqual(got.AsSlice(), tt.want) {
				t.Errorf("FlatMap() = %v, want %v", got.AsSlice(), tt.want)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		name  string
		input lists.List[int]
		init  int
		fn    func(int, int) int
		want  int
	}{
		{
			name:  "sums elements onto init",
			input: lists.NewArray(1, 2, 3, 4),
			init:  10,
			fn:    func(acc, n int) int { return acc + n },
			want:  20,
		},
		{
			name:  "empty input returns init unchanged",
			input: lists.NewArray[int](),
			init:  42,
			fn:    func(acc, n int) int { return acc + n },
			want:  42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lists.Reduce(tt.input, tt.init, tt.fn)
			if got != tt.want {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReduce_DifferentAccumulatorType(t *testing.T) {
	l := lists.NewArray(1, 2, 3)
	got := lists.Reduce(l, "", func(acc string, n int) string {
		return acc + fmt.Sprintf("%d", n)
	})
	if got != "123" {
		t.Errorf("Reduce() = %q, want %q", got, "123")
	}
}
