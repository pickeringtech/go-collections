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

// TestVectorOpsInfinityPropagates guards the compensated accumulator: an
// infinite running sum must stay infinite (not turn into NaN through the
// compensation term), while opposite infinities still cancel to NaN.
func TestVectorOpsInfinityPropagates(t *testing.T) {
	inf := math.Inf(1)
	tests := []struct {
		name string
		fn   func() (float64, bool)
		want float64
	}{
		{name: "Dot with a positive infinite product", fn: func() (float64, bool) { return stats.Dot([]float64{inf, 1}, []float64{1, 1}) }, want: inf},
		{name: "Dot with a negative infinite product", fn: func() (float64, bool) { return stats.Dot([]float64{-inf, 1, 2}, []float64{1, 1, 1}) }, want: -inf},
		{name: "Dot with opposite infinite products", fn: func() (float64, bool) { return stats.Dot([]float64{inf, -inf}, []float64{1, 1}) }, want: math.NaN()},
		{name: "Norm with an infinite element", fn: func() (float64, bool) { return stats.Norm([]float64{inf, 0}) }, want: inf},
		{name: "Norm with a negative infinite element", fn: func() (float64, bool) { return stats.Norm([]float64{0, -inf, 3}) }, want: inf},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.fn()
			if !ok || !floatsClose(got, tt.want) {
				t.Errorf("got %v, %v; want %v, true", got, ok, tt.want)
			}
		})
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

// TestNormAndCosineExtremeMagnitudes guards the sums of squares against
// overflow and underflow: squaring 1e200 overflows to +Inf and squaring 1e-200
// underflows to 0, so an unscaled Norm read +Inf and an unscaled cosine NaN or
// "undefined" for perfectly ordinary directions.
func TestNormAndCosineExtremeMagnitudes(t *testing.T) {
	tests := []struct {
		name string
		fn   func() (float64, bool)
		want float64
	}{
		{name: "Norm of a huge element", fn: func() (float64, bool) { return stats.Norm([]float64{1e200}) }, want: 1e200},
		{name: "Norm of a huge 3-4-5 triple", fn: func() (float64, bool) { return stats.Norm([]float64{3e200, -4e200}) }, want: 5e200},
		{name: "Norm of tiny elements", fn: func() (float64, bool) { return stats.Norm([]float64{3e-200, 4e-200}) }, want: 5e-200},
		{name: "Norm whose squares are subnormal", fn: func() (float64, bool) { return stats.Norm([]float64{3e-160, 4e-160}) }, want: 5e-160},
		{name: "Norm whose squares just stay normal", fn: func() (float64, bool) { return stats.Norm([]float64{3e-150, 4e-150}) }, want: 5e-150},
		{name: "Norm of the largest floats", fn: func() (float64, bool) { return stats.Norm([]float64{math.MaxFloat64, 0}) }, want: math.MaxFloat64},
		{name: "CosineSimilarity of huge vectors", fn: func() (float64, bool) {
			return stats.CosineSimilarity([]float64{1e200, 1e200}, []float64{1e200, 1e200})
		}, want: 1},
		{name: "CosineSimilarity of tiny vectors", fn: func() (float64, bool) {
			return stats.CosineSimilarity([]float64{1e-200, 1e-200}, []float64{1e-200, 1e-200})
		}, want: 1},
		{name: "CosineSimilarity of a huge and a tiny vector", fn: func() (float64, bool) { return stats.CosineSimilarity([]float64{1e200, 0}, []float64{0, 1e-200}) }, want: 0},
		{name: "CosineSimilarity of opposite huge vectors", fn: func() (float64, bool) {
			return stats.CosineSimilarity([]float64{1e300, -1e300}, []float64{-1e-300, 1e-300})
		}, want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.fn()
			tol := 1e-12 * math.Abs(tt.want)
			if tt.want == 0 {
				tol = 1e-12
			}
			if !ok || !(math.Abs(got-tt.want) <= tol) {
				t.Errorf("got %v, %v; want %v, true", got, ok, tt.want)
			}
		})
	}
}

// TestCosineSimilarityIsClamped pins the cosine to [−1, 1]: rounding in the
// dot product and norms can otherwise overshoot (1.0000000000000002), which
// would make a cosine distance negative.
func TestCosineSimilarityIsClamped(t *testing.T) {
	tests := []struct {
		name string
		a, b []float64
		want float64
	}{
		{name: "parallel", a: []float64{0.7000000000000001, 4.8999999999999995, 0.3}, b: []float64{0.7000000000000001, 4.8999999999999995, 0.3}, want: 1},
		{name: "anti-parallel", a: []float64{0.7000000000000001, 4.8999999999999995, 0.3}, b: []float64{-0.7000000000000001, -4.8999999999999995, -0.3}, want: -1},
		{name: "self of decimal vector", a: []float64{7.6, 6.1, 7.8}, b: []float64{7.6, 6.1, 7.8}, want: 1},
		{name: "self of another decimal vector", a: []float64{7.9, 2.3, 0.4}, b: []float64{7.9, 2.3, 0.4}, want: 1},
		{name: "scaled copy", a: []float64{1.8, 6.6, 4.4}, b: []float64{5.4, 19.799999999999997, 13.200000000000001}, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.CosineSimilarity(tt.a, tt.b)
			if !ok || got < -1 || got > 1 || !(math.Abs(got-tt.want) <= 1e-15) {
				t.Errorf("CosineSimilarity() = %v, %v; want %v within [-1, 1], true", got, ok, tt.want)
			}
		})
	}
}

// TestVectorOpsNaNWinsOverInfinity checks the non-finite precedence is
// independent of element order: any NaN makes the result NaN, as the IEEE
// expression √(Σ aᵢ²) would.
func TestVectorOpsNaNWinsOverInfinity(t *testing.T) {
	inf, nan := math.Inf(1), math.NaN()
	tests := []struct {
		name string
		fn   func() (float64, bool)
	}{
		{name: "Norm with Inf before NaN", fn: func() (float64, bool) { return stats.Norm([]float64{inf, nan}) }},
		{name: "Norm with NaN before Inf", fn: func() (float64, bool) { return stats.Norm([]float64{nan, inf}) }},
		{name: "CosineSimilarity with Inf", fn: func() (float64, bool) { return stats.CosineSimilarity([]float64{inf, 1}, []float64{1, 1}) }},
		{name: "CosineSimilarity with NaN against a zero vector", fn: func() (float64, bool) { return stats.CosineSimilarity([]float64{nan, 1}, []float64{0, 0}) }},
		{name: "CosineSimilarity with Inf against a zero vector", fn: func() (float64, bool) { return stats.CosineSimilarity([]float64{0, 0}, []float64{inf, 1}) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.fn()
			if !ok || !math.IsNaN(got) {
				t.Errorf("got %v, %v; want NaN, true", got, ok)
			}
		})
	}
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
