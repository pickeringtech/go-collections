package slices_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/slices"
	"reflect"
	"testing"
)

func ExampleConcatenate() {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}

	joined := slices.Concatenate(a, b)
	fmt.Printf("%v", joined)
	// Output: [1 2 3 4 5 6]
}

func TestConcatenate(t *testing.T) {
	type args struct {
		inputA []int
		inputB []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "inputs are joined consecutively",
			args: args{
				inputA: []int{1, 2, 3},
				inputB: []int{4, 5, 6},
			},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "nil first input means result is second input",
			args: args{
				inputA: nil,
				inputB: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "empty first input means result is second input",
			args: args{
				inputA: []int{},
				inputB: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "nil second input means result is first input",
			args: args{
				inputA: []int{1, 2, 3},
				inputB: nil,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "empty second input means result is first input",
			args: args{
				inputA: []int{1, 2, 3},
				inputB: []int{},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Concatenate(tt.args.inputA, tt.args.inputB)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Concatenate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkConcatenate(b *testing.B) {
	benchmarks := []struct {
		name string
		sliA []int
		sliB []int
	}{
		{
			name: "3 elements",
			sliA: []int{1, 2, 3},
			sliB: []int{4, 5, 6},
		},
		{
			name: "10 elements",
			sliA: slices.Generate(10, slices.NumericIdentityGenerator[int]),
			sliB: slices.Generate(10, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100 elements",
			sliA: slices.Generate(100, slices.NumericIdentityGenerator[int]),
			sliB: slices.Generate(100, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000 elements",
			sliA: slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			sliB: slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "10_000 elements",
			sliA: slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			sliB: slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "100_000 elements",
			sliA: slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			sliB: slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name: "1_000_000 elements",
			sliA: slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			sliB: slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Concatenate(bm.sliA, bm.sliB)
			}
		})
	}
}

func ExampleCopy() {
	original := []int{1, 2, 3}

	sliCopy := slices.Copy(original)

	original = []int{4, 5, 6}

	fmt.Printf("original: %v, copy: %v", original, sliCopy)
	// Output: original: [4 5 6], copy: [1 2 3]
}

func TestCopy(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "duplicates the input into a new slice",
			args: args{
				input: []int{1, 2, 3},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "nil input provides nil output",
			args: args{
				input: nil,
			},
			want: nil,
		},
		{
			name: "empty input provides nil output",
			args: args{
				input: []int{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Copy(tt.args.input)
			tt.args.input = append(tt.args.input, 45)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCopy(b *testing.B) {
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
				_ = slices.Copy(bm.sli)
			}
		})
	}
}

func ExampleDelete() {
	sli := []int{1, 2, 3}

	withoutElement := slices.Delete(sli, 1)

	fmt.Printf("original: %v, with deleted element: %v", sli, withoutElement)
	// Output: original: [1 2 3], with deleted element: [1 3]
}

func TestDelete(t *testing.T) {
	type args struct {
		input []int
		index int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "removes the element at the specified index",
			args: args{
				input: []int{1, 2, 3, 4},
				index: 2,
			},
			want: []int{1, 2, 4},
		},
		{
			name: "removes the element at the last index",
			args: args{
				input: []int{1, 2, 3, 4},
				index: 3,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "removes the zeroth element",
			args: args{
				input: []int{1, 2, 3, 4},
				index: 0,
			},
			want: []int{2, 3, 4},
		},
		{
			name: "if index is beyond range the slice is not modified",
			args: args{
				input: []int{1, 2, 3, 4},
				index: 4,
			},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "if index is below zero the slice is not modified",
			args: args{
				input: []int{1, 2, 3, 4},
				index: -1,
			},
			want: []int{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origInput := slices.Copy(tt.args.input)
			got := slices.Delete(tt.args.input, tt.args.index)
			if !reflect.DeepEqual(origInput, tt.args.input) {
				t.Errorf("Delete() modified input slice - unexpected - original: %v, updated: %v", origInput, tt.args.input)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkDelete(b *testing.B) {
	benchmarks := []struct {
		name  string
		sli   []int
		index int
	}{
		{
			name:  "slice of 3 elements",
			sli:   []int{1, 2, 3},
			index: 1,
		},
		{
			name:  "10 elements",
			sli:   slices.Generate(10, slices.NumericIdentityGenerator[int]),
			index: 5,
		},
		{
			name:  "100 elements",
			sli:   slices.Generate(100, slices.NumericIdentityGenerator[int]),
			index: 50,
		},
		{
			name:  "1_000 elements",
			sli:   slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			index: 500,
		},
		{
			name:  "10_000 elements",
			sli:   slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			index: 5_000,
		},
		{
			name:  "100_000 elements",
			sli:   slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			index: 50_000,
		},
		{
			name:  "1_000_000 elements",
			sli:   slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			index: 500_000,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Delete(bm.sli, bm.index)
			}
		})
	}
}

func ExampleFill() {
	sli := []int{1, 2, 3}
	onlyZeroes := slices.Fill(sli, 0)

	fmt.Printf("original: %v, only zeroes: %v", sli, onlyZeroes)
	// Output: original: [1 2 3], only zeroes: [0 0 0]
}

func TestFill(t *testing.T) {
	type args[T any] struct {
		input []T
		value T
	}
	type testCase[T any] struct {
		name                   string
		args                   args[T]
		want                   []T
		ensureInputIsUnchanged bool
	}
	tests := []testCase[any]{
		{
			name: "fills every value",
			args: args[any]{
				input: []any{1, 2, 3, 4},
				value: 10,
			},
			want:                   []any{10, 10, 10, 10},
			ensureInputIsUnchanged: true,
		},
		{
			name: "nil input causes nil output",
			args: args[any]{
				input: nil,
				value: 10,
			},
			want:                   nil,
			ensureInputIsUnchanged: true,
		},
		{
			name: "empty input causes nil output",
			args: args[any]{
				input: []any{},
				value: 10,
			},
			want:                   nil,
			ensureInputIsUnchanged: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalInput := slices.Copy(tt.args.input)
			got := slices.Fill(tt.args.input, tt.args.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fill() = %v, want %v", got, tt.want)
			}
			if tt.ensureInputIsUnchanged && !reflect.DeepEqual(tt.args.input, originalInput) {
				t.Errorf("Fill() modified original input - original %v, modified input %v", originalInput, tt.args.input)
			}
		})
	}
}

func BenchmarkFill(b *testing.B) {
	benchmarks := []struct {
		name  string
		sli   []int
		value int
	}{
		{
			name:  "3 elements",
			sli:   []int{1, 2, 3},
			value: 999_999_999,
		},
		{
			name:  "10 elements",
			sli:   slices.Generate(10, slices.NumericIdentityGenerator[int]),
			value: 999_999_999,
		},
		{
			name:  "100 elements",
			sli:   slices.Generate(100, slices.NumericIdentityGenerator[int]),
			value: 999_999_999,
		},
		{
			name:  "1_000 elements",
			sli:   slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			value: 999_999_999,
		},
		{
			name:  "10_000 elements",
			sli:   slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			value: 999_999_999,
		},
		{
			name:  "100_000 elements",
			sli:   slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			value: 999_999_999,
		},
		{
			name:  "1_000_000 elements",
			sli:   slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			value: 999_999_999,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Fill(bm.sli, bm.value)
			}
		})
	}
}

func ExampleFillFrom() {
	sli := []int{1, 2, 3}
	filledFrom := slices.FillFrom(sli, 0, 1)

	fmt.Printf("original: %v, filled from element 1: %v", sli, filledFrom)
	// Output: original: [1 2 3], filled from element 1: [1 0 0]
}

func TestFillFrom(t *testing.T) {
	type args[T any] struct {
		input     []T
		value     T
		fromIndex int
	}
	type testCase[T any] struct {
		name                   string
		args                   args[T]
		want                   []T
		ensureInputIsUnchanged bool
	}
	tests := []testCase[any]{
		{
			name: "fills from a given index",
			args: args[any]{
				input:     []any{1, 2, 3, 4, 5},
				value:     10,
				fromIndex: 2,
			},
			want: []any{1, 2, 10, 10, 10},
		},
		{
			name: "if from index is beyond slice length it is unchanged",
			args: args[any]{
				input:     []any{1, 2, 3, 4, 5},
				value:     10,
				fromIndex: 5,
			},
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "if from index is below zero the slice is unchanged",
			args: args[any]{
				input:     []any{1, 2, 3, 4, 5},
				value:     10,
				fromIndex: -1,
			},
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "nil input results in nil output",
			args: args[any]{
				input:     nil,
				value:     10,
				fromIndex: -1,
			},
			want: nil,
		},
		{
			name: "empty input results in nil output",
			args: args[any]{
				input:     []any{},
				value:     10,
				fromIndex: -1,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalInput := slices.Copy(tt.args.input)
			got := slices.FillFrom(tt.args.input, tt.args.value, tt.args.fromIndex)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FillFrom() = %v, want %v", got, tt.want)
			}
			if tt.ensureInputIsUnchanged && !reflect.DeepEqual(tt.args.input, originalInput) {
				t.Errorf("FillFrom() modified original input - original %v, modified input %v", originalInput, tt.args.input)
			}
		})
	}
}

func BenchmarkFillFrom(b *testing.B) {
	benchmarks := []struct {
		name      string
		sli       []int
		value     int
		fromIndex int
	}{
		{
			name:      "3 elements",
			sli:       []int{1, 2, 3},
			value:     999_999_999,
			fromIndex: 1,
		},
		{
			name:      "10 elements",
			sli:       slices.Generate(10, slices.NumericIdentityGenerator[int]),
			value:     999_999_999,
			fromIndex: 5,
		},
		{
			name:      "100 elements",
			sli:       slices.Generate(100, slices.NumericIdentityGenerator[int]),
			value:     999_999_999,
			fromIndex: 50,
		},
		{
			name:      "1_000 elements",
			sli:       slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			value:     999_999_999,
			fromIndex: 500,
		},
		{
			name:      "10_000 elements",
			sli:       slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			value:     999_999_999,
			fromIndex: 5_000,
		},
		{
			name:      "100_000 elements",
			sli:       slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			value:     999_999_999,
			fromIndex: 50_000,
		},
		{
			name:      "1_000_000 elements",
			sli:       slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			value:     999_999_999,
			fromIndex: 500_000,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.FillFrom(bm.sli, bm.value, bm.fromIndex)
			}
		})
	}
}

func ExampleFillFromTo() {
	sli := []int{1, 2, 3, 4, 5}
	filledFrom := slices.FillFromTo(sli, 0, 1, 3)

	fmt.Printf("original: %v, filled in middle: %v", sli, filledFrom)
	// Output: original: [1 2 3 4 5], filled in middle: [1 0 0 4 5]
}

func TestFillFromTo(t *testing.T) {
	type args[T any] struct {
		input     []T
		value     T
		fromIndex int
		toIndex   int
	}
	type testCase[T any] struct {
		name                   string
		args                   args[T]
		want                   []T
		ensureInputIsUnchanged bool
	}
	tests := []testCase[any]{
		{
			name: "fills a range within a slice with a specified value",
			args: args[any]{
				input:     []any{1, 2, 3, 4, 5},
				value:     10,
				fromIndex: 2,
				toIndex:   4,
			},
			want: []any{1, 2, 10, 10, 5},
		},
		{
			name: "from index larger than to index causes no change",
			args: args[any]{
				input:     []any{1, 2, 3, 4, 5},
				value:     10,
				fromIndex: 4,
				toIndex:   2,
			},
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "negative from index causes no change",
			args: args[any]{
				input:     []any{1, 2, 3, 4, 5},
				value:     10,
				fromIndex: -1,
				toIndex:   2,
			},
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "to index beyond length of input causes no change",
			args: args[any]{
				input:     []any{1, 2, 3, 4, 5},
				value:     10,
				fromIndex: 0,
				toIndex:   6,
			},
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "nil input results in nil output",
			args: args[any]{
				input:     nil,
				value:     10,
				fromIndex: 0,
				toIndex:   6,
			},
			want: nil,
		},
		{
			name: "empty input results in nil output",
			args: args[any]{
				input:     []any{},
				value:     10,
				fromIndex: 0,
				toIndex:   6,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalInput := slices.Copy(tt.args.input)
			got := slices.FillFromTo(tt.args.input, tt.args.value, tt.args.fromIndex, tt.args.toIndex)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FillFromTo() = %v, want %v", got, tt.want)
			}
			if tt.ensureInputIsUnchanged && !reflect.DeepEqual(tt.args.input, originalInput) {
				t.Errorf("FillFromTo() modified original input - original %v, modified input %v", originalInput, tt.args.input)
			}
		})
	}
}

func BenchmarkFillFromTo(b *testing.B) {
	benchmarks := []struct {
		name      string
		sli       []int
		value     int
		fromIndex int
		toIndex   int
	}{
		{
			name:      "3 elements",
			sli:       []int{1, 2, 3},
			value:     999_999_999,
			fromIndex: 1,
			toIndex:   2,
		},
		{
			name:      "10 elements",
			sli:       slices.Generate(10, slices.NumericIdentityGenerator[int]),
			value:     999_999_999,
			fromIndex: 5,
			toIndex:   8,
		},
		{
			name:      "100 elements",
			sli:       slices.Generate(100, slices.NumericIdentityGenerator[int]),
			value:     999_999_999,
			fromIndex: 50,
			toIndex:   80,
		},
		{
			name:      "1_000 elements",
			sli:       slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			value:     999_999_999,
			fromIndex: 500,
			toIndex:   800,
		},
		{
			name:      "10_000 elements",
			sli:       slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			value:     999_999_999,
			fromIndex: 5_000,
			toIndex:   8_000,
		},
		{
			name:      "100_000 elements",
			sli:       slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			value:     999_999_999,
			fromIndex: 50_000,
			toIndex:   80_000,
		},
		{
			name:      "1_000_000 elements",
			sli:       slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			value:     999_999_999,
			fromIndex: 500_000,
			toIndex:   800_000,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.FillFromTo(bm.sli, bm.value, bm.fromIndex, bm.toIndex)
			}
		})
	}
}

