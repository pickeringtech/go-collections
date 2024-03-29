package slices_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/slices"
	"reflect"
	"testing"
)

func ExampleAllMatch() {
	sli := []string{"hello", "glorious", "world"}

	allMediumLength := slices.AllMatch(sli, func(s string) bool {
		l := len(s)
		return l > 0 && l < 10
	})

	allShortLength := slices.AllMatch(sli, func(s string) bool {
		l := len(s)
		return l > 0 && l < 5
	})

	fmt.Printf("all medium: %v, all short: %v", allMediumLength, allShortLength)
	// Output: all medium: true, all short: false
}

func TestAllMatch(t *testing.T) {
	type args struct {
		input []int
		fun   slices.FindFunc[int]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "all elements pass test",
			args: args{
				input: []int{1, 2, 3, 4, 5},
				fun: func(element int) bool {
					return element < 6
				},
			},
			want: true,
		},
		{
			name: "some elements fail test",
			args: args{
				input: []int{1, 2, 3, 4, 5},
				fun: func(element int) bool {
					return element < 4
				},
			},
			want: false,
		},
		{
			name: "nil input results in true",
			args: args{
				input: nil,
				fun: func(element int) bool {
					return element < 4
				},
			},
			want: false,
		},
		{
			name: "empty input results in true",
			args: args{
				input: []int{},
				fun: func(element int) bool {
					return element < 4
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.AllMatch(tt.args.input, tt.args.fun)
			if got != tt.want {
				t.Errorf("AllMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkAllMatch(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
		fn   func(int) bool
	}{
		{
			name: "3 elements all match",
			sli:  []int{1, 2, 3},
			fn: func(i int) bool {
				return i < 4
			},
		},
		{
			name: "10 elements all match",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i < 10
			},
		},
		{
			name: "100 elements all match",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i < 100
			},
		},
		{
			name: "1_000 elements all match",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i < 1_000
			},
		},
		{
			name: "10_000 elements all match",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i < 10_000
			},
		},
		{
			name: "100_000 elements all match",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i < 100_000
			},
		},
		{
			name: "1_000_000 elements all match",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i < 1_000_000
			},
		},
		{
			name: "3 elements half match",
			sli:  []int{1, 2, 3},
			fn: func(i int) bool {
				return i < 2
			},
		},
		{
			name: "10 elements half match",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i < 5
			},
		},
		{
			name: "100 elements half match",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i < 50
			},
		},
		{
			name: "1_000 elements half match",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i < 500
			},
		},
		{
			name: "10_000 elements half match",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i < 5_000
			},
		},
		{
			name: "100_000 elements half match",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i < 50_000
			},
		},
		{
			name: "1_000_000 elements half match",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i < 500_000
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				slices.AllMatch(bm.sli, bm.fn)
			}
		})
	}
}

func ExampleAnyMatch() {
	sli := []int{1, 2, 3, 4, 5}

	hasEven := slices.AnyMatch(sli, func(i int) bool {
		return i%2 == 0
	})

	hasNegative := slices.AnyMatch(sli, func(i int) bool {
		return i < 0
	})

	fmt.Printf("has even: %v, has negative: %v", hasEven, hasNegative)
	// Output: has even: true, has negative: false
}

