package distance_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/ml/distance"
)

const floatTol = 1e-9

func floatsClose(a, b float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return math.IsNaN(a) && math.IsNaN(b)
	}
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return a == b
	}
	return math.Abs(a-b) <= floatTol
}

func TestEuclidean(t *testing.T) {
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
			name: "3-4-5 right triangle",
			args: args{a: []float64{0, 0}, b: []float64{3, 4}},
			want: 5,
			ok:   true,
		},
		{
			name: "identical vectors",
			args: args{a: []float64{1, 2, 3}, b: []float64{1, 2, 3}},
			want: 0,
			ok:   true,
		},
		{
			name: "single dimension",
			args: args{a: []float64{0}, b: []float64{7}},
			want: 7,
			ok:   true,
		},
		{
			name: "negative coordinates",
			args: args{a: []float64{-1, -1}, b: []float64{2, 3}},
			want: 5,
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
			got, ok := distance.Euclidean(tc.args.a, tc.args.b)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if ok && !floatsClose(got, tc.want) {
				t.Fatalf("Euclidean() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestEuclideanNaNPropagates(t *testing.T) {
	a := []float64{math.NaN(), 0}
	b := []float64{0, 0}
	got, ok := distance.Euclidean(a, b)
	if !ok {
		t.Fatalf("ok = false, want true (NaN propagates with ok == true)")
	}
	if !math.IsNaN(got) {
		t.Fatalf("Euclidean() = %v, want NaN", got)
	}
}

func TestEuclideanWithIntegers(t *testing.T) {
	a := []int{0, 0}
	b := []int{3, 4}
	got, ok := distance.Euclidean(a, b)
	if !ok {
		t.Fatalf("ok = false, want true")
	}
	if !floatsClose(got, 5) {
		t.Fatalf("Euclidean() = %v, want 5", got)
	}
}

func TestManhattan(t *testing.T) {
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
			name: "basic two dimensions",
			args: args{a: []float64{0, 0}, b: []float64{3, 4}},
			want: 7,
			ok:   true,
		},
		{
			name: "identical vectors",
			args: args{a: []float64{1, 2, 3}, b: []float64{1, 2, 3}},
			want: 0,
			ok:   true,
		},
		{
			name: "negative values",
			args: args{a: []float64{-1, -2}, b: []float64{1, 2}},
			want: 6,
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
			got, ok := distance.Manhattan(tc.args.a, tc.args.b)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if ok && !floatsClose(got, tc.want) {
				t.Fatalf("Manhattan() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestManhattanNaNPropagates(t *testing.T) {
	a := []float64{math.NaN(), 0}
	b := []float64{0, 0}
	got, ok := distance.Manhattan(a, b)
	if !ok {
		t.Fatalf("ok = false, want true (NaN propagates with ok == true)")
	}
	if !math.IsNaN(got) {
		t.Fatalf("Manhattan() = %v, want NaN", got)
	}
}

func TestMinkowski(t *testing.T) {
	type args struct {
		a, b []float64
		p    float64
	}
	tests := []struct {
		name string
		args args
		want float64
		ok   bool
	}{
		{
			name: "p=1 is manhattan",
			args: args{a: []float64{0, 0}, b: []float64{3, 4}, p: 1},
			want: 7,
			ok:   true,
		},
		{
			name: "p=2 is euclidean",
			args: args{a: []float64{0, 0}, b: []float64{3, 4}, p: 2},
			want: 5,
			ok:   true,
		},
		{
			name: "p less than 1 is undefined",
			args: args{a: []float64{0, 0}, b: []float64{1, 1}, p: 0.5},
			want: 0,
			ok:   false,
		},
		{
			name: "p=0 is undefined",
			args: args{a: []float64{0, 0}, b: []float64{1, 1}, p: 0},
			want: 0,
			ok:   false,
		},
		{
			name: "negative p is undefined",
			args: args{a: []float64{0, 0}, b: []float64{1, 1}, p: -1},
			want: 0,
			ok:   false,
		},
		{
			name: "NaN p is undefined",
			args: args{a: []float64{0, 0}, b: []float64{1, 1}, p: math.NaN()},
			want: 0,
			ok:   false,
		},
		{
			name: "+Inf p is undefined",
			args: args{a: []float64{0, 0}, b: []float64{1, 1}, p: math.Inf(1)},
			want: 0,
			ok:   false,
		},
		{
			name: "-Inf p is undefined",
			args: args{a: []float64{0, 0}, b: []float64{1, 1}, p: math.Inf(-1)},
			want: 0,
			ok:   false,
		},
		{
			name: "all-zero differences give zero",
			args: args{a: []float64{2, 3}, b: []float64{2, 3}, p: 3},
			want: 0,
			ok:   true,
		},
		{
			name: "length mismatch is undefined",
			args: args{a: []float64{1, 2, 3}, b: []float64{1, 2}, p: 2},
			want: 0,
			ok:   false,
		},
		{
			name: "empty is undefined",
			args: args{a: []float64{}, b: []float64{}, p: 2},
			want: 0,
			ok:   false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := distance.Minkowski(tc.args.a, tc.args.b, tc.args.p)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if ok && !floatsClose(got, tc.want) {
				t.Fatalf("Minkowski() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestMinkowskiNaNInputPropagates(t *testing.T) {
	a := []float64{math.NaN(), 0}
	b := []float64{0, 0}
	got, ok := distance.Minkowski(a, b, 3)
	if !ok {
		t.Fatalf("ok = false, want true (NaN propagates with ok == true)")
	}
	if !math.IsNaN(got) {
		t.Fatalf("Minkowski() = %v, want NaN", got)
	}
}

func TestMinkowskiInfInputPropagates(t *testing.T) {
	a := []float64{math.Inf(1), 0}
	b := []float64{0, 0}
	got, ok := distance.Minkowski(a, b, 3)
	if !ok {
		t.Fatalf("ok = false, want true (Inf propagates with ok == true)")
	}
	if !math.IsInf(got, 0) && !math.IsNaN(got) {
		t.Fatalf("Minkowski() = %v, want a non-finite result", got)
	}
}

// TestMinkowskiLargePStable verifies that large p does not overflow. A naive
// Σ diffᵖ implementation drives math.Pow(1e4, 100) to +Inf and then
// math.Pow(+Inf, 1/100) back to +Inf, losing the answer entirely. Factoring
// out the max term keeps the result finite and close to the dominating
// coordinate (the Chebyshev limit the metric approaches as p grows).
func TestMinkowskiLargePStable(t *testing.T) {
	a := []float64{0, 0}
	b := []float64{1e4, 5e3}
	got, ok := distance.Minkowski(a, b, 100)
	if !ok {
		t.Fatalf("ok = false, want true")
	}
	if math.IsInf(got, 0) || math.IsNaN(got) {
		t.Fatalf("Minkowski() = %v, want a finite result (no overflow)", got)
	}
	// With p=100 the smaller coordinate contributes (5e3/1e4)^100 ≈ 8e-31, so
	// the result is the larger coordinate to well within tolerance.
	if !floatsClose(got, 1e4) {
		t.Fatalf("Minkowski() = %v, want ~1e4", got)
	}
}

func TestCosineDistance(t *testing.T) {
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
			name: "identical direction",
			args: args{a: []float64{1, 2, 3}, b: []float64{2, 4, 6}},
			want: 0,
			ok:   true,
		},
		{
			name: "orthogonal vectors",
			args: args{a: []float64{1, 0}, b: []float64{0, 1}},
			want: 1,
			ok:   true,
		},
		{
			name: "anti-parallel vectors",
			args: args{a: []float64{1, 0}, b: []float64{-1, 0}},
			want: 2,
			ok:   true,
		},
		{
			name: "zero vector is undefined",
			args: args{a: []float64{0, 0}, b: []float64{1, 2}},
			want: 0,
			ok:   false,
		},
		{
			name: "length mismatch is undefined",
			args: args{a: []float64{1, 2}, b: []float64{1, 2, 3}},
			want: 0,
			ok:   false,
		},
		{
			name: "empty is undefined",
			args: args{a: []float64{}, b: []float64{}},
			want: 0,
			ok:   false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := distance.CosineDistance(tc.args.a, tc.args.b)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if ok && !floatsClose(got, tc.want) {
				t.Fatalf("CosineDistance() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHamming(t *testing.T) {
	type args struct {
		a, b []string
	}
	tests := []struct {
		name string
		args args
		want int
		ok   bool
	}{
		{
			name: "no differences",
			args: args{a: []string{"a", "b", "c"}, b: []string{"a", "b", "c"}},
			want: 0,
			ok:   true,
		},
		{
			name: "one difference",
			args: args{a: []string{"a", "b", "c"}, b: []string{"a", "x", "c"}},
			want: 1,
			ok:   true,
		},
		{
			name: "all different",
			args: args{a: []string{"a", "b", "c"}, b: []string{"x", "y", "z"}},
			want: 3,
			ok:   true,
		},
		{
			name: "both empty",
			args: args{a: []string{}, b: []string{}},
			want: 0,
			ok:   true,
		},
		{
			name: "both nil",
			args: args{a: nil, b: nil},
			want: 0,
			ok:   true,
		},
		{
			name: "length mismatch is undefined",
			args: args{a: []string{"a", "b"}, b: []string{"a", "b", "c"}},
			want: 0,
			ok:   false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := distance.Hamming(tc.args.a, tc.args.b)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if ok && got != tc.want {
				t.Fatalf("Hamming() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHammingWithInts(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := []int{1, 0, 3, 0}
	got, ok := distance.Hamming(a, b)
	if !ok {
		t.Fatalf("ok = false, want true")
	}
	if got != 2 {
		t.Fatalf("Hamming() = %v, want 2", got)
	}
}

func TestLevenshtein(t *testing.T) {
	type args struct {
		a, b string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "kitten to sitting",
			args: args{a: "kitten", b: "sitting"},
			want: 3,
		},
		{
			name: "identical strings",
			args: args{a: "hello", b: "hello"},
			want: 0,
		},
		{
			name: "empty to abc",
			args: args{a: "", b: "abc"},
			want: 3,
		},
		{
			name: "abc to empty",
			args: args{a: "abc", b: ""},
			want: 3,
		},
		{
			name: "both empty",
			args: args{a: "", b: ""},
			want: 0,
		},
		{
			name: "single insert",
			args: args{a: "cat", b: "cats"},
			want: 1,
		},
		{
			name: "single delete",
			args: args{a: "cats", b: "cat"},
			want: 1,
		},
		{
			name: "single substitute",
			args: args{a: "cat", b: "bat"},
			want: 1,
		},
		{
			name: "unicode runes counted as single characters",
			args: args{a: "café", b: "cafe"},
			want: 1, // accent substitution
		},
		{
			name: "symmetric",
			args: args{a: "saturday", b: "sunday"},
			want: 3,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := distance.Levenshtein(tc.args.a, tc.args.b)
			if got != tc.want {
				t.Fatalf("Levenshtein(%q, %q) = %v, want %v", tc.args.a, tc.args.b, got, tc.want)
			}
		})
	}
}

func TestLevenshteinIsSymmetric(t *testing.T) {
	a := "saturday"
	b := "sunday"
	ab := distance.Levenshtein(a, b)
	ba := distance.Levenshtein(b, a)
	if ab != ba {
		t.Fatalf("Levenshtein(%q, %q) = %d, Levenshtein(%q, %q) = %d; want equal", a, b, ab, b, a, ba)
	}
}

func TestLevenshteinAllThreeTransitions(t *testing.T) {
	// "ab" → "ba" exercises all three DP transitions (delete, insert, substitute)
	// in the edit-distance computation. Expected distance is 2.
	got := distance.Levenshtein("ab", "ba")
	if got != 2 {
		t.Fatalf("Levenshtein(\"ab\", \"ba\") = %d, want 2", got)
	}
}