func ExampleFillTo() {
	sli := []int{1, 2, 3, 4, 5}
	filledFrom := slices.FillTo(sli, 0, 3)

	fmt.Printf("original: %v, filled to middle: %v", sli, filledFrom)
	// Output: original: [1 2 3 4 5], filled to middle: [0 0 0 4 5]
}

func TestFillTo(t *testing.T) {
	type args[T any] struct {
		input   []T
		value   T
		toIndex int
	}
	type testCase[T any] struct {
		name                   string
		args                   args[T]
		want                   []T
		ensureInputIsUnchanged bool
	}
	tests := []testCase[any]{
		{
			name: "fills a range within a slice with a specified value",
			args: args[any]{
				input:   []any{1, 2, 3, 4, 5},
				value:   10,
				toIndex: 4,
			},
			want: []any{10, 10, 10, 10, 5},
		},
		{
			name: "negative to index causes no change",
			args: args[any]{
				input:   []any{1, 2, 3, 4, 5},
				value:   10,
				toIndex: -1,
			},
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "to index beyond length of input causes no change",
			args: args[any]{
				input:   []any{1, 2, 3, 4, 5},
				value:   10,
				toIndex: 6,
			},
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "nil input results in nil output",
			args: args[any]{
				input:   nil,
				value:   10,
				toIndex: 6,
			},
			want: nil,
		},
		{
			name: "empty input results in nil output",
			args: args[any]{
				input:   []any{},
				value:   10,
				toIndex: 6,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalInput := slices.Copy(tt.args.input)
			got := slices.FillTo(tt.args.input, tt.args.value, tt.args.toIndex)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FillTo() = %v, want %v", got, tt.want)
			}
			if tt.ensureInputIsUnchanged && !reflect.DeepEqual(tt.args.input, originalInput) {
				t.Errorf("FillFromTo() modified original input - original %v, modified input %v", originalInput, tt.args.input)
			}
		})
	}
}

