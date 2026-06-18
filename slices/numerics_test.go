package slices_test

import (
	"fmt"
	"github.com/pickeringtech/go-collections/constraints"
	"github.com/pickeringtech/go-collections/slices"
	"testing"
)

func ExampleNumericSlice_Avg() {
	sli := slices.NumericSlice[int]([]int{1, 2, 3, 4, 5})

	avg, ok := sli.Avg()
	fmt.Printf("average: %v, ok: %v, slice: %v", avg, ok, sli)
	// Output: average: 3, ok: true, slice: [1 2 3 4 5]
}

func TestNumericSlice_Avg(t *testing.T) {
	type testCase[T constraints.Numeric] struct {
		name   string
		n      slices.NumericSlice[T]
		want   float64
		wantOK bool
	}
	tests := []testCase[int]{
		{
			name:   "averages out correctly",
			n:      []int{1, 2, 3, 4, 5},
			want:   3,
			wantOK: true,
		},
		{
			name:   "empty input is undefined",
			n:      []int{},
			want:   0,
			wantOK: false,
		},
		{
			name:   "nil input is undefined",
			n:      nil,
			want:   0,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOK := tt.n.Avg()
			if got != tt.want || gotOK != tt.wantOK {
				t.Errorf("Avg() = (%v, %v), want (%v, %v)", got, gotOK, tt.want, tt.wantOK)
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
				_, _ = bm.sli.Avg()
			}
		})
	}
}

func ExampleNumericSlice_Max() {
	sli := slices.NumericSlice[int]([]int{1, 10, 1000, -10, -1, 0, 30})

	max, ok := sli.Max()
	fmt.Printf("max: %v, ok: %v, slice: %v", max, ok, sli)
	// Output: max: 1000, ok: true, slice: [1 10 1000 -10 -1 0 30]
}

func TestNumericSlice_Max(t *testing.T) {
	type testCase[T constraints.Numeric] struct {
		name   string
		n      slices.NumericSlice[T]
		want   T
		wantOK bool
	}
	tests := []testCase[int]{
		{
			name:   "selects the highest value",
			n:      []int{1, 10, 1000, -10, -1, 0, 340},
			want:   1000,
			wantOK: true,
		},
		{
			name:   "selects the highest value from all-negative input",
			n:      []int{-10, -3, -7},
			want:   -3,
			wantOK: true,
		},
		{
			name:   "empty input is undefined",
			n:      []int{},
			want:   0,
			wantOK: false,
		},
		{
			name:   "nil input is undefined",
			n:      nil,
			want:   0,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOK := tt.n.Max()
			if got != tt.want || gotOK != tt.wantOK {
				t.Errorf("Max() = (%v, %v), want (%v, %v)", got, gotOK, tt.want, tt.wantOK)
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
				_, _ = bm.sli.Max()
			}
		})
	}
}

func ExampleNumericSlice_Min() {
	sli := slices.NumericSlice[int]([]int{1, 10, 1000, -10, -1, 0, 30})

	min, ok := sli.Min()
	fmt.Printf("min: %v, ok: %v, slice: %v", min, ok, sli)
	// Output: min: -10, ok: true, slice: [1 10 1000 -10 -1 0 30]
}

func TestNumericSlice_Min(t *testing.T) {
	type testCase[T constraints.Numeric] struct {
		name   string
		n      slices.NumericSlice[T]
		want   T
		wantOK bool
	}
	tests := []testCase[int]{
		{
			name:   "selects the smallest value",
			n:      []int{1, 10, 1000, 340, -1, -100, 0, 20},
			want:   -100,
			wantOK: true,
		},
		{
			name:   "empty input is undefined",
			n:      []int{},
			want:   0,
			wantOK: false,
		},
		{
			name:   "nil input is undefined",
			n:      nil,
			want:   0,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOK := tt.n.Min()
			if got != tt.want || gotOK != tt.wantOK {
				t.Errorf("Min() = (%v, %v), want (%v, %v)", got, gotOK, tt.want, tt.wantOK)
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
				_, _ = bm.sli.Min()
			}
		})
	}
}

func ExampleNumericSlice_Sum() {
	sli := slices.NumericSlice[int]([]int{1, 2, 3, 4, 5})

	sum, ok := sli.Sum()
	fmt.Printf("sum: %v, ok: %v, slice: %v", sum, ok, sli)
	// Output: sum: 15, ok: true, slice: [1 2 3 4 5]
}

func TestNumericSlice_Sum(t *testing.T) {
	type testCase[T constraints.Numeric] struct {
		name   string
		n      slices.NumericSlice[T]
		want   T
		wantOK bool
	}
	tests := []testCase[int]{
		{
			name:   "calculates sum correctly, including negative numbers",
			n:      []int{1, 2, -1, 3, 4, 5},
			want:   14,
			wantOK: true,
		},
		{
			name:   "empty input is undefined",
			n:      []int{},
			want:   0,
			wantOK: false,
		},
		{
			name:   "nil input is undefined",
			n:      nil,
			want:   0,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOK := tt.n.Sum()
			if got != tt.want || gotOK != tt.wantOK {
				t.Errorf("Sum() = (%v, %v), want (%v, %v)", got, gotOK, tt.want, tt.wantOK)
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
				_, _ = bm.sli.Sum()
			}
		})
	}
}

