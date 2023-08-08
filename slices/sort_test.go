package slices_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/constraints"
	"github.com/pickeringtech/go-collections/slices"
	"reflect"
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
			name: "a == b == false",
			args: args[int]{
				a: 0,
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
			got := slices.AscendingSortFunc(tt.args.a, tt.args.b)
			if got != tt.want {
				t.Errorf("AscendingSortFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDescendingSortFunc(t *testing.T) {
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
			name: "a < b == false",
			args: args[int]{
				a: 0,
				b: 1,
			},
			want: false,
		},
		{
			name: "a > b == true",
			args: args[int]{
				a: 1,
				b: 0,
			},
			want: true,
		},
		{
			name: "a == b == true",
			args: args[int]{
				a: 0,
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
			want: false,
		},
		{
			name: "(a > b) < 0",
			args: args[int]{
				a: -1,
				b: -2,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.DescendingSortFunc(tt.args.a, tt.args.b)
			if got != tt.want {
				t.Errorf("DescendingSortFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleSort() {
	sli := []int{10, 2, -1, 1000, -10, 0, 1}

	sorted := slices.Sort(sli, slices.AscendingSortFunc[int])

	fmt.Printf("sorted: %v, original: %v", sorted, sli)
	// Output: sorted: [-10 -1 0 1 2 10 1000], original: [10 2 -1 1000 -10 0 1]
}

func TestSort(t *testing.T) {
	type args[T any] struct {
		input []T
		fun   slices.SortFunc[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "sorts numbers ascending",
			args: args[int]{
				input: []int{5, 2, 1, 3, 4, 9, 6, 8, 7},
				fun:   slices.AscendingSortFunc[int],
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "sorts numbers descending",
			args: args[int]{
				input: []int{5, 2, 1, 3, 4, 9, 6, 8, 7},
				fun:   slices.DescendingSortFunc[int],
			},
			want: []int{9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name: "handles nil input",
			args: args[int]{
				input: nil,
				fun:   slices.DescendingSortFunc[int],
			},
			want: nil,
		},
		{
			name: "handles empty input",
			args: args[int]{
				input: []int{},
				fun:   slices.DescendingSortFunc[int],
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.args.input
			orig := append(input[:0:0], input...)
			got := slices.Sort(input, tt.args.fun)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sort() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(orig, input) {
				t.Errorf("Sort() changed input - no changes expected")
			}
		})
	}
}

func BenchmarkSort(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
		fn   slices.SortFunc[int]
	}{
		{
			name: "3 elements desc",
			sli:  []int{1, 2, 3},
			fn:   slices.DescendingSortFunc[int],
		},
		{
			name: "10 elements desc",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn:   slices.DescendingSortFunc[int],
		},
		{
			name: "100 elements desc",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn:   slices.DescendingSortFunc[int],
		},
		{
			name: "1_000 elements desc",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn:   slices.DescendingSortFunc[int],
		},
		{
			name: "10_000 elements desc",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn:   slices.DescendingSortFunc[int],
		},
		{
			name: "100_000 elements desc",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn:   slices.DescendingSortFunc[int],
		},
		{
			name: "1_000_000 elements desc",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn:   slices.DescendingSortFunc[int],
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Sort(bm.sli, bm.fn)
			}
		})
	}
}

func ExampleSortInPlace() {
	sli := []int{10, 2, -1, 1000, -10, 0, 1}

	slices.SortInPlace(sli, slices.AscendingSortFunc[int])

	fmt.Printf("sorted: %v", sli)
	// Output: sorted: [-10 -1 0 1 2 10 1000]
}

func TestSortInPlace(t *testing.T) {
	type args[T any] struct {
		input []T
		fun   slices.SortFunc[T]
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "sorts numbers ascending",
			args: args[int]{
				input: []int{5, 2, 1, 3, 4, 9, 6, 8, 7},
				fun:   slices.AscendingSortFunc[int],
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "sorts numbers descending",
			args: args[int]{
				input: []int{5, 2, 1, 3, 4, 9, 6, 8, 7},
				fun:   slices.DescendingSortFunc[int],
			},
			want: []int{9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name: "handles nil input",
			args: args[int]{
				input: nil,
				fun:   slices.DescendingSortFunc[int],
			},
			want: nil,
		},
		{
			name: "handles empty input",
			args: args[int]{
				input: []int{},
				fun:   slices.DescendingSortFunc[int],
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.args.input
			slices.SortInPlace(tt.args.input, tt.args.fun)
			if !reflect.DeepEqual(input, tt.want) {
				t.Errorf("SortInPlace() = %v, want %v", input, tt.want)
			}
		})
	}
}

func BenchmarkSortInPlace(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
		fn   slices.SortFunc[int]
	}{
		{
			name: "3 elements desc",
			sli:  []int{1, 2, 3},
			fn:   slices.DescendingSortFunc[int],
		},
		{
			name: "10 elements desc",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
			fn:   slices.DescendingSortFunc[int],
		},
		{
			name: "100 elements desc",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
			fn:   slices.DescendingSortFunc[int],
		},
		{
			name: "1_000 elements desc",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			fn:   slices.DescendingSortFunc[int],
		},
		{
			name: "10_000 elements desc",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			fn:   slices.DescendingSortFunc[int],
		},
		{
			name: "100_000 elements desc",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			fn:   slices.DescendingSortFunc[int],
		},
		{
			name: "1_000_000 elements desc",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			fn:   slices.DescendingSortFunc[int],
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				slices.SortInPlace(bm.sli, bm.fn)
			}
		})
	}
}

func ExampleSortOrderedAsc() {
	sli := []int{10, 2, -1, 1000, -10, 0, 1}

	sorted := slices.SortOrderedAsc(sli)

	fmt.Printf("sorted: %v, original: %v", sorted, sli)
	// Output: sorted: [-10 -1 0 1 2 10 1000], original: [10 2 -1 1000 -10 0 1]
}

func TestSortOrderedAsc(t *testing.T) {
	type args[T constraints.Ordered] struct {
		input []T
	}
	type testCase[T constraints.Ordered] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "sorts numbers ascending",
			args: args[int]{
				input: []int{5, 2, 1, 3, 4, 9, 6, 8, 7},
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "handles nil input",
			args: args[int]{
				input: nil,
			},
			want: nil,
		},
		{
			name: "handles empty input",
			args: args[int]{
				input: []int{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.args.input
			orig := append(input[:0:0], input...)
			got := slices.SortOrderedAsc(input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortOrderedAsc() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(orig, input) {
				t.Errorf("SortOrderedAsc() changed input - no changes expected")
			}
		})
	}
}

func BenchmarkSortOrderedAsc(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
	}{
		{
			name: "3 elements desc",
			sli:  []int{1, 2, 3},
		},
		{
			name: "10 elements desc",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100 elements desc",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000 elements desc",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "10_000 elements desc",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100_000 elements desc",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000_000 elements desc",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.SortOrderedAsc(bm.sli)
			}
		})
	}
}

func ExampleSortOrderedAscInPlace() {
	sli := []int{10, 2, -1, 1000, -10, 0, 1}

	slices.SortOrderedAscInPlace(sli)

	fmt.Printf("sorted: %v", sli)
	// Output: sorted: [-10 -1 0 1 2 10 1000]
}

func TestSortOrderedAscInPlace(t *testing.T) {
	type args[T constraints.Ordered] struct {
		input []T
	}
	type testCase[T constraints.Ordered] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "sorts numbers ascending",
			args: args[int]{
				input: []int{5, 2, 1, 3, 4, 9, 6, 8, 7},
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "handles nil input",
			args: args[int]{
				input: nil,
			},
			want: nil,
		},
		{
			name: "handles empty input",
			args: args[int]{
				input: []int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slices.SortOrderedAscInPlace(tt.args.input)
			if !reflect.DeepEqual(tt.args.input, tt.want) {
				t.Errorf("SortAscInPlace() = %v, want %v", tt.args.input, tt.want)
			}
		})
	}
}