func BenchmarkFillTo(b *testing.B) {
	benchmarks := []struct {
		name    string
		sli     []int
		value   int
		toIndex int
	}{
		{
			name:    "3 elements",
			sli:     []int{1, 2, 3},
			value:   999_999_999,
			toIndex: 1,
		},
		{
			name:    "10 elements",
			sli:     slices.Generate(10, slices.NumericIdentityGenerator[int]),
			value:   999_999_999,
			toIndex: 5,
		},
		{
			name:    "100 elements",
			sli:     slices.Generate(100, slices.NumericIdentityGenerator[int]),
			value:   999_999_999,
			toIndex: 50,
		},
		{
			name:    "1_000 elements",
			sli:     slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			value:   999_999_999,
			toIndex: 500,
		},
		{
			name:    "10_000 elements",
			sli:     slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			value:   999_999_999,
			toIndex: 5_000,
		},
		{
			name:    "100_000 elements",
			sli:     slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			value:   999_999_999,
			toIndex: 50_000,
		},
		{
			name:    "1_000_000 elements",
			sli:     slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			value:   999_999_999,
			toIndex: 500_000,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.FillTo([]int{1, 2, 3}, 0, 2)
			}
		})
	}
}

func ExampleInsert() {
	sli := []int{1, 2, 3, 4, 5}
	inserted := slices.Insert(sli, 2, 10, 11, 12)

	fmt.Printf("original: %v, inserted: %v", sli, inserted)
	// Output: original: [1 2 3 4 5], inserted: [1 2 10 11 12 3 4 5]
}

