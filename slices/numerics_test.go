package slices_test

import (
	"github.com/pickeringtech/go-collectionutil/constraints"
	"github.com/pickeringtech/go-collectionutil/slices"
	"testing"
)

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