func BenchmarkSortOrderedAscInPlace(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
	}{
		{
			name: "3 elements desc",
			sli:  []int{1, 2, 3},
		},
		{
			name: "10 elements desc",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100 elements desc",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000 elements desc",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "10_000 elements desc",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100_000 elements desc",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000_000 elements desc",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				slices.SortOrderedAscInPlace(bm.sli)
			}
		})
	}
}

func ExampleSortOrderedDesc() {
	sli := []int{10, 2, -1, 1000, -10, 0, 1}

	sorted := slices.SortOrderedDesc(sli)

	fmt.Printf("sorted: %v, original: %v", sorted, sli)
	// Output: sorted: [1000 10 2 1 0 -1 -10], original: [10 2 -1 1000 -10 0 1]
}

func TestSortOrderedDesc(t *testing.T) {
	type args[T constraints.Ordered] struct {
		input []T
	}
	type testCase[T constraints.Ordered] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "sorts numbers descending",
			args: args[int]{
				input: []int{5, 2, 1, 3, 4, 9, 6, 8, 7},
			},
			want: []int{9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name: "handles nil input",
			args: args[int]{
				input: nil,
			},
			want: nil,
		},
		{
			name: "handles empty input",
			args: args[int]{
				input: []int{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.args.input
			orig := append(input[:0:0], input...)
			got := slices.SortOrderedDesc(input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortOrderedDesc() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(orig, input) {
				t.Errorf("SortOrderedDesc() changed input - no changes expected")
			}
		})
	}
}

func BenchmarkSortOrderedDesc(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
	}{
		{
			name: "3 elements desc",
			sli:  []int{1, 2, 3},
		},
		{
			name: "10 elements desc",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100 elements desc",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000 elements desc",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "10_000 elements desc",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100_000 elements desc",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000_000 elements desc",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.SortOrderedDesc(bm.sli)
			}
		})
	}
}

