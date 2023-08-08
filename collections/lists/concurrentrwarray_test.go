package lists_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/collections/lists"
	"github.com/pickeringtech/go-collections/maps"
	"github.com/pickeringtech/go-collections/slices"
	"reflect"
	"testing"
)

func ExampleConcurrentRWArray_AllMatch() {
	a := lists.NewConcurrentRWArray(3, 4)
	match := a.AllMatch(func(a int) bool {
		return a > 2 && a < 5
	})
	fmt.Printf("Matches 1: %v\n", match)

	a = lists.NewConcurrentRWArray(2, 3, 4)
	match = a.AllMatch(func(a int) bool {
		return a > 2 && a < 5
	})
	fmt.Printf("Matches 2: %v\n", match)

	// Output:
	// Matches 1: true
	// Matches 2: false
}

func TestConcurrentRWArray_AllMatch(t *testing.T) {
	type args[T any] struct {
		fn func(T) bool
	}
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "all matches",
			a:    lists.NewConcurrentRWArray(3, 4),
			args: args[int]{
				fn: func(a int) bool {
					return a > 2 && a < 5
				},
			},
			want: true,
		},
		{
			name: "do not all match",
			a:    lists.NewConcurrentRWArray(2, 3, 4),
			args: args[int]{
				fn: func(a int) bool {
					return a > 2 && a < 5
				},
			},
			want: false,
		},
		{
			name: "empty input provides true",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				fn: func(a int) bool {
					return a > 2 && a < 5
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.AllMatch(tt.args.fn)
			if got != tt.want {
				t.Errorf("AllMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleConcurrentRWArray_AnyMatch() {
	arr := lists.NewConcurrentRWArray(4, 5, 3, 1, 2)

	match := arr.AnyMatch(func(a int) bool {
		return a == 3
	})
	fmt.Printf("Matches 1: %v\n", match)

	arr = lists.NewConcurrentRWArray(4, 5, 1, 2)
	match = arr.AnyMatch(func(a int) bool {
		return a == 3
	})
	fmt.Printf("Matches 2: %v\n", match)

	// Output:
	// Matches 1: true
	// Matches 2: false
}

func TestConcurrentRWArray_AnyMatch(t *testing.T) {
	type args[T any] struct {
		fn func(T) bool
	}
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "matches with first element",
			a:    lists.NewConcurrentRWArray(3, 4, 5, 1, 2),
			args: args[int]{
				fn: func(i int) bool {
					return i == 3
				},
			},
			want: true,
		},
		{
			name: "matches with last element",
			a:    lists.NewConcurrentRWArray(4, 5, 1, 2, 3),
			args: args[int]{
				fn: func(i int) bool {
					return i == 3
				},
			},
			want: true,
		},
		{
			name: "matches with middle element",
			a:    lists.NewConcurrentRWArray(4, 5, 3, 1, 2),
			args: args[int]{
				fn: func(i int) bool {
					return i == 3
				},
			},
			want: true,
		},
		{
			name: "no match",
			a:    lists.NewConcurrentRWArray(4, 5, 1, 2),
			args: args[int]{
				fn: func(i int) bool {
					return i == 3
				},
			},
			want: false,
		},
		{
			name: "empty input provides false",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				fn: func(i int) bool {
					return i == 3
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.AnyMatch(tt.args.fn); got != tt.want {
				t.Errorf("AnyMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleConcurrentRWArray_Dequeue() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	first, ok, rest := arr.Dequeue()
	fmt.Printf("First: %v\n", first)
	fmt.Printf("OK: %v\n", ok)
	fmt.Printf("Rest: %v\n", rest)

	// Output:
	// First: 1
	// OK: true
	// Rest: [2 3 4 5]
}

func TestConcurrentRWArray_Dequeue(t *testing.T) {
	type testCase[T any] struct {
		name    string
		a       *lists.ConcurrentRWArray[T]
		want    T
		wantOK  bool
		wantSli []T
	}
	tests := []testCase[int]{
		{
			name:    "dequeues first element",
			a:       lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			want:    1,
			wantOK:  true,
			wantSli: []int{2, 3, 4, 5},
		},
		{
			name:    "dequeueing last element returns true, but nil slice",
			a:       lists.NewConcurrentRWArray[int](1),
			want:    1,
			wantOK:  true,
			wantSli: nil,
		},
		{
			name:    "dequeueing empty input returns false, and nil slice",
			a:       lists.NewConcurrentRWArray[int](),
			want:    0,
			wantOK:  false,
			wantSli: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOK, gotSli := tt.a.Dequeue()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dequeue() got = %v, want %v", got, tt.want)
			}
			if gotOK != tt.wantOK {
				t.Errorf("Dequeue() gotOK = %v, want %v", gotOK, tt.wantOK)
			}
			if !reflect.DeepEqual(gotSli, tt.wantSli) {
				t.Errorf("Dequeue() gotSli = %v, want %v", gotSli, tt.wantSli)
			}
		})
	}
}

func TestConcurrentRWArray_DequeueInPlace(t *testing.T) {
	type testCase[T any] struct {
		name     string
		a        *lists.ConcurrentRWArray[T]
		wantVal  T
		wantOK   bool
		wantRest []T
	}
	tests := []testCase[int]{
		{
			name:     "dequeues first element",
			a:        lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			wantVal:  1,
			wantOK:   true,
			wantRest: []int{2, 3, 4, 5},
		},
		{
			name:     "dequeueing last element returns true, but nil slice",
			a:        lists.NewConcurrentRWArray[int](1),
			wantVal:  1,
			wantOK:   true,
			wantRest: nil,
		},
		{
			name:     "dequeueing empty input returns false, and nil slice",
			a:        lists.NewConcurrentRWArray[int](),
			wantVal:  0,
			wantOK:   false,
			wantRest: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOK := tt.a.DequeueInPlace()
			gotRest := tt.a.GetAsSlice()
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("DequeueInPlace() gotVal = %v, wantVal %v", gotVal, tt.wantVal)
			}
			if gotOK != tt.wantOK {
				t.Errorf("DequeueInPlace() gotOK = %v, wantOK %v", gotOK, tt.wantOK)
			}
			if !reflect.DeepEqual(gotRest, tt.wantRest) {
				t.Errorf("DequeueInPlace() gotRest = %v, wantRest %v", gotRest, tt.wantRest)
			}
		})
	}
}

func ExampleConcurrentRWArray_Enqueue() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	arr.Enqueue(10)
	fmt.Printf("ConcurrentRWArray: %v\n", arr.GetAsSlice())

	// Output:
	// ConcurrentRWArray: [1 2 3 4 5]
}

func TestConcurrentRWArray_Enqueue(t *testing.T) {
	type args[T any] struct {
		element T
	}
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "adds element to end of ConcurrentRWArray",
			a:    lists.NewConcurrentRWArray(1, 2, 3, 4, 5),
			args: args[int]{
				element: 10,
			},
			want: []int{1, 2, 3, 4, 5, 10},
		},
		{
			name: "adding to empty ConcurrentRWArray works",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				element: 10,
			},
			want: []int{10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Enqueue(tt.args.element)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Enqueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentRWArray_EnqueueInPlace(t *testing.T) {
	type args[T any] struct {
		element T
	}
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "adds element to end of ConcurrentRWArray",
			a:    lists.NewConcurrentRWArray(1, 2, 3, 4, 5),
			args: args[int]{
				element: 10,
			},
			want: []int{1, 2, 3, 4, 5, 10},
		},
		{
			name: "adding to empty ConcurrentRWArray works",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				element: 10,
			},
			want: []int{10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.EnqueueInPlace(tt.args.element)

			got := tt.a.GetAsSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnqueueInPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleConcurrentRWArray_Filter() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	out := arr.Filter(func(i int) bool {
		return i > 2 && i < 5
	})
	fmt.Printf("ConcurrentRWArray: %v\n", out)

	// Output:
	// ConcurrentRWArray: [3 4]
}

func TestConcurrentRWArray_Filter(t *testing.T) {
	type args[T any] struct {
		fn func(T) bool
	}
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "filters out values outside range",
			a:    lists.NewConcurrentRWArray(1, 2, 3, 4, 5),
			args: args[int]{
				fn: func(i int) bool {
					return i > 2 && i < 5
				},
			},
			want: []int{3, 4},
		},
		{
			name: "empty input provides nil output",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				fn: func(i int) bool {
					return i > 2 && i < 5
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Filter(tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentRWArray_FilterInPlace(t *testing.T) {
	type args[T any] struct {
		fn func(T) bool
	}
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "filters out values outside range",
			a:    lists.NewConcurrentRWArray(1, 2, 3, 4, 5),
			args: args[int]{
				fn: func(i int) bool {
					return i > 2 && i < 5
				},
			},
			want: []int{3, 4},
		},
		{
			name: "empty input provides nil output",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				fn: func(i int) bool {
					return i > 2 && i < 5
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.FilterInPlace(tt.args.fn)

			got := tt.a.GetAsSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterInPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleConcurrentRWArray_Find() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	out, ok := arr.Find(func(i int) bool {
		return i == 3
	})
	fmt.Printf("Found: %v, %v\n", out, ok)

	// Output:
	// Found: 3, true
}

func TestConcurrentRWArray_Find(t *testing.T) {
	type args[T any] struct {
		fn func(T) bool
	}
	type testCase[T any] struct {
		name   string
		a      *lists.ConcurrentRWArray[T]
		args   args[T]
		want   T
		wantOK bool
	}
	tests := []testCase[int]{
		{
			name: "finds first element",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				fn: func(i int) bool {
					return i == 1
				},
			},
			want:   1,
			wantOK: true,
		},
		{
			name: "finds last element",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				fn: func(i int) bool {
					return i == 5
				},
			},
			want:   5,
			wantOK: true,
		},
		{
			name: "does not find element",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				fn: func(i int) bool {
					return false
				},
			},
			want:   0,
			wantOK: false,
		},
		{
			name: "empty input triggers not found",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				fn: func(i int) bool {
					return true
				},
			},
			want:   0,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.a.Find(tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOK {
				t.Errorf("Find() got1 = %v, want %v", got1, tt.wantOK)
			}
		})
	}
}

func ExampleConcurrentRWArray_FindIndex() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	out := arr.FindIndex(func(i int) bool {
		return i == 3
	})
	fmt.Printf("Index: %v\n", out)

	// Output:
	// Index: 2
}

func TestConcurrentRWArray_FindIndex(t *testing.T) {
	type args[T any] struct {
		fn func(T) bool
	}
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		args args[T]
		want int
	}
	tests := []testCase[int]{
		{
			name: "finds first element",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				fn: func(i int) bool {
					return i == 1
				},
			},
			want: 0,
		},
		{
			name: "finds last element",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				fn: func(i int) bool {
					return i == 5
				},
			},
			want: 4,
		},
		{
			name: "does not find element",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				fn: func(i int) bool {
					return false
				},
			},
			want: -1,
		},
		{
			name: "empty input triggers not found",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				fn: func(i int) bool {
					return true
				},
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.FindIndex(tt.args.fn); got != tt.want {
				t.Errorf("FindIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleConcurrentRWArray_ForEach() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	arr.ForEach(func(element int) {
		fmt.Printf("%v\n", element)
	})

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func TestConcurrentRWArray_ForEach(t *testing.T) {
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		want []int
	}
	tests := []testCase[int]{
		{
			name: "iterates over each element in order",
			a:    lists.NewConcurrentRWArray(1, 2, 3, 4, 5),
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "does not iterate over empty input",
			a:    lists.NewConcurrentRWArray[int](),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []int

			tt.a.ForEach(func(element int) {
				got = append(got, element)
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ForEach() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleConcurrentRWArray_ForEachWithIndex() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	arr.ForEachWithIndex(func(index int, element int) {
		fmt.Printf("%v: %v\n", index, element)
	})

	// Output:
	// 0: 1
	// 1: 2
	// 2: 3
	// 3: 4
	// 4: 5
}

func TestConcurrentRWArray_ForEachWithIndex(t *testing.T) {
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		want []maps.Entry[int, int]
	}
	tests := []testCase[int]{
		{
			name: "iterates over each element in order",
			a:    lists.NewConcurrentRWArray(1, 2, 3, 4, 5),
			want: []maps.Entry[int, int]{
				{Key: 0, Value: 1},
				{Key: 1, Value: 2},
				{Key: 2, Value: 3},
				{Key: 3, Value: 4},
				{Key: 4, Value: 5},
			},
		},
		{
			name: "does not iterate over empty input",
			a:    lists.NewConcurrentRWArray[int](),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []maps.Entry[int, int]

			tt.a.ForEachWithIndex(func(idx, element int) {
				got = append(got, maps.Entry[int, int]{Key: idx, Value: element})
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ForEachWithIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleConcurrentRWArray_Get() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	fmt.Printf("%v\n", arr.Get(2, -1))
	fmt.Printf("%v\n", arr.Get(-1, -1))

	// Output:
	// 3
	// -1
}

func TestConcurrentRWArray_Get(t *testing.T) {
	type args[T any] struct {
		index        int
		defaultValue T
	}
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "retrieves element at index",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				index:        2,
				defaultValue: -1,
			},
			want: 3,
		},
		{
			name: "returns default at index -1",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				index:        -1,
				defaultValue: -1,
			},
			want: -1,
		},
		{
			name: "returns default at index 5",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				index:        5,
				defaultValue: -1,
			},
			want: -1,
		},
		{
			name: "returns default for empty input",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				index:        0,
				defaultValue: -1,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Get(tt.args.index, tt.args.defaultValue)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleConcurrentRWArray_GetAsSlice() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	fmt.Printf("%v\n", arr.GetAsSlice())

	// Output:
	// [1 2 3 4 5]
}

func TestConcurrentRWArray_GetAsSlice(t *testing.T) {
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "converts input to slice",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "empty input provides nil output",
			a:    lists.NewConcurrentRWArray[int](),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.GetAsSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAsSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleConcurrentRWArray_Insert() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	out := arr.Insert(2, 6, 7, 8)

	fmt.Printf("%v\n", out)

	// Output:
	// [1 2 6 7 8 3 4 5]
}

func TestConcurrentRWArray_Insert(t *testing.T) {
	type args[T any] struct {
		index    int
		elements []T
	}
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "inserts elements at specified index",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				index:    2,
				elements: []int{6, 7, 8},
			},
			want: []int{1, 2, 6, 7, 8, 3, 4, 5},
		},
		{
			name: "inserting empty elements slice does nothing",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				index:    2,
				elements: []int{},
			},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "inserting nil elements slice does nothing",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				index:    2,
				elements: nil,
			},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "inserting into empty ConcurrentRWArray yields nil output",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				index:    0,
				elements: []int{6, 7, 8},
			},
			want: nil,
		},
		{
			name: "empty ConcurrentRWArray and empty elements slice yields nil output",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				index:    2,
				elements: []int{},
			},
			want: nil,
		},
		{
			name: "empty ConcurrentRWArray and nil elements slice yields nil output",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				index:    2,
				elements: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Insert(tt.args.index, tt.args.elements...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentRWArray_InsertInPlace(t *testing.T) {
	type args[T any] struct {
		index    int
		elements []T
	}
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "inserts elements at specified index",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				index:    2,
				elements: []int{6, 7, 8},
			},
			want: []int{1, 2, 6, 7, 8, 3, 4, 5},
		},
		{
			name: "inserting empty elements slice does nothing",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				index:    2,
				elements: []int{},
			},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "inserting nil elements slice does nothing",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				index:    2,
				elements: nil,
			},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "inserting into empty ConcurrentRWArray yields nil output",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				index:    0,
				elements: []int{6, 7, 8},
			},
			want: nil,
		},
		{
			name: "empty ConcurrentRWArray and empty elements slice yields nil output",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				index:    2,
				elements: []int{},
			},
			want: nil,
		},
		{
			name: "empty ConcurrentRWArray and nil elements slice yields nil output",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				index:    2,
				elements: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.InsertInPlace(tt.args.index, tt.args.elements...)

			got := tt.a.GetAsSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertInPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleConcurrentRWArray_Length() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	fmt.Printf("%v\n", arr.Length())

	// Output:
	// 5
}

