package stats_test

import (
	"math"
	"testing"

	"github.com/pickeringtech/go-collections/stats"
)

const epsilon = 1e-9

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) <= epsilon
}

func TestWeightedMean(t *testing.T) {
	tests := []struct {
		name    string
		values  []float64
		weights []float64
		want    float64
		wantOK  bool
	}{
		{
			name:    "weighted average of prices",
			values:  []float64{10, 20, 30},
			weights: []float64{1, 2, 3},
			want:    140.0 / 6.0,
			wantOK:  true,
		},
		{
			name:    "equal weights match arithmetic mean",
			values:  []float64{2, 4, 6, 8},
			weights: []float64{1, 1, 1, 1},
			want:    5,
			wantOK:  true,
		},
		{
			name:    "weights need not sum to one",
			values:  []float64{1, 2, 3},
			weights: []float64{10, 10, 10},
			want:    2,
			wantOK:  true,
		},
		{
			name:    "single value",
			values:  []float64{42},
			weights: []float64{7},
			want:    42,
			wantOK:  true,
		},
		{
			name:    "zero weight ignores its value",
			values:  []float64{1, 1000},
			weights: []float64{1, 0},
			want:    1,
			wantOK:  true,
		},
		{
			name:    "empty input",
			values:  []float64{},
			weights: []float64{},
			wantOK:  false,
		},
		{
			name:    "nil input",
			values:  nil,
			weights: nil,
			wantOK:  false,
		},
		{
			name:    "mismatched lengths",
			values:  []float64{1, 2, 3},
			weights: []float64{1, 2},
			wantOK:  false,
		},
		{
			name:    "negative weight rejected",
			values:  []float64{1, 2, 3},
			weights: []float64{1, -1, 1},
			wantOK:  false,
		},
		{
			name:    "zero-sum weights rejected",
			values:  []float64{1, 2, 3},
			weights: []float64{0, 0, 0},
			wantOK:  false,
		},
		{
			name:    "NaN value rejected",
			values:  []float64{1, math.NaN(), 3},
			weights: []float64{1, 1, 1},
			wantOK:  false,
		},
		{
			name:    "Inf weight rejected",
			values:  []float64{1, 2, 3},
			weights: []float64{1, math.Inf(1), 1},
			wantOK:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.WeightedMean(tt.values, tt.weights)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !tt.wantOK {
				if got != 0 {
					t.Errorf("result = %v, want 0 when ok is false", got)
				}
				return
			}
			if !approxEqual(got, tt.want) {
				t.Errorf("result = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeightedMeanIntegers(t *testing.T) {
	got, ok := stats.WeightedMean([]int{1, 2, 3, 4}, []int{4, 3, 2, 1})
	if !ok {
		t.Fatal("ok = false, want true")
	}
	// (1*4 + 2*3 + 3*2 + 4*1) / (4+3+2+1) = 20/10 = 2
	if !approxEqual(got, 2) {
		t.Errorf("result = %v, want 2", got)
	}
}

func TestGeometricMean(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
		wantOK bool
	}{
		{
			name:   "perfect cube",
			values: []float64{1, 10, 100},
			want:   10,
			wantOK: true,
		},
		{
			name:   "equal values return that value",
			values: []float64{5, 5, 5},
			want:   5,
			wantOK: true,
		},
		{
			name:   "two values is the square root of their product",
			values: []float64{2, 8},
			want:   4,
			wantOK: true,
		},
		{
			name:   "single value",
			values: []float64{42},
			want:   42,
			wantOK: true,
		},
		{
			name:   "empty input",
			values: []float64{},
			wantOK: false,
		},
		{
			name:   "nil input",
			values: nil,
			wantOK: false,
		},
		{
			name:   "zero value rejected",
			values: []float64{1, 0, 3},
			wantOK: false,
		},
		{
			name:   "negative value rejected",
			values: []float64{1, -2, 3},
			wantOK: false,
		},
		{
			name:   "NaN rejected",
			values: []float64{1, math.NaN(), 3},
			wantOK: false,
		},
		{
			name:   "Inf rejected",
			values: []float64{1, math.Inf(1), 3},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.GeometricMean(tt.values)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !tt.wantOK {
				if got != 0 {
					t.Errorf("result = %v, want 0 when ok is false", got)
				}
				return
			}
			if !approxEqual(got, tt.want) {
				t.Errorf("result = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeometricMeanNoOverflow(t *testing.T) {
	// A naive Π then nth-root would overflow float64; the log-space
	// implementation must not.
	values := make([]float64, 1000)
	for i := range values {
		values[i] = 1e300
	}
	got, ok := stats.GeometricMean(values)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if math.IsInf(got, 0) || math.IsNaN(got) {
		t.Fatalf("result = %v, want a finite value (no overflow)", got)
	}
	if !approxEqual(got/1e300, 1) {
		t.Errorf("result = %v, want ~1e300", got)
	}
}

func TestHarmonicMean(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
		wantOK bool
	}{
		{
			name:   "classic example",
			values: []float64{1, 2, 4},
			want:   3.0 / 1.75,
			wantOK: true,
		},
		{
			name:   "equal values return that value",
			values: []float64{5, 5, 5},
			want:   5,
			wantOK: true,
		},
		{
			name:   "single value",
			values: []float64{42},
			want:   42,
			wantOK: true,
		},
		{
			name:   "empty input",
			values: []float64{},
			wantOK: false,
		},
		{
			name:   "nil input",
			values: nil,
			wantOK: false,
		},
		{
			name:   "zero value rejected",
			values: []float64{1, 0, 4},
			wantOK: false,
		},
		{
			name:   "negative value rejected",
			values: []float64{1, -2, 4},
			wantOK: false,
		},
		{
			name:   "NaN rejected",
			values: []float64{1, math.NaN(), 4},
			wantOK: false,
		},
		{
			name:   "Inf rejected",
			values: []float64{1, math.Inf(1), 4},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := stats.HarmonicMean(tt.values)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !tt.wantOK {
				if got != 0 {
					t.Errorf("result = %v, want 0 when ok is false", got)
				}
				return
			}
			if !approxEqual(got, tt.want) {
				t.Errorf("result = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMeanOrdering checks the textbook inequality
// harmonic ≤ geometric ≤ arithmetic for positive values, which is a strong
// cross-check that all three implementations agree.
func TestMeanOrdering(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ones := make([]float64, len(values))
	for i := range ones {
		ones[i] = 1
	}

	arithmetic, ok := stats.WeightedMean(values, ones)
	if !ok {
		t.Fatal("WeightedMean ok = false")
	}
	geometric, ok := stats.GeometricMean(values)
	if !ok {
		t.Fatal("GeometricMean ok = false")
	}
	harmonic, ok := stats.HarmonicMean(values)
	if !ok {
		t.Fatal("HarmonicMean ok = false")
	}

	if !(harmonic <= geometric && geometric <= arithmetic) {
		t.Errorf("expected harmonic (%v) <= geometric (%v) <= arithmetic (%v)",
			harmonic, geometric, arithmetic)
	}
}

func FuzzWeightedMean(f *testing.F) {
	f.Add(1.0, 2.0, 3.0, 1.0, 1.0, 1.0)
	f.Fuzz(func(t *testing.T, v0, v1, v2, w0, w1, w2 float64) {
		values := []float64{v0, v1, v2}
		weights := []float64{w0, w1, w2}
		got, ok := stats.WeightedMean(values, weights)
		if !ok {
			if got != 0 {
				t.Errorf("result = %v, want 0 when ok is false", got)
			}
			return
		}
		// A successful weighted mean of finite values must itself be finite
		// and bounded by the min and max input value.
		if nonFiniteHelper(got) {
			t.Errorf("result = %v is non-finite despite ok", got)
		}
		lo := math.Min(v0, math.Min(v1, v2))
		hi := math.Max(v0, math.Max(v1, v2))
		if got < lo-epsilon || got > hi+epsilon {
			t.Errorf("result = %v outside [%v, %v]", got, lo, hi)
		}
	})
}

func nonFiniteHelper(x float64) bool {
	return math.IsNaN(x) || math.IsInf(x, 0)
}
