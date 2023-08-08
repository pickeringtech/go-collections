package slices_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/constraints"
	"github.com/pickeringtech/go-collections/slices"
	"testing"
)

func ExampleNumericSlice_Avg() {
	sli := slices.NumericSlice[int]([]int{1, 2, 3, 4, 5})

	avg := sli.Avg()
	fmt.Printf("average: %v, slice: %v", avg, sli)
	// Output: average: 3, slice: [1 2 3 4 5]
}

func TestNumericSlice_Avg(t *testing.T) {
	type testCase[T constraints.Numeric] struct {
		name string
		n    slices.NumericSlice[T]
		want float64
	}
	tests := []testCase[int]{
		{
			name: "averages out correctly",
			n:    []int{1, 2, 3, 4, 5},
			want: 3,
		},
		{
			name: "empty input provides zero output",
			n:    []int{},
			want: 0,
		},
		{
			name: "nil input provides zero output",
			n:    nil,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Avg(); got != tt.want {
				t.Errorf("Avg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkNumericSlice_Avg(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  slices.NumericSlice[int]
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
				_ = bm.sli.Avg()
			}
		})
	}
}

func ExampleNumericSlice_Max() {
	sli := slices.NumericSlice[int]([]int{1, 10, 1000, -10, -1, 0, 30})

	max := sli.Max()
	fmt.Printf("max: %v, slice: %v", max, sli)
	// Output: max: 1000, slice: [1 10 1000 -10 -1 0 30]
}

func TestNumericSlice_Max(t *testing.T) {
	type testCase[T constraints.Numeric] struct {
		name string
		n    slices.NumericSlice[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "selects the highest value",
			n:    []int{1, 10, 1000, -10, -1, 0, 340},
			want: 1000,
		},
		{
			name: "empty input results in zero output",
			n:    []int{},
			want: 0,
		},
		{
			name: "nil input results in zero output",
			n:    nil,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Max(); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkNumericSlice_Max(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  slices.NumericSlice[int]
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
				_ = bm.sli.Max()
			}
		})
	}
}

func ExampleNumericSlice_Min() {
	sli := slices.NumericSlice[int]([]int{1, 10, 1000, -10, -1, 0, 30})

	min := sli.Min()
	fmt.Printf("min: %v, slice: %v", min, sli)
	// Output: min: -10, slice: [1 10 1000 -10 -1 0 30]
}

func TestNumericSlice_Min(t *testing.T) {
	type testCase[T constraints.Numeric] struct {
		name string
		n    slices.NumericSlice[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "selects the smallest value",
			n:    []int{1, 10, 1000, 340, -1, -100, 0, 20},
			want: -100,
		},
		{
			name: "empty input results in zero output",
			n:    []int{},
			want: 0,
		},
		{
			name: "nil input results in zero output",
			n:    nil,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Min(); got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkNumericSlice_Min(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  slices.NumericSlice[int]
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
				_ = bm.sli.Min()
			}
		})
	}
}

func ExampleNumericSlice_Sum() {
	sli := slices.NumericSlice[int]([]int{1, 2, 3, 4, 5})

	sum := sli.Sum()
	fmt.Printf("sum: %v, slice: %v", sum, sli)
	// Output: sum: 15, slice: [1 2 3 4 5]
}

func TestNumericSlice_Sum(t *testing.T) {
	type testCase[T constraints.Numeric] struct {
		name string
		n    slices.NumericSlice[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "calculates sum correctly, including negative numbers",
			n:    []int{1, 2, -1, 3, 4, 5},
			want: 14,
		},
		{
			name: "empty input provides zero output",
			n:    []int{},
			want: 0,
		},
		{
			name: "nil input provides zero output",
			n:    nil,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.n.Sum()
			if got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkNumericSlice_Sum(b *testing.B) {
	benchmarks := []struct {
		name string
		sli  slices.NumericSlice[int]
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
				_ = bm.sli.Sum()
			}
		})
	}
}

func ExampleAvg() {
	sli := []int{1, 2, 3, 4, 5}

	avg := slices.Avg(sli)

	fmt.Printf("avg: %v, slice: %v", avg, sli)
	// Output: avg: 3, slice: [1 2 3 4 5]
}

func TestAvg(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "calculates expected average result",
			args: args{
				input: []int{1, 2, 3, 4, 5},
			},
			want: 3,
		},
		{
			name: "nil input results in zero",
			args: args{
				input: nil,
			},
			want: 0,
		},
		{
			name: "empty input results in zero",
			args: args{
				input: []int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Avg(tt.args.input)
			if got != tt.want {
				t.Errorf("Avg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkAvg(b *testing.B) {
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
				_ = slices.Avg(bm.sli)
			}
		})
	}
}

func ExampleMax() {
	sli := []int{1, 10, 1000, -10, -1, 0, 30}

	max := slices.Max(sli)
	fmt.Printf("max: %v, slice: %v", max, sli)
	// Output: max: 1000, slice: [1 10 1000 -10 -1 0 30]
}

func TestMax(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "finds the largest element in the input",
			args: args{
				input: []int{1, 2, 1, 1, 5, 0, 3, 4},
			},
			want: 5,
		},
		{
			name: "nil input provides zero",
			args: args{
				input: nil,
			},
			want: 0,
		},
		{
			name: "empty input provides zero",
			args: args{
				input: []int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Max(tt.args.input)
			if got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkMax(b *testing.B) {
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
				_ = slices.Max(bm.sli)
			}
		})
	}
}

func ExampleMin() {
	sli := []int{1, 10, 1000, -10, -1, 0, 30}

	min := slices.Min(sli)
	fmt.Printf("min: %v, slice: %v", min, sli)
	// Output: min: -10, slice: [1 10 1000 -10 -1 0 30]
}

func TestMin(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "finds the minimal value in the input",
			args: args{
				input: []int{1, 2, 1, 3, -3, 10},
			},
			want: -3,
		},
		{
			name: "nil input provides zero output",
			args: args{
				input: nil,
			},
			want: 0,
		},
		{
			name: "empty input provides zero output",
			args: args{
				input: []int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Min(tt.args.input)
			if got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkMin(b *testing.B) {
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
				_ = slices.Min(bm.sli)
			}
		})
	}
}

func ExampleSum() {
	sli := []int{1, 2, 3, 4, 5}

	sum := slices.Sum(sli)
	fmt.Printf("sum: %v, slice: %v", sum, sli)
	// Output: sum: 15, slice: [1 2 3 4 5]
}

func TestSum(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "results add up to expected amount",
			args: args{
				input: []int{1, 2, 3, 4, 5},
			},
			want: 15,
		},
		{
			name: "nil input results in zero",
			args: args{
				input: nil,
			},
			want: 0,
		},
		{
			name: "empty input results in zero",
			args: args{
				input: []int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Sum(tt.args.input)
			if got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSum(b *testing.B) {
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
				_ = slices.Sum(bm.sli)
			}
		})
	}
}
