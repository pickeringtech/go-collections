package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

// The canonical Wikipedia worked example: mean 5, Σ(x−mean)² = 32 over 8 values.
// Population variance = 32/8 = 4; sample variance = 32/7 ≈ 4.5714.
var canonical = []float64{2, 4, 4, 4, 5, 5, 7, 9}

func TestPopulationVariance(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
		want  float64
		ok    bool
	}{
		{name: "canonical dataset", input: canonical, want: 4, ok: true},
		{name: "single element is defined as zero", input: []float64{42}, want: 0, ok: true},
		{name: "two identical elements", input: []float64{7, 7}, want: 0, ok: true},
		{name: "empty is undefined", input: []float64{}, want: 0, ok: false},
		{name: "nil is undefined", input: nil, want: 0, ok: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := stats.PopulationVariance(tc.input)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if !approxEqual(got, tc.want) {
				t.Fatalf("variance = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSampleVariance(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
		want  float64
		ok    bool
	}{
		{name: "canonical dataset", input: canonical, want: 32.0 / 7.0, ok: true},
		{name: "two elements", input: []float64{1, 5}, want: 8, ok: true},
		{name: "single element is undefined", input: []float64{42}, want: 0, ok: false},
		{name: "empty is undefined", input: []float64{}, want: 0, ok: false},
		{name: "nil is undefined", input: nil, want: 0, ok: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := stats.SampleVariance(tc.input)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if !approxEqual(got, tc.want) {
				t.Fatalf("variance = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestPopulationStdDev(t *testing.T) {
	got, ok := stats.PopulationStdDev(canonical)
	if !ok || !approxEqual(got, 2) {
		t.Fatalf("PopulationStdDev = %v, %v; want 2, true", got, ok)
	}
	_, ok = stats.PopulationStdDev([]float64{})
	if ok {
		t.Fatalf("PopulationStdDev on empty: ok = true, want false")
	}
}

func TestSampleStdDev(t *testing.T) {
	got, ok := stats.SampleStdDev(canonical)
	if !ok || !approxEqual(got, math.Sqrt(32.0/7.0)) {
		t.Fatalf("SampleStdDev = %v, %v; want %v, true", got, ok, math.Sqrt(32.0/7.0))
	}
	_, ok = stats.SampleStdDev([]float64{42})
	if ok {
		t.Fatalf("SampleStdDev on single element: ok = true, want false")
	}
}

// TestNearConstantLargeMagnitude is the differentiator test called out in the
// issue. With a billion-scale offset the squares (~1e18) exceed float64's
// exact-integer range, so the naive Σx² − (Σx)²/n formula suffers catastrophic
// cancellation and reports garbage (often negative). Welford works only with
// small deltas from the running mean, so it recovers the exact variance of the
// underlying {4, 7, 13, 16} pattern: mean 10, Σ(x−mean)² = 90.
func TestNearConstantLargeMagnitude(t *testing.T) {
	const offset = 1e9
	input := []float64{offset + 4, offset + 7, offset + 13, offset + 16}

	popVar, ok := stats.PopulationVariance(input)
	if !ok || popVar != 22.5 { // 90 / 4 — exact, not approximate.
		t.Fatalf("PopulationVariance = %v, %v; want exactly 22.5, true", popVar, ok)
	}

	sampVar, ok := stats.SampleVariance(input)
	if !ok || sampVar != 30 { // 90 / 3 — exact.
		t.Fatalf("SampleVariance = %v, %v; want exactly 30, true", sampVar, ok)
	}

	// A naive sum-of-squares implementation, shown here purely to demonstrate it
	// disagrees with the correct answer on this very dataset.
	naive := naivePopulationVariance(input)
	if approxEqual(naive, 22.5) {
		t.Fatalf("naive formula unexpectedly accurate (%v); the dataset no longer guards against it", naive)
	}
}

// naivePopulationVariance is the lossy textbook formula the package deliberately
// avoids. It exists only so the test above can prove the chosen dataset really
// does expose the precision gap that Welford closes.
func naivePopulationVariance(input []float64) float64 {
	var sum, sumSq float64
	for _, x := range input {
		sum += x
		sumSq += x * x
	}
	n := float64(len(input))
	return sumSq/n - (sum/n)*(sum/n)
}

// TestNonFinitePropagates pins the variance family's policy: a non-finite input
// is surfaced, not silently dropped. The statistic comes back non-finite with
// ok == true, matching the documented behaviour of the covariance/correlation
// family it shares an algorithm with.
func TestNonFinitePropagates(t *testing.T) {
	fns := []struct {
		name string
		f    func([]float64) (float64, bool)
	}{
		{"PopulationVariance", stats.PopulationVariance[float64]},
		{"SampleVariance", stats.SampleVariance[float64]},
		{"PopulationStdDev", stats.PopulationStdDev[float64]},
		{"SampleStdDev", stats.SampleStdDev[float64]},
	}

	// A NaN poisons the running mean, so every statistic is NaN (ok stays true).
	for _, fn := range fns {
		got, ok := fn.f([]float64{1, 2, math.NaN(), 4})
		if !ok || !math.IsNaN(got) {
			t.Fatalf("%s(NaN input) = %v, %v; want NaN, true", fn.name, got, ok)
		}
	}

	// An ±Inf yields a non-finite (Inf or NaN) statistic, again with ok == true.
	for _, input := range [][]float64{{1, 2, math.Inf(1), 4}, {1, 2, math.Inf(-1), 4}} {
		for _, fn := range fns {
			got, ok := fn.f(input)
			if !ok || (!math.IsInf(got, 0) && !math.IsNaN(got)) {
				t.Fatalf("%s(%v) = %v, %v; want non-finite, true", fn.name, input, got, ok)
			}
		}
	}
}

// TestIntegerInput exercises the generic constraint over a non-float type and
// confirms the float64 conversion path.
func TestIntegerInput(t *testing.T) {
	got, ok := stats.PopulationVariance([]int{2, 4, 4, 4, 5, 5, 7, 9})
	if !ok || !approxEqual(got, 4) {
		t.Fatalf("PopulationVariance(ints) = %v, %v; want 4, true", got, ok)
	}
}