func TestAnyMatch(t *testing.T) {
	type args struct {
		input []int
		fun   slices.FindFunc[int]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "finds a match",
			args: args{
				input: []int{1, 3, 5, 7, 8},
				fun: func(element int) bool {
					return element%2 == 0
				},
			},
			want: true,
		},
		{
			name: "no match found",
			args: args{
				input: []int{1, 3, 5, 7},
				fun: func(element int) bool {
					return element%2 == 0
				},
			},
			want: false,
		},
		{
			name: "no match is found on nil input",
			args: args{
				input: nil,
				fun: func(element int) bool {
					return element%2 == 0
				},
			},
			want: false,
		},
		{
			name: "no match is found on empty input",
			args: args{
				input: []int{},
				fun: func(element int) bool {
					return element%2 == 0
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.AnyMatch(tt.args.input, tt.args.fun)
			if got != tt.want {
				t.Errorf("AnyMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkAnyMatch(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
		fn   func(int) bool
	}{
		{
			name: "3 elements with match",
			sli:  []int{1, 2, 3},
			fn: func(i int) bool {
				return i == 3
			},
		},
		{
			name: "10 elements with match",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 9
			},
		},
		{
			name: "100 elements with match",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 99
			},
		},
		{
			name: "1_000 elements with match",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 999
			},
		},
		{
			name: "10_000 elements with match",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 9_999
			},
		},
		{
			name: "100_000 elements with match",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 99_999
			},
		},
		{
			name: "1_000_000 elements with match",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 999_999
			},
		},
		{
			name: "3 elements no match",
			sli:  []int{1, 2, 3},
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "10 elements no match",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "100 elements no match",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "1_000 elements no match",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "10_000 elements no match",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "100_000 elements no match",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "1_000_000 elements no match",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				slices.AnyMatch(bm.sli, bm.fn)
			}
		})
	}
}

func ExampleFind() {
	sli := []int{1, 2, 3, 4, 5}

	firstEven, ok := slices.Find(sli, func(i int) bool {
		return i%2 == 0
	})

	fmt.Printf("first even: %v, ok: %v", firstEven, ok)
	// Output: first even: 2, ok: true
}

func TestFind(t *testing.T) {
	type args struct {
		input []int
		fun   slices.FindFunc[int]
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
		wantOk     bool
	}{
		{
			name: "selects the expected element",
			args: args{
				input: []int{2, 4, 6, 8, 10},
				fun: func(element int) bool {
					return element > 4
				},
			},
			wantResult: 6,
			wantOk:     true,
		},
		{
			name: "nil input results in zero value and boolean false",
			args: args{
				input: nil,
				fun: func(element int) bool {
					return element > 4
				},
			},
			wantResult: 0,
			wantOk:     false,
		},
		{
			name: "empty input results in zero value and boolean false",
			args: args{
				input: []int{},
				fun: func(element int) bool {
					return element > 4
				},
			},
			wantResult: 0,
			wantOk:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotOk := slices.Find(tt.args.input, tt.args.fun)
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("FindAny() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotOk != tt.wantOk {
				t.Errorf("FindAny() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func BenchmarkFind(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
		fn   func(i int) bool
	}{
		{
			name: "3 elements with match",
			sli:  []int{1, 2, 3},
			fn: func(i int) bool {
				return i == 3
			},
		},
		{
			name: "10 elements with match",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 9
			},
		},
		{
			name: "100 elements with match",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 99
			},
		},
		{
			name: "1_000 elements with match",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 999
			},
		},
		{
			name: "10_000 elements with match",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 9_999
			},
		},
		{
			name: "100_000 elements with match",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 99_999
			},
		},
		{
			name: "1_000_000 elements with match",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 999_999
			},
		},
		{
			name: "3 elements no match",
			sli:  []int{1, 2, 3},
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "10 elements no match",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "100 elements no match",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "1_000 elements no match",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "10_000 elements no match",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "100_000 elements no match",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "1_000_000 elements no match",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = slices.Find(bm.sli, bm.fn)
			}
		})
	}
}

func ExampleFindIndex() {
	sli := []int{1, 2, 3, 4, 5}

	firstEvenIdx := slices.FindIndex(sli, func(i int) bool {
		return i%2 == 0
	})

	missingIdx := slices.FindIndex(sli, func(i int) bool {
		return false
	})

	fmt.Printf("first even index: %v, not found index: %v", firstEvenIdx, missingIdx)
	// Output: first even index: 1, not found index: -1
}

