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

func TestEuclideanDistance(t *testing.T) {
	got, ok := stats.EuclideanDistance([]float64{0, 0}, []float64{3, 4})
	if !ok || !approxEqual(got, 5) {
		t.Fatalf("EuclideanDistance = %v, %v; want 5, true", got, ok)
	}

	t.Run("identical points are distance zero", func(t *testing.T) {
		d, ok := stats.EuclideanDistance([]int{1, 2, 3}, []int{1, 2, 3})
		if !ok || d != 0 {
			t.Fatalf("EuclideanDistance = %v, %v; want 0, true", d, ok)
		}
	})

	t.Run("rejects empty and mismatched lengths", func(t *testing.T) {
		_, ok := stats.EuclideanDistance([]float64{}, []float64{})
		if ok {
			t.Errorf("empty reported ok")
		}
		_, ok = stats.EuclideanDistance([]float64{1}, []float64{1, 2})
		if ok {
			t.Errorf("mismatched lengths reported ok")
		}
	})

	t.Run("NaN propagates with ok", func(t *testing.T) {
		got, ok := stats.EuclideanDistance([]float64{math.NaN(), 0}, []float64{0, 0})
		if !ok {
			t.Fatalf("ok = false, want true (NaN propagates with ok == true)")
		}
		if !math.IsNaN(got) {
			t.Fatalf("EuclideanDistance = %v, want NaN", got)
		}
	})

	t.Run("Inf propagates with ok", func(t *testing.T) {
		got, ok := stats.EuclideanDistance([]float64{math.Inf(1), 0}, []float64{0, 0})
		if !ok {
			t.Fatalf("ok = false, want true (Inf propagates with ok == true)")
		}
		if !math.IsInf(got, 1) {
			t.Fatalf("EuclideanDistance = %v, want +Inf", got)
		}
	})

	t.Run("overflow-safe for huge coordinates", func(t *testing.T) {
		// A naive Σ(aᵢ−bᵢ)² squares 1e200 to 1e400, which overflows float64
		// to +Inf. The scaled sum-of-squares keeps the exact magnitude.
		got, ok := stats.EuclideanDistance([]float64{1e200, 0}, []float64{0, 0})
		if !ok {
			t.Fatalf("ok = false, want true")
		}
		if got != 1e200 {
			t.Fatalf("EuclideanDistance = %v, want 1e200 (no overflow)", got)
		}
	})

	t.Run("underflow-safe for tiny coordinates", func(t *testing.T) {
		// 3e-200 and 4e-200 form a scaled 3-4-5 triangle whose squares would
		// underflow to zero without scaling; the distance is exactly 5e-200.
		got, ok := stats.EuclideanDistance([]float64{0, 0}, []float64{3e-200, 4e-200})
		if !ok {
			t.Fatalf("ok = false, want true")
		}
		if !approxEqual(got/1e-200, 5) {
			t.Fatalf("EuclideanDistance = %v, want 5e-200 (no underflow)", got)
		}
	})
}

func TestCosineSimilarity(t *testing.T) {
	cases := map[string]struct {
		a, b []float64
		want float64
	}{
		"identical direction": {[]float64{1, 1}, []float64{2, 2}, 1},
		"orthogonal":          {[]float64{1, 0}, []float64{0, 1}, 0},
		"opposite":            {[]float64{1, 0}, []float64{-1, 0}, -1},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, ok := stats.CosineSimilarity(tc.a, tc.b)
			if !ok || !approxEqual(got, tc.want) {
				t.Fatalf("CosineSimilarity = %v, %v; want %v, true", got, ok, tc.want)
			}
		})
	}

	t.Run("rejects undefined inputs", func(t *testing.T) {
		_, ok := stats.CosineSimilarity([]float64{}, []float64{})
		if ok {
			t.Errorf("empty reported ok")
		}
		_, ok = stats.CosineSimilarity([]float64{1, 2}, []float64{1})
		if ok {
			t.Errorf("mismatched lengths reported ok")
		}
		_, ok = stats.CosineSimilarity([]float64{0, 0}, []float64{1, 2})
		if ok {
			t.Errorf("zero vector reported ok")
		}
	})

	t.Run("non-finite propagates", func(t *testing.T) {
		got, ok := stats.CosineSimilarity([]float64{1, math.Inf(1)}, []float64{1, 1})
		if !ok || (!math.IsNaN(got) && !math.IsInf(got, 0)) {
			t.Fatalf("CosineSimilarity = %v, %v; want non-finite, true", got, ok)
		}
	})
}