func TestConcurrentRWArray_Length(t *testing.T) {
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		want int
	}
	tests := []testCase[int]{
		{
			name: "counts 5 elements",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			want: 5,
		},
		{
			name: "empty ConcurrentRWArray has length 0",
			a:    lists.NewConcurrentRWArray[int](),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Length()
			if got != tt.want {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleConcurrentRWArray_PeekEnd() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	out, ok := arr.PeekEnd()

	fmt.Printf("%v, %t\n", out, ok)

	// Output:
	// 5, true
}

func TestConcurrentRWArray_PeekEnd(t *testing.T) {
	type testCase[T any] struct {
		name   string
		a      *lists.ConcurrentRWArray[T]
		want   T
		wantOK bool
	}
	tests := []testCase[int]{
		{
			name:   "peeks last element",
			a:      lists.NewConcurrentRWArray(1, 2, 3, 4, 5),
			want:   5,
			wantOK: true,
		},
		{
			name:   "empty ConcurrentRWArray yields nil output and false ok",
			a:      lists.NewConcurrentRWArray[int](),
			want:   0,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOK := tt.a.PeekEnd()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeekEnd() got = %v, want %v", got, tt.want)
			}
			if gotOK != tt.wantOK {
				t.Errorf("PeekEnd() got1 = %v, want %v", gotOK, tt.wantOK)
			}
		})
	}
}