func ExampleSortOrderedDescInPlace() {
	sli := []int{10, 2, -1, 1000, -10, 0, 1}

	slices.SortOrderedDescInPlace(sli)

	fmt.Printf("sorted: %v", sli)
	// Output: sorted: [1000 10 2 1 0 -1 -10]
}

func TestSortOrderedDescInPlace(t *testing.T) {
	type args[T constraints.Ordered] struct {
		input []T
	}
	type testCase[T constraints.Ordered] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "sorts numbers descending",
			args: args[int]{
				input: []int{5, 2, 1, 3, 4, 9, 6, 8, 7},
			},
			want: []int{9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name: "handles nil input",
			args: args[int]{
				input: nil,
			},
			want: nil,
		},
		{
			name: "handles empty input",
			args: args[int]{
				input: []int{},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.args.input
			slices.SortOrderedDescInPlace(input)
			if !reflect.DeepEqual(input, tt.want) {
				t.Errorf("SortOrderedDescInPlace() = %v, want %v", input, tt.want)
			}
		})
	}
}

func BenchmarkSortOrderedDescInPlace(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  []int
	}{
		{
			name: "3 elements desc",
			sli:  []int{1, 2, 3},
		},
		{
			name: "10 elements desc",
			sli:  slices.Generate(10, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100 elements desc",
			sli:  slices.Generate(100, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000 elements desc",
			sli:  slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "10_000 elements desc",
			sli:  slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100_000 elements desc",
			sli:  slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000_000 elements desc",
			sli:  slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				slices.SortOrderedDescInPlace(bm.sli)
			}
		})
	}
}

func ExampleSortByOrderedField() {
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

	sorted := slices.SortByOrderedField(sli, slices.DescendingSortFunc[int], func(value numberAndSquare) int {
		return value.square
	})

	fmt.Printf("sorted: %v, original: %v", sorted, sli)
	// Output: sorted: [{9 81} {8 64} {7 49} {6 36} {5 25} {4 16} {3 9} {2 4} {1 1} {0 0}], original: [{0 0} {1 1} {2 4} {3 9} {4 16} {5 25} {6 36} {7 49} {8 64} {9 81}]
}

func TestSortByOrderedField(t *testing.T) {
	type language struct {
		name          string
		yearOfRelease int
	}
	type args[T any, S constraints.Ordered] struct {
		input     []T
		fun       slices.SortFunc[S]
		extractor slices.SortFieldExtractorFunc[T, S]
	}
	type testCase[T any, S constraints.Ordered] struct {
		name string
		args args[T, S]
		want []T
	}
	tests := []testCase[language, int]{
		{
			name: "sorts by year of release ascending as expected",
			args: args[language, int]{
				input: []language{
					{
						name:          "golang",
						yearOfRelease: 2009,
					},
					{
						name:          "c",
						yearOfRelease: 1972,
					},
					{
						name:          "rust",
						yearOfRelease: 2015,
					},
				},
				fun: slices.AscendingSortFunc[int],
				extractor: func(l language) int {
					return l.yearOfRelease
				},
			},
			want: []language{
				{
					name:          "c",
					yearOfRelease: 1972,
				},
				{
					name:          "golang",
					yearOfRelease: 2009,
				},
				{
					name:          "rust",
					yearOfRelease: 2015,
				},
			},
		},
		{
			name: "empty input provides nil output",
			args: args[language, int]{
				input: []language{},
				fun:   slices.AscendingSortFunc[int],
				extractor: func(l language) int {
					return l.yearOfRelease
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.SortByOrderedField[language, int](tt.args.input, tt.args.fun, tt.args.extractor)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortByOrderedField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSortByOrderedField(b *testing.B) {
	type person struct {
		age      int
		shoeSize int
	}

	generator := func(index int) person {
		return person{
			age:      index,
			shoeSize: index / 2,
		}
	}

	benchmarks := []struct {
		name string
		sli  []person
	}{
		{
			name: "10 elements desc",
			sli:  slices.Generate(10, generator),
		},
		{
			name: "100 elements desc",
			sli:  slices.Generate(100, generator),
		},
		{
			name: "1_000 elements desc",
			sli:  slices.Generate(1_000, generator),
		},
		{
			name: "10_000 elements desc",
			sli:  slices.Generate(10_000, generator),
		},
		{
			name: "100_000 elements desc",
			sli:  slices.Generate(100_000, generator),
		},
		{
			name: "1_000_000 elements desc",
			sli:  slices.Generate(1_000_000, generator),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.SortByOrderedField(bm.sli, slices.DescendingSortFunc[int], func(p person) int {
					return p.shoeSize
				})
			}
		})
	}
}
