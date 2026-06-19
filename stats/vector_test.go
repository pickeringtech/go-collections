package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

func TestDot(t *testing.T) {
	type args struct {
		a, b []float64
	}
	tests := []struct {
		name string
		args args
		want float64
		ok   bool
	}{
		{
			name: "basic dot product",
			args: args{a: []float64{1, 2, 3}, b: []float64{4, 5, 6}},
			want: 32, // 1*4 + 2*5 + 3*6
			ok:   true,
		},
		{
			name: "orthogonal vectors",
			args: args{a: []float64{1, 0}, b: []float64{0, 1}},
			want: 0,
			ok:   true,
		},
		{
			name: "single element",
			args: args{a: []float64{5}, b: []float64{3}},
			want: 15,
			ok:   true,
		},
		{
			name: "negative values",
			args: args{a: []float64{-1, -2}, b: []float64{3, 4}},
			want: -11,
			ok:   true,
		},
		{
			name: "length mismatch is undefined",
			args: args{a: []float64{1, 2, 3}, b: []float64{1, 2}},
			want: 0,
			ok:   false,
		},
		{
			name: "empty is undefined",
			args: args{a: []float64{}, b: []float64{}},
			want: 0,
			ok:   false,
		},
		{
			name: "nil is undefined",
			args: args{a: nil, b: nil},
			want: 0,
			ok:   false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := stats.Dot(tc.args.a, tc.args.b)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if ok && !floatsClose(got, tc.want) {
				t.Fatalf("Dot() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDotNaNPropagates(t *testing.T) {
	a := []float64{1, math.NaN(), 3}
	b := []float64{4, 5, 6}
	got, ok := stats.Dot(a, b)
	if !ok {
		t.Fatalf("ok = false, want true (NaN propagates with ok == true)")
	}
	if !math.IsNaN(got) {
		t.Fatalf("Dot() = %v, want NaN", got)
	}
}

func TestDotWithIntegers(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	got, ok := stats.Dot(a, b)
	if !ok {
		t.Fatalf("ok = false, want true")
	}
	if !floatsClose(got, 32) {
		t.Fatalf("Dot() = %v, want 32", got)
	}
}

func TestNorm(t *testing.T) {
	type args struct {
		input []float64
	}
	tests := []struct {
		name string
		args args
		want float64
		ok   bool
	}{
		{
			name: "3-4-5 triple",
			args: args{input: []float64{3, 4}},
			want: 5,
			ok:   true,
		},
		{
			name: "unit vector",
			args: args{input: []float64{1, 0, 0}},
			want: 1,
			ok:   true,
		},
		{
			name: "single element",
			args: args{input: []float64{7}},
			want: 7,
			ok:   true,
		},
		{
			name: "all ones three elements",
			args: args{input: []float64{1, 1, 1}},
			want: math.Sqrt(3),
			ok:   true,
		},
		{
			name: "negative values",
			args: args{input: []float64{-3, 4}},
			want: 5,
			ok:   true,
		},
		{
			name: "empty is undefined",
			args: args{input: []float64{}},
			want: 0,
			ok:   false,
		},
		{
			name: "nil is undefined",
			args: args{input: nil},
			want: 0,
			ok:   false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := stats.Norm(tc.args.input)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if ok && !floatsClose(got, tc.want) {
				t.Fatalf("Norm() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestNormNaNPropagates(t *testing.T) {
	input := []float64{1, math.NaN(), 3}
	got, ok := stats.Norm(input)
	if !ok {
		t.Fatalf("ok = false, want true (NaN propagates with ok == true)")
	}
	if !math.IsNaN(got) {
		t.Fatalf("Norm() = %v, want NaN", got)
	}
}

func TestNormWithIntegers(t *testing.T) {
	input := []int{3, 4}
	got, ok := stats.Norm(input)
	if !ok {
		t.Fatalf("ok = false, want true")
	}
	if !floatsClose(got, 5) {
		t.Fatalf("Norm() = %v, want 5", got)
	}
}