func ExampleConcurrentRWArray_PeekFront() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	out, ok := arr.PeekFront()

	fmt.Printf("%v, %t\n", out, ok)

	// Output:
	// 1, true
}

func TestConcurrentRWArray_PeekFront(t *testing.T) {
	type testCase[T any] struct {
		name   string
		a      *lists.ConcurrentRWArray[T]
		want   T
		wantOK bool
	}
	tests := []testCase[int]{
		{
			name:   "peeks start of ConcurrentRWArray",
			a:      lists.NewConcurrentRWArray(1, 2, 3, 4, 5),
			want:   1,
			wantOK: true,
		},
		{
			name:   "empty ConcurrentRWArray yields nil output and false ok",
			a:      lists.NewConcurrentRWArray[int](),
			want:   0,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOK := tt.a.PeekFront()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeekFront() got = %v, want %v", got, tt.want)
			}
			if gotOK != tt.wantOK {
				t.Errorf("PeekFront() got1 = %v, want %v", gotOK, tt.wantOK)
			}
		})
	}
}

func ExampleConcurrentRWArray_Pop() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	out, ok, outSli := arr.Pop()

	fmt.Printf("%v, %v, %v\n", out, ok, outSli)

	// Output:
	// 5, true, [1 2 3 4]
}

func TestConcurrentRWArray_Pop(t *testing.T) {
	type testCase[T any] struct {
		name    string
		a       *lists.ConcurrentRWArray[T]
		want    T
		wantOK  bool
		wantSli []T
	}
	tests := []testCase[int]{
		{
			name:    "pops last element",
			a:       lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			want:    5,
			wantOK:  true,
			wantSli: []int{1, 2, 3, 4},
		},
		{
			name:    "empty ConcurrentRWArray yields zero value output and nil slice",
			a:       lists.NewConcurrentRWArray[int](),
			want:    0,
			wantOK:  false,
			wantSli: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOK, gotSli := tt.a.Pop()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pop() got = %v, want %v", got, tt.want)
			}
			if gotOK != tt.wantOK {
				t.Errorf("Pop() gotOK = %v, want %v", gotOK, tt.wantOK)
			}
			if !reflect.DeepEqual(gotSli, tt.wantSli) {
				t.Errorf("Pop() gotSli = %v, want %v", gotSli, tt.wantSli)
			}
		})
	}
}