func TestFindIndex(t *testing.T) {
	type args[T any] struct {
		input []T
		fun   slices.FindFunc[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want int
	}
	tests := []testCase[int]{
		{
			name: "finds expected element index",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
				fun: func(a int) bool {
					return a > 2
				},
			},
			want: 2,
		},
		{
			name: "no match results in -1",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
				fun: func(a int) bool {
					return a > 10
				},
			},
			want: -1,
		},
		{
			name: "nil input results in -1",
			args: args[int]{
				input: nil,
				fun: func(a int) bool {
					return a > 10
				},
			},
			want: -1,
		},
		{
			name: "empty input results in -1",
			args: args[int]{
				input: nil,
				fun: func(a int) bool {
					return a > 10
				},
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.FindIndex(tt.args.input, tt.args.fun)
			if got != tt.want {
				t.Errorf("FindIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkFindIndex(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
		fn   func(i int) bool
	}{
		{
			name: "3 elements with match",
			sli:  []int{1, 2, 3},
			fn: func(i int) bool {
				return i == 3
			},
		},
		{
			name: "10 elements with match",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 9
			},
		},
		{
			name: "100 elements with match",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 99
			},
		},
		{
			name: "1_000 elements with match",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 999
			},
		},
		{
			name: "10_000 elements with match",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 9_999
			},
		},
		{
			name: "100_000 elements with match",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 99_999
			},
		},
		{
			name: "1_000_000 elements with match",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 999_999
			},
		},
		{
			name: "3 elements no match",
			sli:  []int{1, 2, 3},
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "10 elements no match",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "100 elements no match",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "1_000 elements no match",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "10_000 elements no match",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "100_000 elements no match",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "1_000_000 elements no match",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.FindIndex(bm.sli, bm.fn)
			}
		})
	}
}

func ExampleFindLast() {
	sli := []int{1, 2, 3, 4, 5}

	lastEven, ok := slices.FindLast(sli, func(i int) bool {
		return i%2 == 0
	})

	fmt.Printf("last even: %v, ok: %v", lastEven, ok)
	// Output: last even: 4, ok: true
}

func TestFindLast(t *testing.T) {
	type args[T any] struct {
		input []T
		fun   slices.FindFunc[T]
	}
	type testCase[T any] struct {
		name       string
		args       args[T]
		wantResult T
		wantOk     bool
	}
	tests := []testCase[int]{
		{
			name: "finds the last entry matching the test function",
			args: args[int]{
				input: []int{5, 4, 3, 2, 1},
				fun: func(a int) bool {
					return a > 3
				},
			},
			wantResult: 4,
			wantOk:     true,
		},
		{
			name: "no match causes zero value and falsy boolean returns",
			args: args[int]{
				input: []int{5, 4, 3, 2, 1},
				fun: func(a int) bool {
					return a > 10
				},
			},
			wantResult: 0,
			wantOk:     false,
		},
		{
			name: "nil input causes zero value and falsy boolean returns",
			args: args[int]{
				input: nil,
				fun: func(a int) bool {
					return a > 10
				},
			},
			wantResult: 0,
			wantOk:     false,
		},
		{
			name: "empty input causes zero value and falsy boolean returns",
			args: args[int]{
				input: []int{},
				fun: func(a int) bool {
					return a > 10
				},
			},
			wantResult: 0,
			wantOk:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotOk := slices.FindLast(tt.args.input, tt.args.fun)
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("FindLast() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotOk != tt.wantOk {
				t.Errorf("FindLast() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func BenchmarkFindLast(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
		fn   func(i int) bool
	}{
		{
			name: "3 elements with match",
			sli:  []int{1, 2, 3},
			fn: func(i int) bool {
				return i == 1
			},
		},
		{
			name: "10 elements with match",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 0
			},
		},
		{
			name: "100 elements with match",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 0
			},
		},
		{
			name: "1_000 elements with match",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 0
			},
		},
		{
			name: "10_000 elements with match",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 0
			},
		},
		{
			name: "100_000 elements with match",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 0
			},
		},
		{
			name: "1_000_000 elements with match",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 0
			},
		},
		{
			name: "3 elements no match",
			sli:  []int{1, 2, 3},
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "10 elements no match",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "100 elements no match",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "1_000 elements no match",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "10_000 elements no match",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "100_000 elements no match",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "1_000_000 elements no match",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = slices.FindLast(bm.sli, bm.fn)
			}
		})
	}
}

