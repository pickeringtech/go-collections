package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

const floatTol = 1e-9

// floatsClose compares two float64s with NaN/Inf awareness: NaN equals only
// NaN, infinities must match exactly in sign, and finite values must agree
// within floatTol.
func floatsClose(a, b float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return math.IsNaN(a) && math.IsNaN(b)
	}
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return a == b
	}
	return math.Abs(a-b) <= floatTol
}

// linX and yLinear are perfectly linearly related (yLinear = 2·linX), giving a known
// covariance of 4 (population) / 5 (sample) and a correlation of exactly 1.
var (
	linX    = []float64{1, 2, 3, 4, 5}
	yLinear = []float64{2, 4, 6, 8, 10}
	yInvert = []float64{10, 8, 6, 4, 2}
)

func TestPopulationCovariance(t *testing.T) {
	tests := []struct {
		name string
		x, y []float64
		want float64
		ok   bool
	}{
		{name: "perfectly linear", x: linX, y: yLinear, want: 4, ok: true},
		{name: "perfectly inverse", x: linX, y: yInvert, want: -4, ok: true},
		{name: "single pair is defined as zero", x: []float64{3}, y: []float64{7}, want: 0, ok: true},
		{name: "length mismatch is undefined", x: []float64{1, 2, 3}, y: []float64{1, 2}, want: 0, ok: false},
		{name: "empty is undefined", x: []float64{}, y: []float64{}, want: 0, ok: false},
		{name: "nil is undefined", x: nil, y: nil, want: 0, ok: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := stats.PopulationCovariance(tc.x, tc.y)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if !floatsClose(got, tc.want) {
				t.Fatalf("covariance = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSampleCovariance(t *testing.T) {
	tests := []struct {
		name string
		x, y []float64
		want float64
		ok   bool
	}{
		{name: "perfectly linear", x: linX, y: yLinear, want: 5, ok: true},
		{name: "two pairs", x: []float64{1, 3}, y: []float64{2, 6}, want: 4, ok: true},
		{name: "single pair is undefined", x: []float64{3}, y: []float64{7}, want: 0, ok: false},
		{name: "length mismatch is undefined", x: []float64{1, 2}, y: []float64{1, 2, 3}, want: 0, ok: false},
		{name: "empty is undefined", x: nil, y: nil, want: 0, ok: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := stats.SampleCovariance(tc.x, tc.y)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if !floatsClose(got, tc.want) {
				t.Fatalf("covariance = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCorrelation(t *testing.T) {
	tests := []struct {
		name string
		x, y []float64
		want float64
		ok   bool
	}{
		{name: "perfect positive", x: linX, y: yLinear, want: 1, ok: true},
		{name: "perfect negative", x: linX, y: yInvert, want: -1, ok: true},
		{name: "shifted is still perfectly correlated", x: linX, y: []float64{5, 7, 9, 11, 13}, want: 1, ok: true},
		{name: "uncorrelated", x: []float64{1, 2, 3, 4}, y: []float64{1, 2, 2, 1}, want: 0, ok: true},
		{name: "constant x is undefined", x: []float64{2, 2, 2}, y: []float64{1, 2, 3}, want: 0, ok: false},
		{name: "constant y is undefined", x: []float64{1, 2, 3}, y: []float64{4, 4, 4}, want: 0, ok: false},
		{name: "single pair is undefined", x: []float64{1}, y: []float64{2}, want: 0, ok: false},
		{name: "length mismatch is undefined", x: []float64{1, 2}, y: []float64{1, 2, 3}, want: 0, ok: false},
		{name: "empty is undefined", x: nil, y: nil, want: 0, ok: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := stats.Correlation(tc.x, tc.y)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if !floatsClose(got, tc.want) {
				t.Fatalf("correlation = %v, want %v", got, tc.want)
			}
		})
	}
}

// Correlation must stay within [−1, 1] (it is covariance after normalization)
// and is symmetric: the correlation of (x, y) equals that of (y, x).
func TestCorrelationIsSymmetric(t *testing.T) {
	a := []float64{3, 1, 4, 1, 5, 9, 2, 6}
	b := []float64{2, 7, 1, 8, 2, 8, 1, 8}

	rab, okAB := stats.Correlation(a, b)
	rba, okBA := stats.Correlation(b, a)
	if !okAB || !okBA {
		t.Fatalf("ok = %v / %v, want both true", okAB, okBA)
	}
	if !floatsClose(rab, rba) {
		t.Fatalf("Correlation(a,b) = %v, Correlation(b,a) = %v; want equal", rab, rba)
	}
	if rab < -1-floatTol || rab > 1+floatTol {
		t.Fatalf("correlation %v out of [-1, 1]", rab)
	}
}

// The operations are generic over constraints.Numeric, so integer slices work
// without conversion and produce the same float64 results.
func TestRelationalWithIntegers(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}
	ys := []int{2, 4, 6, 8, 10}

	if cov, ok := stats.PopulationCovariance(xs, ys); !ok || !floatsClose(cov, 4) {
		t.Fatalf("PopulationCovariance = %v, %v; want 4, true", cov, ok)
	}
	if r, ok := stats.Correlation(xs, ys); !ok || !floatsClose(r, 1) {
		t.Fatalf("Correlation = %v, %v; want 1, true", r, ok)
	}
}

// A non-finite input propagates to a non-finite statistic rather than being
// silently dropped — the package's documented NaN/Inf policy.
func TestRelationalNaNPropagates(t *testing.T) {
	xs := []float64{1, 2, math.NaN(), 4}
	ys := []float64{2, 4, 6, 8}

	cov, ok := stats.PopulationCovariance(xs, ys)
	if !ok {
		t.Fatalf("ok = false, want true (NaN propagates with ok == true)")
	}
	if !math.IsNaN(cov) {
		t.Fatalf("covariance = %v, want NaN", cov)
	}
}