func TestConcurrentRWArray_PopInPlace(t *testing.T) {
	type testCase[T any] struct {
		name     string
		a        *lists.ConcurrentRWArray[T]
		wantVal  T
		wantOK   bool
		wantRest []T
	}
	tests := []testCase[int]{
		{
			name:     "pops last element",
			a:        lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			wantVal:  5,
			wantOK:   true,
			wantRest: []int{1, 2, 3, 4},
		},
		{
			name:     "empty ConcurrentRWArray yields zero value output and nil slice",
			a:        lists.NewConcurrentRWArray[int](),
			wantVal:  0,
			wantOK:   false,
			wantRest: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOK := tt.a.PopInPlace()
			gotRest := tt.a.GetAsSlice()
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("PopInPlace() got = %v, want %v", gotVal, tt.wantVal)
			}
			if gotOK != tt.wantOK {
				t.Errorf("PopInPlace() gotOK = %v, want %v", gotOK, tt.wantOK)
			}
			if !reflect.DeepEqual(gotRest, tt.wantRest) {
				t.Errorf("PopInPlace() gotRest = %v, wantRest %v", gotRest, tt.wantRest)
			}
		})
	}
}

func ExampleConcurrentRWArray_Push() {
	arr := lists.NewConcurrentRWArray(1, 2, 3, 4, 5)

	out := arr.Push(10)

	fmt.Printf("%v\n", out)

	// Output:
	// [1 2 3 4 5 10]
}