func ExampleAvg() {
	sli := []int{1, 2, 3, 4, 5}

	avg, ok := slices.Avg(sli)

	fmt.Printf("avg: %v, ok: %v, slice: %v", avg, ok, sli)
	// Output: avg: 3, ok: true, slice: [1 2 3 4 5]
}

func TestAvg(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name   string
		args   args
		want   float64
		wantOK bool
	}{
		{
			name: "calculates expected average result",
			args: args{
				input: []int{1, 2, 3, 4, 5},
			},
			want:   3,
			wantOK: true,
		},
		{
			name: "all-zero input averages to zero and is defined",
			args: args{
				input: []int{0, 0, 0},
			},
			want:   0,
			wantOK: true,
		},
		{
			name: "nil input is undefined",
			args: args{
				input: nil,
			},
			want:   0,
			wantOK: false,
		},
		{
			name: "empty input is undefined",
			args: args{
				input: []int{},
			},
			want:   0,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOK := slices.Avg(tt.args.input)
			if got != tt.want || gotOK != tt.wantOK {
				t.Errorf("Avg() = (%v, %v), want (%v, %v)", got, gotOK, tt.want, tt.wantOK)
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
				_, _ = slices.Avg(bm.sli)
			}
		})
	}
}

func ExampleMax() {
	sli := []int{1, 10, 1000, -10, -1, 0, 30}

	max, ok := slices.Max(sli)
	fmt.Printf("max: %v, ok: %v, slice: %v", max, ok, sli)
	// Output: max: 1000, ok: true, slice: [1 10 1000 -10 -1 0 30]
}

func TestMax(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name   string
		args   args
		want   int
		wantOK bool
	}{
		{
			name: "finds the largest element in the input",
			args: args{
				input: []int{1, 2, 1, 1, 5, 0, 3, 4},
			},
			want:   5,
			wantOK: true,
		},
		{
			name: "finds the largest element in all-negative input",
			args: args{
				input: []int{-10, -3, -7},
			},
			want:   -3,
			wantOK: true,
		},
		{
			name: "nil input is undefined",
			args: args{
				input: nil,
			},
			want:   0,
			wantOK: false,
		},
		{
			name: "empty input is undefined",
			args: args{
				input: []int{},
			},
			want:   0,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOK := slices.Max(tt.args.input)
			if got != tt.want || gotOK != tt.wantOK {
				t.Errorf("Max() = (%v, %v), want (%v, %v)", got, gotOK, tt.want, tt.wantOK)
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
				_, _ = slices.Max(bm.sli)
			}
		})
	}
}

func ExampleMin() {
	sli := []int{1, 10, 1000, -10, -1, 0, 30}

	min, ok := slices.Min(sli)
	fmt.Printf("min: %v, ok: %v, slice: %v", min, ok, sli)
	// Output: min: -10, ok: true, slice: [1 10 1000 -10 -1 0 30]
}

func TestMin(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name   string
		args   args
		want   int
		wantOK bool
	}{
		{
			name: "finds the minimal value in the input",
			args: args{
				input: []int{1, 2, 1, 3, -3, 10},
			},
			want:   -3,
			wantOK: true,
		},
		{
			name: "nil input is undefined",
			args: args{
				input: nil,
			},
			want:   0,
			wantOK: false,
		},
		{
			name: "empty input is undefined",
			args: args{
				input: []int{},
			},
			want:   0,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOK := slices.Min(tt.args.input)
			if got != tt.want || gotOK != tt.wantOK {
				t.Errorf("Min() = (%v, %v), want (%v, %v)", got, gotOK, tt.want, tt.wantOK)
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
				_, _ = slices.Min(bm.sli)
			}
		})
	}
}

func ExampleSum() {
	sli := []int{1, 2, 3, 4, 5}

	sum, ok := slices.Sum(sli)
	fmt.Printf("sum: %v, ok: %v, slice: %v", sum, ok, sli)
	// Output: sum: 15, ok: true, slice: [1 2 3 4 5]
}

func TestSum(t *testing.T) {
	type args struct {
		input []int
	}
	tests := []struct {
		name   string
		args   args
		want   int
		wantOK bool
	}{
		{
			name: "results add up to expected amount",
			args: args{
				input: []int{1, 2, 3, 4, 5},
			},
			want:   15,
			wantOK: true,
		},
		{
			name: "nil input is undefined",
			args: args{
				input: nil,
			},
			want:   0,
			wantOK: false,
		},
		{
			name: "empty input is undefined",
			args: args{
				input: []int{},
			},
			want:   0,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOK := slices.Sum(tt.args.input)
			if got != tt.want || gotOK != tt.wantOK {
				t.Errorf("Sum() = (%v, %v), want (%v, %v)", got, gotOK, tt.want, tt.wantOK)
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
				_, _ = slices.Sum(bm.sli)
			}
		})
	}
}