func TestInsert(t *testing.T) {
	type args[T any] struct {
		input    []T
		startIdx int
		elements []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "inserts one element",
			args: args[int]{
				input:    []int{1, 2, 3, 4, 5},
				startIdx: 2,
				elements: []int{10},
			},
			want: []int{1, 2, 10, 3, 4, 5},
		},
		{
			name: "inserts multiple elements",
			args: args[int]{
				input:    []int{1, 2, 3, 4, 5},
				startIdx: 2,
				elements: []int{10, 11, 12},
			},
			want: []int{1, 2, 10, 11, 12, 3, 4, 5},
		},
		{
			name: "empty input results in nil",
			args: args[int]{
				input:    []int{},
				startIdx: 0,
				elements: []int{10, 11, 12},
			},
			want: nil,
		},
		{
			name: "nil input results in nil",
			args: args[int]{
				input:    nil,
				startIdx: 0,
				elements: []int{10, 11, 12},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Insert(tt.args.input, tt.args.startIdx, tt.args.elements...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkInsert(b *testing.B) {
	benchmarks := []struct {
		name     string
		sli      []int
		startIdx int
		elements []int
	}{
		{
			name:     "3 elements",
			sli:      []int{1, 2, 3},
			startIdx: 1,
			elements: []int{10, 11, 12},
		},
		{
			name:     "10 elements",
			sli:      slices.Generate(10, slices.NumericIdentityGenerator[int]),
			startIdx: 5,
			elements: slices.Generate(10, slices.NumericIdentityGenerator[int]),
		},
		{
			name:     "100 elements",
			sli:      slices.Generate(100, slices.NumericIdentityGenerator[int]),
			startIdx: 50,
			elements: slices.Generate(100, slices.NumericIdentityGenerator[int]),
		},
		{
			name:     "1_000 elements",
			sli:      slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			startIdx: 500,
			elements: slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name:     "10_000 elements",
			sli:      slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			startIdx: 5_000,
			elements: slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name:     "100_000 elements",
			sli:      slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			startIdx: 50_000,
			elements: slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name:     "1_000_000 elements",
			sli:      slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			startIdx: 500_000,
			elements: slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Insert(bm.sli, bm.startIdx, bm.elements...)
			}
		})
	}
}

func ExampleJoinToString() {
	sli := []int{1, 2, 3}
	result := slices.JoinToString(sli, " + ")

	fmt.Printf("%s = 6", result)
	// Output: 1 + 2 + 3 = 6
}

func TestJoinToString(t *testing.T) {
	type args[T any] struct {
		input     []T
		separator string
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want string
	}
	tests := []testCase[any]{
		{
			name: "joins correctly using separator",
			args: args[any]{
				input:     []any{"Earth", "Wind", "Fire", "Water"},
				separator: "-and-",
			},
			want: "Earth-and-Wind-and-Fire-and-Water",
		},
		{
			name: "joins correctly using separator and varying types",
			args: args[any]{
				input:     []any{"Earth", 10, "Fire", []string{"a", "b", "c"}},
				separator: "-and-",
			},
			want: "Earth-and-10-and-Fire-and-[a b c]",
		},
		{
			name: "nil input results in empty output",
			args: args[any]{
				input:     nil,
				separator: "-and-",
			},
			want: "",
		},
		{
			name: "empty input results in empty output",
			args: args[any]{
				input:     []any{},
				separator: "-and-",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.JoinToString(tt.args.input, tt.args.separator)
			if got != tt.want {
				t.Errorf("JoinToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkJoinToString(b *testing.B) {
	benchmarks := []struct {
		name      string
		sli       []int
		separator string
	}{
		{
			name:      "3 elements",
			sli:       []int{1, 2, 3},
			separator: " + ",
		},
		{
			name:      "10 elements",
			sli:       slices.Generate(10, slices.NumericIdentityGenerator[int]),
			separator: " + ",
		},
		{
			name:      "100 elements",
			sli:       slices.Generate(100, slices.NumericIdentityGenerator[int]),
			separator: " + ",
		},
		{
			name:      "1_000 elements",
			sli:       slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			separator: " + ",
		},
		{
			name:      "10_000 elements",
			sli:       slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			separator: " + ",
		},
		{
			name:      "100_000 elements",
			sli:       slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			separator: " + ",
		},
		{
			name:      "1000_000 elements",
			sli:       slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			separator: " + ",
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.JoinToString(bm.sli, bm.separator)
			}
		})
	}
}

func ExamplePop() {
	sli := []int{1, 2, 3, 4, 5}

	lastElement, ok, shorterSli := slices.Pop(sli)
	fmt.Printf("last element: %v, ok: %v, shorter slice: %v, original slice: %v", lastElement, ok, shorterSli, sli)
	// Output: last element: 5, ok: true, shorter slice: [1 2 3 4], original slice: [1 2 3 4 5]
}

func TestPop(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantOK  bool
		wantSli []int
	}{
		{
			name: "gets the last element from the input returning the smaller slice",
			args: args{
				input: []int{1, 2, 3, 4},
			},
			want:    4,
			wantOK:  true,
			wantSli: []int{1, 2, 3},
		},
		{
			name: "popping last element from a slice with one element returns the element and a nil slice",
			args: args{
				input: []int{1},
			},
			want:    1,
			wantOK:  true,
			wantSli: nil,
		},
		{
			name: "nil input provides zero value for type and nil resulting slice",
			args: args{
				input: nil,
			},
			want:    0,
			wantOK:  false,
			wantSli: nil,
		},
		{
			name: "empty input provides zero value for type and nil resulting slice",
			args: args{
				input: []int{},
			},
			want:    0,
			wantOK:  false,
			wantSli: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOK, gotSli := slices.Pop(tt.args.input)
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

func BenchmarkPop(b *testing.B) {
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
				_, _, _ = slices.Pop(bm.sli)
			}
		})
	}
}

func ExamplePopFront() {
	sli := []int{1, 2, 3, 4, 5}

	firstElement, ok, shorterSli := slices.PopFront(sli)
	fmt.Printf("first element: %v, ok: %v, shorter slice: %v, original slice: %v", firstElement, ok, shorterSli, sli)
	// Output: first element: 1, ok: true, shorter slice: [2 3 4 5], original slice: [1 2 3 4 5]
}

func TestPopFront(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name             string
		args             args
		wantFirstElement int
		wantOK           bool
		wantNewSlice     []int
	}{
		{
			name: "first element is removed from input slice and returned",
			args: args{
				input: []int{5, 4, 3, 2, 1},
			},
			wantFirstElement: 5,
			wantOK:           true,
			wantNewSlice:     []int{4, 3, 2, 1},
		},
		{
			name: "popping last element from a slice with one element returns the element and a nil slice",
			args: args{
				input: []int{1},
			},
			wantFirstElement: 1,
			wantOK:           true,
			wantNewSlice:     nil,
		},
		{
			name: "nil input provides zero value output and nil resulting slice",
			args: args{
				input: nil,
			},
			wantFirstElement: 0,
			wantOK:           false,
			wantNewSlice:     nil,
		},
		{
			name: "empty input provides zero value output and nil resulting slice",
			args: args{
				input: []int{},
			},
			wantFirstElement: 0,
			wantOK:           false,
			wantNewSlice:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFirstElement, gotOK, gotNewSlice := slices.PopFront(tt.args.input)
			if !reflect.DeepEqual(gotFirstElement, tt.wantFirstElement) {
				t.Errorf("PopFront() gotFirstElement = %v, want %v", gotFirstElement, tt.wantFirstElement)
			}
			if gotOK != tt.wantOK {
				t.Errorf("PopFront() gotOK = %v, want %v", gotOK, tt.wantOK)
			}
			if !reflect.DeepEqual(gotNewSlice, tt.wantNewSlice) {
				t.Errorf("PopFront() gotNewSlice = %v, want %v", gotNewSlice, tt.wantNewSlice)
			}
		})
	}
}

func BenchmarkPopFront(b *testing.B) {
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
				_, _, _ = slices.PopFront(bm.sli)
			}
		})
	}
}