func ExampleFindLastIndex() {
	sli := []int{1, 2, 3, 4, 5}

	lastEvenIdx := slices.FindLastIndex(sli, func(i int) bool {
		return i%2 == 0
	})

	missingIdx := slices.FindLastIndex(sli, func(i int) bool {
		return false
	})

	fmt.Printf("last even index: %v, not found index: %v", lastEvenIdx, missingIdx)
	// Output: last even index: 3, not found index: -1
}

func TestFindLastIndex(t *testing.T) {
	type args[T any] struct {
		input []T
		fun   slices.FindFunc[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want int
	}
	tests := []testCase[int]{
		{
			name: "finds the last entry matching the test function",
			args: args[int]{
				input: []int{5, 4, 3, 2, 1},
				fun: func(a int) bool {
					return a > 3
				},
			},
			want: 1,
		},
		{
			name: "no match causes zero value and falsy boolean returns",
			args: args[int]{
				input: []int{5, 4, 3, 2, 1},
				fun: func(a int) bool {
					return a > 10
				},
			},
			want: -1,
		},
		{
			name: "nil input causes zero value and falsy boolean returns",
			args: args[int]{
				input: nil,
				fun: func(a int) bool {
					return a > 10
				},
			},
			want: -1,
		},
		{
			name: "empty input causes zero value and falsy boolean returns",
			args: args[int]{
				input: []int{},
				fun: func(a int) bool {
					return a > 10
				},
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.FindLastIndex(tt.args.input, tt.args.fun)
			if got != tt.want {
				t.Errorf("FindLastIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkFindLastIndex(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
		fn   func(i int) bool
	}{
		{
			name: "3 elements with match",
			sli:  []int{1, 2, 3},
			fn: func(i int) bool {
				return i == 1
			},
		},
		{
			name: "10 elements with match",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 0
			},
		},
		{
			name: "100 elements with match",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 0
			},
		},
		{
			name: "1_000 elements with match",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 0
			},
		},
		{
			name: "10_000 elements with match",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 0
			},
		},
		{
			name: "100_000 elements with match",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 0
			},
		},
		{
			name: "1_000_000 elements with match",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return i == 0
			},
		},
		{
			name: "3 elements no match",
			sli:  []int{1, 2, 3},
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "10 elements no match",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "100 elements no match",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "1_000 elements no match",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "10_000 elements no match",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "100_000 elements no match",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
		{
			name: "1_000_000 elements no match",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn: func(i int) bool {
				return false
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.FindLastIndex(bm.sli, bm.fn)
			}
		})
	}
}

func ExampleFirst() {
	sli := []int{1, 2, 3, 4, 5}

	firstElement, ok := slices.First(sli)
	missingElement, notOk := slices.First([]int{})

	fmt.Printf("first element: %v, ok: %v, missing element: %v, missing ok: %v", firstElement, ok, missingElement, notOk)
	// Output: first element: 1, ok: true, missing element: 0, missing ok: false
}

func TestFirst(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
		wantOk     bool
	}{
		{
			name: "finds the first element",
			args: args{
				input: []int{1, 2, 3, 4},
			},
			wantResult: 1,
			wantOk:     true,
		},
		{
			name: "returns zero value and false if nil input",
			args: args{
				input: nil,
			},
			wantResult: 0,
			wantOk:     false,
		},
		{
			name: "returns zero value and false if empty input",
			args: args{
				input: []int{},
			},
			wantResult: 0,
			wantOk:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotOk := slices.First(tt.args.input)
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("FindFirst() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotOk != tt.wantOk {
				t.Errorf("FindFirst() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func BenchmarkFirst(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
	}{
		{
			name: "3 elements",
			sli:  []int{1, 2, 3},
		},
		{
			name: "10 elements",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100 elements",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000 elements",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "10_000 elements",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100_000 elements",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000_000 elements",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = slices.First(bm.sli)
			}
		})
	}
}

