package slices

import "testing"

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
			if got := Sum(tt.args.input); got != tt.want {
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
			if got := Avg(tt.args.input); got != tt.want {
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
			if got := Max(tt.args.input); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}