func ExamplePush() {
	sli := []int{1, 2, 3, 4}

	longerSli := slices.Push(sli, 5, 6)
	fmt.Printf("longer slice: %v, original slice: %v", longerSli, sli)
	// Output: longer slice: [1 2 3 4 5 6], original slice: [1 2 3 4]
}

func TestPush(t *testing.T) {
	type args struct {
		input       []int
		newElements []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "pushes the new elements to the end of the input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: []int{4, 5, 6},
			},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "nil input slice results in only the new elements",
			args: args{
				input:       nil,
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "empty input slice results in only the new elements",
			args: args{
				input:       []int{},
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "nil new elements results in only original input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: nil,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "empty new elements results in only original input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: []int{},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Push(tt.args.input, tt.args.newElements...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Push() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkPush(b *testing.B) {
	benchmarks := []struct {
		name   string
		sli    []int
		values []int
	}{
		{
			name:   "3 elements",
			sli:    []int{1, 2, 3},
			values: []int{4, 5, 6},
		},
		{
			name:   "10 elements",
			sli:    slices.Generate(10, slices.NumericIdentityGenerator[int]),
			values: slices.Generate(10, slices.NumericIdentityGenerator[int]),
		},
		{
			name:   "100 elements",
			sli:    slices.Generate(100, slices.NumericIdentityGenerator[int]),
			values: slices.Generate(100, slices.NumericIdentityGenerator[int]),
		},
		{
			name:   "1_000 elements",
			sli:    slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			values: slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name:   "10_000 elements",
			sli:    slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			values: slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name:   "100_000 elements",
			sli:    slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			values: slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name:   "1_000_000 elements",
			sli:    slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			values: slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.Push(bm.sli, bm.values...)
			}
		})
	}
}

func ExamplePushFront() {
	sli := []int{1, 2, 3, 4}

	longerSli := slices.PushFront(sli, -1, 0)
	fmt.Printf("longer slice: %v, original slice: %v", longerSli, sli)
	// Output: longer slice: [-1 0 1 2 3 4], original slice: [1 2 3 4]
}

func TestPushFront(t *testing.T) {
	type args struct {
		input       []int
		newElements []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "adds the new elements to the front of the input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6, 1, 2, 3},
		},
		{
			name: "nil new elements results in original input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: nil,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "empty new elements results in original input slice",
			args: args{
				input:       []int{1, 2, 3},
				newElements: []int{},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "nil input slice results in only new elements",
			args: args{
				input:       nil,
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
		{
			name: "empty input slice results in only new elements",
			args: args{
				input:       []int{},
				newElements: []int{4, 5, 6},
			},
			want: []int{4, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.PushFront(tt.args.input, tt.args.newElements...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PushFront() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkPushFront(b *testing.B) {
	benchmarks := []struct {
		name   string
		sli    []int
		values []int
	}{
		{
			name:   "3 elements",
			sli:    []int{1, 2, 3},
			values: []int{4, 5, 6},
		},
		{
			name:   "10 elements",
			sli:    slices.Generate(10, slices.NumericIdentityGenerator[int]),
			values: slices.Generate(10, slices.NumericIdentityGenerator[int]),
		},
		{
			name:   "100 elements",
			sli:    slices.Generate(100, slices.NumericIdentityGenerator[int]),
			values: slices.Generate(100, slices.NumericIdentityGenerator[int]),
		},
		{
			name:   "1_000 elements",
			sli:    slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
			values: slices.Generate(1_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name:   "10_000 elements",
			sli:    slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
			values: slices.Generate(10_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name:   "100_000 elements",
			sli:    slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
			values: slices.Generate(100_000, slices.NumericIdentityGenerator[int]),
		},
		{
			name:   "1_000_000 elements",
			sli:    slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
			values: slices.Generate(1_000_000, slices.NumericIdentityGenerator[int]),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slices.PushFront(bm.sli, bm.values...)
			}
		})
	}
}