func ExampleGet() {
	sli := []int{1, 2, 3, 4, 5}

	firstElement := slices.Get(sli, 0, -1)
	missingElement := slices.Get(sli, 10, -1)

	fmt.Printf("first element: %v, missing element: %v", firstElement, missingElement)
	// Output: first element: 1, missing element: -1
}

func TestGet(t *testing.T) {
	type args[T any] struct {
		input        []T
		index        int
		defaultValue T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "retrieves the value at the specified index",
			args: args[int]{
				input:        []int{1, 2, 3, 4, 5},
				index:        2,
				defaultValue: -1,
			},
			want: 3,
		},
		{
			name: "receives the default value if the index is negative",
			args: args[int]{
				input:        []int{1, 2, 3, 4, 5},
				index:        -1,
				defaultValue: -1,
			},
			want: -1,
		},
		{
			name: "receives the default value if the index is equal to the length of the input",
			args: args[int]{
				input:        []int{1, 2, 3, 4, 5},
				index:        5,
				defaultValue: -1,
			},
			want: -1,
		},
		{
			name: "receives the default value if the input is empty",
			args: args[int]{
				input:        []int{},
				index:        5,
				defaultValue: -1,
			},
			want: -1,
		},
		{
			name: "receives the default value if the input is nil",
			args: args[int]{
				input:        nil,
				index:        5,
				defaultValue: -1,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.Get(tt.args.input, tt.args.index, tt.args.defaultValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	benchmarks := []struct {
		name         string
		sli          []int
		index        int
		defaultValue int
	}{
		{
			name:         "3 elements",
			sli:          []int{1, 2, 3},
			index:        1,
			defaultValue: -1,
		},
		{
			name:         "10 elements",
			sli:          slices.Generate(10, slices.NumericIdentityGenerator[int]),
			index:        5,
			defaultValue: -1,
		},
		{
			name:         "100 elements",
			sli:          slices.Generate(100, slices.NumericIdentityGenerator[int]),
			index:        50,
			defaultValue: -1,
		},
		{
			name:         "1_000 elements",
			sli:          slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			index:        500,
			defaultValue: -1,
		},
		{
			name:         "10_000 elements",
			sli:          slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			index:        5_000,
			defaultValue: -1,
		},
		{
			name:         "100_000 elements",
			sli:          slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			index:        50_000,
			defaultValue: -1,
		},
		{
			name:         "1_000_000 elements",
			sli:          slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			index:        500_000,
			defaultValue: -1,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Get(bm.sli, 0, bm.defaultValue)
			}
		})
	}
}

func ExampleIncludes() {
	sli := []int{1, 2, 3, 4, 5}

	isIncluded := slices.Includes(sli, 3)
	fmt.Printf("is included: %v", isIncluded)
	// Output: is included: true
}

func TestIncludes(t *testing.T) {
	type args[T comparable] struct {
		input []T
		value T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "finds included value",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
				value: 3,
			},
			want: true,
		},
		{
			name: "does not find if value is not included",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
				value: 6,
			},
			want: false,
		},
		{
			name: "nil input provides falsy return",
			args: args[int]{
				input: nil,
				value: 6,
			},
			want: false,
		},
		{
			name: "empty input provides falsy return",
			args: args[int]{
				input: []int{},
				value: 6,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Includes(tt.args.input, tt.args.value)
			if got != tt.want {
				t.Errorf("Includes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIncludes(b *testing.B) {
	benchmarks := []struct {
		name  string
		sli   []int
		value int
	}{
		{
			name:  "3 elements",
			sli:   []int{1, 2, 3},
			value: 2,
		},
		{
			name:  "10 elements",
			sli:   slices.Generate(10, slices.NumericIdentityGenerator[int]),
			value: 5,
		},
		{
			name:  "100 elements",
			sli:   slices.Generate(100, slices.NumericIdentityGenerator[int]),
			value: 50,
		},
		{
			name:  "1_000 elements",
			sli:   slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			value: 500,
		},
		{
			name:  "10_000 elements",
			sli:   slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			value: 5_000,
		},
		{
			name:  "100_000 elements",
			sli:   slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			value: 50_000,
		},
		{
			name:  "1_000_000 elements",
			sli:   slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			value: 500_000,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Includes(bm.sli, bm.value)
			}
		})
	}
}

func ExampleIndexOf() {
	sli := []int{1, 2, 3, 4, 5}

	idx := slices.IndexOf(sli, 3)

	fmt.Printf("index: %v", idx)
	// Output: index: 2
}

func TestIndexOf(t *testing.T) {
	type args[T comparable] struct {
		input []T
		value T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want int
	}
	tests := []testCase[int]{
		{
			name: "finds index of element in input",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
				value: 3,
			},
			want: 2,
		},
		{
			name: "not finding value results in -1",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
				value: 6,
			},
			want: -1,
		},
		{
			name: "nil input results in -1",
			args: args[int]{
				input: nil,
				value: 6,
			},
			want: -1,
		},
		{
			name: "empty input results in -1",
			args: args[int]{
				input: []int{},
				value: 6,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.IndexOf(tt.args.input, tt.args.value)
			if got != tt.want {
				t.Errorf("IndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIndexOf(b *testing.B) {
	benchmarks := []struct {
		name  string
		sli   []int
		value int
	}{
		{
			name:  "3 elements",
			sli:   []int{1, 2, 3},
			value: 2,
		},
		{
			name:  "10 elements",
			sli:   slices.Generate(10, slices.NumericIdentityGenerator[int]),
			value: 5,
		},
		{
			name:  "100 elements",
			sli:   slices.Generate(100, slices.NumericIdentityGenerator[int]),
			value: 50,
		},
		{
			name:  "1_000 elements",
			sli:   slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			value: 500,
		},
		{
			name:  "10_000 elements",
			sli:   slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			value: 5_000,
		},
		{
			name:  "100_000 elements",
			sli:   slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			value: 50_000,
		},
		{
			name:  "1_000_000 elements",
			sli:   slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			value: 500_000,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.IndexOf(bm.sli, bm.value)
			}
		})
	}
}

func ExampleIsEmpty() {
	sli := []int{1, 2, 3, 4, 5}

	isEmpty := slices.IsEmpty(sli)
	fmt.Printf("is empty: %v\n", isEmpty)

	sli = []int{}
	isEmpty = slices.IsEmpty(sli)
	fmt.Printf("is empty: %v\n", isEmpty)

	// Output:
	// is empty: false
	// is empty: true
}

func TestIsEmpty(t *testing.T) {
	type args[T any] struct {
		input []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "returns true if input is empty",
			args: args[int]{
				input: []int{},
			},
			want: true,
		},
		{
			name: "returns true if input is empty",
			args: args[int]{
				input: nil,
			},
			want: true,
		},
		{
			name: "returns false if input has a single element",
			args: args[int]{
				input: []int{1},
			},
			want: false,
		},
		{
			name: "returns false if input has multiple elements",
			args: args[int]{
				input: []int{1, 2},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.IsEmpty(tt.args.input); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIsEmpty(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
	}{
		{
			name: "empty sli",
			sli:  []int{},
		},
		{
			name: "nil sli",
			sli:  nil,
		},
		{
			name: "3 elements",
			sli:  []int{1, 2, 3},
		},
		{
			name: "10 elements",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100 elements",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000 elements",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "10_000 elements",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100_000 elements",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000_000 elements",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.IsEmpty(bm.sli)
			}
		})
	}
}

func ExampleLength() {
	sli := []int{1, 2, 3, 4, 5}

	length := slices.Length(sli)

	fmt.Printf("length: %v", length)
	// Output: length: 5
}

func TestLength(t *testing.T) {
	type args[T any] struct {
		input []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want int
	}
	tests := []testCase[int]{
		{
			name: "5 elements",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
			},
			want: 5,
		},
		{
			name: "4 elements",
			args: args[int]{
				input: []int{1, 2, 3, 4},
			},
			want: 4,
		},
		{
			name: "1 element",
			args: args[int]{
				input: []int{1},
			},
			want: 1,
		},
		{
			name: "empty input",
			args: args[int]{
				input: []int{},
			},
			want: 0,
		},
		{
			name: "nil input",
			args: args[int]{
				input: nil,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.Length(tt.args.input); got != tt.want {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkLength(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
	}{
		{
			name: "empty sli",
			sli:  []int{},
		},
		{
			name: "nil sli",
			sli:  nil,
		},
		{
			name: "3 elements",
			sli:  []int{1, 2, 3},
		},
		{
			name: "10 elements",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100 elements",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000 elements",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "10_000 elements",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100_000 elements",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000_000 elements",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Length(bm.sli)
			}
		})
	}
}

func ExamplePeekEnd() {
	sli := []int{1, 2, 3, 4, 5}

	lastElement, ok := slices.PeekEnd(sli)

	fmt.Printf("last element: %v, ok: %v", lastElement, ok)
	// Output: last element: 5, ok: true
}

func TestPeekEnd(t *testing.T) {
	type args[T any] struct {
		input []T
	}
	type testCase[T any] struct {
		name            string
		args            args[T]
		wantLastElement T
		wantOk          bool
	}
	tests := []testCase[int]{
		{
			name: "returns the last element",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
			},
			wantLastElement: 5,
			wantOk:          true,
		},
		{
			name: "empty input returns nothing and falsy boolean",
			args: args[int]{
				input: []int{},
			},
			wantLastElement: 0,
			wantOk:          false,
		},
		{
			name: "nil input returns nothing and falsy boolean",
			args: args[int]{
				input: nil,
			},
			wantLastElement: 0,
			wantOk:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLastElement, gotOk := slices.PeekEnd(tt.args.input)
			if !reflect.DeepEqual(gotLastElement, tt.wantLastElement) {
				t.Errorf("PeekEnd() gotLastElement = %v, want %v", gotLastElement, tt.wantLastElement)
			}
			if gotOk != tt.wantOk {
				t.Errorf("PeekEnd() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func BenchmarkPeekEnd(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
	}{
		{
			name: "empty input",
			sli:  []int{},
		},
		{
			name: "nil input",
			sli:  nil,
		},
		{
			name: "3 elements",
			sli:  []int{1, 2, 3},
		},
		{
			name: "10 elements",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100 elements",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000 elements",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "10_000 elements",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100_000 elements",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000_000 elements",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = slices.PeekEnd(bm.sli)
			}
		})
	}
}

func TestPeekFront(t *testing.T) {
	type args[T any] struct {
		input []T
	}
	type testCase[T any] struct {
		name            string
		args            args[T]
		wantLastElement T
		wantOk          bool
	}
	tests := []testCase[int]{
		{
			name: "returns the last element",
			args: args[int]{
				input: []int{1, 2, 3, 4, 5},
			},
			wantLastElement: 1,
			wantOk:          true,
		},
		{
			name: "empty input returns nothing and falsy boolean",
			args: args[int]{
				input: []int{},
			},
			wantLastElement: 0,
			wantOk:          false,
		},
		{
			name: "nil input returns nothing and falsy boolean",
			args: args[int]{
				input: nil,
			},
			wantLastElement: 0,
			wantOk:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLastElement, gotOk := slices.PeekFront(tt.args.input)
			if !reflect.DeepEqual(gotLastElement, tt.wantLastElement) {
				t.Errorf("PeekEnd() gotLastElement = %v, want %v", gotLastElement, tt.wantLastElement)
			}
			if gotOk != tt.wantOk {
				t.Errorf("PeekEnd() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func BenchmarkPeekFront(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
	}{
		{
			name: "empty input",
			sli:  []int{},
		},
		{
			name: "nil input",
			sli:  nil,
		},
		{
			name: "3 elements",
			sli:  []int{1, 2, 3},
		},
		{
			name: "10 elements",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100 elements",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000 elements",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "10_000 elements",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100_000 elements",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000_000 elements",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = slices.PeekFront(bm.sli)
			}
		})
	}
}

func ExampleSubSlice() {
	sli := []int{1, 2, 3, 4, 5}

	subSlice := slices.SubSlice(sli, 1, 3)

	fmt.Printf("sub-slice: %v, original: %v", subSlice, sli)
	// Output: sub-slice: [2 3], original: [1 2 3 4 5]
}

func TestSubSlice(t *testing.T) {
	type args[T any] struct {
		input     []T
		fromIndex int
		toIndex   int
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "creates sub-slice within valid range",
			args: args[int]{
				input:     []int{1, 2, 3, 4, 5},
				fromIndex: 1,
				toIndex:   4,
			},
			want: []int{2, 3, 4},
		},
		{
			name: "if range goes beyond end of input, length of result is reduced",
			args: args[int]{
				input:     []int{1, 2, 3, 4, 5},
				fromIndex: 3,
				toIndex:   6,
			},
			want: []int{4, 5},
		},
		{
			name: "if range goes is before start of input, length of result is reduced",
			args: args[int]{
				input:     []int{1, 2, 3, 4, 5},
				fromIndex: -1,
				toIndex:   2,
			},
			want: []int{1, 2},
		},
		{
			name: "if range is entirely before input, result is nil",
			args: args[int]{
				input:     []int{1, 2, 3, 4, 5},
				fromIndex: -2,
				toIndex:   -1,
			},
			want: nil,
		},
		{
			name: "if range is entirely after input, result is nil",
			args: args[int]{
				input:     []int{1, 2, 3, 4, 5},
				fromIndex: 6,
				toIndex:   7,
			},
			want: nil,
		},
		{
			name: "if fromIndex > toIndex, result is nil",
			args: args[int]{
				input:     []int{1, 2, 3, 4, 5},
				fromIndex: 3,
				toIndex:   2,
			},
			want: nil,
		},
		{
			name: "nil input produces nil output",
			args: args[int]{
				input:     nil,
				fromIndex: 0,
				toIndex:   1,
			},
			want: nil,
		},
		{
			name: "empty input produces nil output",
			args: args[int]{
				input:     []int{},
				fromIndex: 0,
				toIndex:   1,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.SubSlice(tt.args.input, tt.args.fromIndex, tt.args.toIndex)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SubSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSubSlice(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
		from int
		to   int
	}{
		{
			name: "3 elements",
			sli:  []int{1, 2, 3},
			from: 1,
			to:   2,
		},
		{
			name: "10 elements",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			from: 3,
			to:   8,
		},
		{
			name: "100 elements",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			from: 30,
			to:   80,
		},
		{
			name: "1_000 elements",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			from: 300,
			to:   800,
		},
		{
			name: "10_000 elements",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			from: 3_000,
			to:   8_000,
		},
		{
			name: "100_000 elements",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			from: 30_000,
			to:   80_000,
		},
		{
			name: "1_000_000 elements",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			from: 300_000,
			to:   800_000,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.SubSlice(bm.sli, bm.from, bm.to)
			}
		})
	}
}