func TestConcurrentRWArray_Push(t *testing.T) {
	type args[T any] struct {
		element T
	}
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "push adds element to end of ConcurrentRWArray",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				element: 10,
			},
			want: []int{1, 2, 3, 4, 5, 10},
		},
		{
			name: "pushing to empty ConcurrentRWArray adds element to end of ConcurrentRWArray",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				element: 10,
			},
			want: []int{10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Push(tt.args.element)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Push() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentRWArray_PushInPlace(t *testing.T) {
	type args[T any] struct {
		element T
	}
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "push adds element to end of ConcurrentRWArray",
			a:    lists.NewConcurrentRWArray[int](1, 2, 3, 4, 5),
			args: args[int]{
				element: 10,
			},
			want: []int{1, 2, 3, 4, 5, 10},
		},
		{
			name: "pushing to empty ConcurrentRWArray adds element to end of ConcurrentRWArray",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				element: 10,
			},
			want: []int{10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.PushInPlace(tt.args.element)

			got := tt.a.GetAsSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PushInPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleConcurrentRWArray_Sort() {
	arr := lists.NewConcurrentRWArray(5, 4, 3, 2, 1)

	out := arr.Sort(slices.AscendingSortFunc[int])

	fmt.Printf("%v\n", out)

	// Output:
	// [1 2 3 4 5]
}

func TestConcurrentRWArray_Sort(t *testing.T) {
	type args[T any] struct {
		fn func(T, T) bool
	}
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "sorts descending",
			a:    lists.NewConcurrentRWArray(1, 2, 3, 4, 5),
			args: args[int]{
				fn: slices.DescendingSortFunc[int],
			},
			want: []int{5, 4, 3, 2, 1},
		},
		{
			name: "sorts ascending",
			a:    lists.NewConcurrentRWArray(5, 4, 3, 2, 1),
			args: args[int]{
				fn: slices.AscendingSortFunc[int],
			},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "sorting empty ConcurrentRWArray results in nil",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				fn: slices.AscendingSortFunc[int],
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Sort(tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleConcurrentRWArray_SortInPlace() {
	arr := lists.NewConcurrentRWArray(5, 4, 3, 2, 1)

	arr.SortInPlace(slices.AscendingSortFunc[int])

	fmt.Printf("%v\n", arr.GetAsSlice())

	// Output:
	// [1 2 3 4 5]
}

func TestConcurrentRWArray_SortInPlace(t *testing.T) {
	type args[T any] struct {
		fn func(T, T) bool
	}
	type testCase[T any] struct {
		name string
		a    *lists.ConcurrentRWArray[T]
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "sorts descending",
			a:    lists.NewConcurrentRWArray(1, 2, 3, 4, 5),
			args: args[int]{
				fn: slices.DescendingSortFunc[int],
			},
			want: []int{5, 4, 3, 2, 1},
		},
		{
			name: "sorts ascending",
			a:    lists.NewConcurrentRWArray(5, 4, 3, 2, 1),
			args: args[int]{
				fn: slices.AscendingSortFunc[int],
			},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "sorting empty ConcurrentRWArray results in nil",
			a:    lists.NewConcurrentRWArray[int](),
			args: args[int]{
				fn: slices.AscendingSortFunc[int],
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.SortInPlace(tt.args.fn)

			got := tt.a.GetAsSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortInPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}
